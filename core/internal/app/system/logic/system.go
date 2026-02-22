// Package logic 系统管理业务逻辑
// 实现版本信息查询、GitHub Releases 更新检测、一键更新、回滚等功能
package logic

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"

	jsoniter "github.com/json-iterator/go"

	"github.com/inconshreveable/go-update"

	systemApi "your-finance/allfi/api/v1/system"
	"your-finance/allfi/internal/app/system/service"
	"your-finance/allfi/internal/version"
)

// json 使用 json-iterator 替代标准库以提升性能
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 更新历史文件路径（data 目录为 Docker volume 挂载点）
const updateHistoryFile = "data/update_history.json"

// GitHub Releases API 地址
const githubReleasesURL = "https://api.github.com/repos/your-finance/allfi/releases/latest"

// Docker 内部 updater 服务地址
const dockerUpdaterURL = "http://updater:8081/update"

// sSystem 系统管理服务实现
type sSystem struct {
	mu          sync.Mutex // 保护更新状态的并发访问
	updateState string     // 当前状态：idle/updating/completed/failed
	updateStep  int        // 当前步骤序号
	updateTotal int        // 总步骤数
	stepName    string     // 当前步骤名称
	updateMsg   string     // 详细信息
}

// New 创建系统管理服务实例
func New() service.ISystem {
	return &sSystem{
		updateState: "idle",
	}
}

// GetVersion 获取当前版本信息
//
// 功能说明:
// 1. 从 version 包获取构建时注入的版本号、构建时间、Git 提交哈希
// 2. 通过检测 /.dockerenv 文件判断运行模式（docker/host）
// 3. 通过 runtime.Version() 获取 Go 版本
func (s *sSystem) GetVersion(ctx context.Context) (*systemApi.GetVersionRes, error) {
	// 检测运行模式：Docker 容器内存在 /.dockerenv 文件
	runMode := "host"
	if _, err := os.Stat("/.dockerenv"); err == nil {
		runMode = "docker"
	}

	return &systemApi.GetVersionRes{
		Version:   version.Version,
		BuildTime: version.BuildTime,
		GitCommit: version.GitCommit,
		RunMode:   runMode,
		GoVersion: runtime.Version(),
	}, nil
}

// CheckUpdate 检查 GitHub Releases 是否有新版本
//
// 功能说明:
// 1. 通过 GitHub API 获取最新 Release 信息
// 2. 解析 tag_name（格式为 v0.3.0），去掉 v 前缀后与当前版本比较
// 3. 使用语义版本比较（major.minor.patch 逐级比较）
// 4. 网络错误或 API 异常时返回 has_update=false + 错误信息
func (s *sSystem) CheckUpdate(ctx context.Context) (*systemApi.CheckUpdateRes, error) {
	currentVer := version.Version

	// 构造带超时的 HTTP 请求
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubReleasesURL, nil)
	if err != nil {
		g.Log().Warning(ctx, "构造 GitHub API 请求失败", "error", err)
		return &systemApi.CheckUpdateRes{
			HasUpdate:      false,
			CurrentVersion: currentVer,
			LatestVersion:  "",
			ReleaseNotes:   fmt.Sprintf("检查更新失败: %v", err),
		}, nil
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "AllFi/"+currentVer)

	resp, err := client.Do(req)
	if err != nil {
		g.Log().Warning(ctx, "请求 GitHub API 失败", "error", err)
		return &systemApi.CheckUpdateRes{
			HasUpdate:      false,
			CurrentVersion: currentVer,
			LatestVersion:  "",
			ReleaseNotes:   fmt.Sprintf("网络请求失败: %v", err),
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		g.Log().Warning(ctx, "GitHub API 返回非 200 状态码", "status", resp.StatusCode)
		return &systemApi.CheckUpdateRes{
			HasUpdate:      false,
			CurrentVersion: currentVer,
			LatestVersion:  "",
			ReleaseNotes:   fmt.Sprintf("GitHub API 返回状态码 %d", resp.StatusCode),
		}, nil
	}

	// 解析 GitHub Release 响应
	var release struct {
		TagName     string `json:"tag_name"`
		Body        string `json:"body"`
		HTMLURL     string `json:"html_url"`
		PublishedAt string `json:"published_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		g.Log().Warning(ctx, "解析 GitHub Release 响应失败", "error", err)
		return &systemApi.CheckUpdateRes{
			HasUpdate:      false,
			CurrentVersion: currentVer,
			LatestVersion:  "",
			ReleaseNotes:   fmt.Sprintf("解析响应失败: %v", err),
		}, nil
	}

	// tag_name 格式为 v0.3.0，去掉 v 前缀
	latestVer := strings.TrimPrefix(release.TagName, "v")

	// 语义版本比较
	hasUpdate := compareSemanticVersion(latestVer, currentVer) > 0

	return &systemApi.CheckUpdateRes{
		HasUpdate:      hasUpdate,
		CurrentVersion: currentVer,
		LatestVersion:  latestVer,
		ReleaseNotes:   release.Body,
		ReleaseURL:     release.HTMLURL,
		PublishedAt:    release.PublishedAt,
	}, nil
}

// ApplyUpdate 执行版本更新
//
// 功能说明:
// 1. 检测运行模式，Docker 模式转发给 updater 服务
// 2. 宿主机模式下异步执行 OTA 二进制热替换更新
// 3. 使用 sync.Mutex 保护更新状态
// 4. 立即返回 status: "started"，后台执行实际更新操作
func (s *sSystem) ApplyUpdate(ctx context.Context, targetVersion string) (*systemApi.ApplyUpdateRes, error) {
	s.mu.Lock()
	// 检查是否已在更新中
	if s.updateState == "updating" {
		s.mu.Unlock()
		return &systemApi.ApplyUpdateRes{
			Status:  "failed",
			Message: "已有更新任务在执行中，请等待完成",
		}, nil
	}
	s.mu.Unlock()

	// 检测运行模式
	if _, err := os.Stat("/.dockerenv"); err == nil {
		// Docker 模式：转发给 updater 服务
		return s.applyUpdateDocker(ctx, targetVersion)
	}

	// 宿主机模式：异步执行 OTA 二进制热替换
	s.setUpdateState("updating", 0, 3, "准备更新", fmt.Sprintf("目标版本: %s", targetVersion))

	go s.doHostUpdate(targetVersion, false)

	return &systemApi.ApplyUpdateRes{
		Status:  "started",
		Message: fmt.Sprintf("正在更新到版本 %s", targetVersion),
	}, nil
}

// Rollback 回滚到指定历史版本
//
// 功能说明:
// 与 ApplyUpdate 类似，但在历史记录中标记为回滚操作
// Docker 模式转发给 updater 服务，宿主机模式执行 git checkout
func (s *sSystem) Rollback(ctx context.Context, targetVersion string) (*systemApi.RollbackRes, error) {
	s.mu.Lock()
	if s.updateState == "updating" {
		s.mu.Unlock()
		return &systemApi.RollbackRes{
			Status:  "failed",
			Message: "已有更新任务在执行中，请等待完成",
		}, nil
	}
	s.mu.Unlock()

	// 检测运行模式
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return s.rollbackDocker(ctx, targetVersion)
	}

	// 宿主机模式：异步执行 OTA 二进制热替换回滚
	s.setUpdateState("updating", 0, 3, "准备回滚", fmt.Sprintf("目标版本: %s", targetVersion))

	go s.doHostUpdate(targetVersion, true)

	return &systemApi.RollbackRes{
		Status:  "started",
		Message: fmt.Sprintf("正在回滚到版本 %s", targetVersion),
	}, nil
}

// GetUpdateStatus 获取当前更新/回滚操作的进度
func (s *sSystem) GetUpdateStatus(ctx context.Context) (*systemApi.GetUpdateStatusRes, error) {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return s.getUpdateStatusDocker(ctx)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return &systemApi.GetUpdateStatusRes{
		State:    s.updateState,
		Step:     s.updateStep,
		Total:    s.updateTotal,
		StepName: s.stepName,
		Message:  s.updateMsg,
	}, nil
}

// getUpdateStatusDocker Docker 模式下获取 updater 服务的状态
func (s *sSystem) getUpdateStatusDocker(ctx context.Context) (*systemApi.GetUpdateStatusRes, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dockerUpdaterURL+"/status", nil)
	if err != nil {
		return nil, fmt.Errorf("构造 updater 请求失败: %w", err)
	}
	// dockerUpdaterURL 是 "http://updater:8081/update"，所以需要替换掉 path
	req.URL.Path = "/status"

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("连接 updater 服务失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("updater 服务返回状态码 %d", resp.StatusCode)
	}

	var status struct {
		State    string `json:"state"`
		Step     int    `json:"step"`
		Total    int    `json:"total"`
		StepName string `json:"step_name"`
		Message  string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("解析 updater 状态失败: %w", err)
	}

	return &systemApi.GetUpdateStatusRes{
		State:    status.State,
		Step:     status.Step,
		Total:    status.Total,
		StepName: status.StepName,
		Message:  status.Message,
	}, nil
}

// GetUpdateHistory 获取历史更新记录
//
// 从 data/update_history.json 文件中读取更新记录
// 如果文件不存在，返回空列表
func (s *sSystem) GetUpdateHistory(ctx context.Context) (*systemApi.GetUpdateHistoryRes, error) {
	records, err := s.loadHistory()
	if err != nil {
		g.Log().Warning(ctx, "读取更新历史失败", "error", err)
		return &systemApi.GetUpdateHistoryRes{
			Records: []systemApi.UpdateRecord{},
		}, nil
	}

	return &systemApi.GetUpdateHistoryRes{
		Records: records,
	}, nil
}

// ====================== 内部辅助方法 ======================

// setUpdateState 设置更新状态（线程安全）
func (s *sSystem) setUpdateState(state string, step, total int, stepName, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.updateState = state
	s.updateStep = step
	s.updateTotal = total
	s.stepName = stepName
	s.updateMsg = msg
}

// doHostUpdate 宿主机模式下执行 OTA 更新/回滚
//
// 步骤:
// 1. 根据当前架构下载对应平台的 Release 压缩包 (.tar.gz)
// 2. 解压并提取 allfi 二进制文件进行自身覆盖替换 (Self-Update)
// 3. 执行成功后重启进程
func (s *sSystem) doHostUpdate(targetVersion string, isRollback bool) {
	actionName := "更新"
	if isRollback {
		actionName = "回滚"
	}

	// 记录当前版本到历史 (预存 started 状态)
	record := systemApi.UpdateRecord{
		Version:   targetVersion,
		GitCommit: version.GitCommit,
		Timestamp: time.Now().Format(time.RFC3339),
		Status:    "started",
	}

	// 1. 拼接要下载的 Release 包全名
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	tarballName := fmt.Sprintf("allfi-%s-%s.tar.gz", targetVersion, platform)
	downloadURL := fmt.Sprintf("https://github.com/your-finance/allfi/releases/download/v%s/%s", targetVersion, tarballName)

	s.setUpdateState("updating", 1, 3, "下载新版本", fmt.Sprintf("正在下载包: %s", tarballName))
	g.Log().Info(context.Background(), fmt.Sprintf("准备执行 %s，下载地址: %s", actionName, downloadURL))

	// 2. 发起下载请求
	client := &http.Client{Timeout: 300 * time.Second} // 下载二进制包需要较长超时
	resp, err := client.Get(downloadURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("下载 %s 失败，请检查网络（状态码: %d）", tarballName, resp.StatusCode)
		if err != nil {
			errMsg = fmt.Sprintf("下载异常被中断: %v", err)
		}
		g.Log().Error(context.Background(), errMsg)
		s.setUpdateState("failed", 1, 3, "下载失败", errMsg)
		record.Status = "failed"
		_ = s.appendHistory(record)
		return
	}
	defer resp.Body.Close()

	// 3. 解析 tar.gz，并在里面找到可执行文件 "allfi"
	s.setUpdateState("updating", 2, 3, "替换新版本", "正在提取执行文件并替换当前进程...")
	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		s.setUpdateState("failed", 2, 3, "读取压缩包失败", fmt.Sprintf("解压包异常: %v", err))
		record.Status = "failed"
		_ = s.appendHistory(record)
		return
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	var binaryReader io.Reader
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.setUpdateState("failed", 2, 3, "解析压缩包失败", fmt.Sprintf("读取错误: %v", err))
			record.Status = "failed"
			_ = s.appendHistory(record)
			return
		}

		// release 结构是 "allfi-{version}-{platform}/allfi"
		if !header.FileInfo().IsDir() && filepath.Base(header.Name) == "allfi" {
			// 对于 tarball 中名为 allfi 的可扫码可执行文件截获
			binaryReader = tr
			break
		}
	}

	if binaryReader == nil {
		s.setUpdateState("failed", 2, 3, "更新失败", "在安装包中未找到 allfi 二进制文件")
		record.Status = "failed"
		_ = s.appendHistory(record)
		return
	}

	// 4. 执行原位替换（将自身的可执行文件无缝覆盖）
	err = update.Apply(binaryReader, update.Options{})
	if err != nil {
		s.setUpdateState("failed", 2, 3, "覆盖二进制文件失败", fmt.Sprintf("替换发生错误: %v", err))
		record.Status = "failed"
		_ = s.appendHistory(record)
		return
	}

	// 5. 更新成功，标记重启
	if isRollback {
		record.Status = "rolled_back"
	} else {
		record.Status = "success"
	}
	_ = s.appendHistory(record)

	s.setUpdateState("completed", 3, 3, fmt.Sprintf("%s完成", actionName), fmt.Sprintf("已可平滑重启应用新版本 %s", targetVersion))

	g.Log().Info(context.Background(), fmt.Sprintf("OTA 自升级（%s）已成功将二进制文件替换到版本 %s", actionName, targetVersion))

	// 重启当前服务进程
	go func() {
		// 休眠2秒让接口成功返回 status: "completed" 给前端后，再干掉老进程
		time.Sleep(2 * time.Second)
		exe, err := os.Executable()
		if err != nil {
			g.Log().Error(context.Background(), "重启失败，无法找到可执行文件位置: ", err)
			return
		}

		// 将老参数原样带给新进程重启
		cmd := exec.Command(exe, os.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		g.Log().Info(context.Background(), "开始重启，启动新版本进程...")
		if err := cmd.Start(); err != nil {
			g.Log().Error(context.Background(), "新版本进程启动失败: ", err)
			return
		}
		// 正确退出旧进程
		os.Exit(0)
	}()
}

// applyUpdateDocker Docker 模式下转发更新请求给 updater 服务
func (s *sSystem) applyUpdateDocker(ctx context.Context, targetVersion string) (*systemApi.ApplyUpdateRes, error) {
	body, _ := json.Marshal(map[string]string{
		"target_version": targetVersion,
	})

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, dockerUpdaterURL, bytes.NewReader(body))
	if err != nil {
		return &systemApi.ApplyUpdateRes{
			Status:  "failed",
			Message: fmt.Sprintf("构造 updater 请求失败: %v", err),
		}, nil
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return &systemApi.ApplyUpdateRes{
			Status:  "failed",
			Message: fmt.Sprintf("连接 updater 服务失败: %v", err),
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &systemApi.ApplyUpdateRes{
			Status:  "failed",
			Message: fmt.Sprintf("updater 服务返回状态码 %d", resp.StatusCode),
		}, nil
	}

	return &systemApi.ApplyUpdateRes{
		Status:  "started",
		Message: fmt.Sprintf("Docker 更新已触发，目标版本: %s", targetVersion),
	}, nil
}

// rollbackDocker Docker 模式下转发回滚请求给 updater 服务
func (s *sSystem) rollbackDocker(ctx context.Context, targetVersion string) (*systemApi.RollbackRes, error) {
	body, _ := json.Marshal(map[string]string{
		"target_version": targetVersion,
		"action":         "rollback",
	})

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, dockerUpdaterURL, bytes.NewReader(body))
	if err != nil {
		return &systemApi.RollbackRes{
			Status:  "failed",
			Message: fmt.Sprintf("构造 updater 请求失败: %v", err),
		}, nil
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return &systemApi.RollbackRes{
			Status:  "failed",
			Message: fmt.Sprintf("连接 updater 服务失败: %v", err),
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &systemApi.RollbackRes{
			Status:  "failed",
			Message: fmt.Sprintf("updater 服务返回状态码 %d", resp.StatusCode),
		}, nil
	}

	return &systemApi.RollbackRes{
		Status:  "started",
		Message: fmt.Sprintf("Docker 回滚已触发，目标版本: %s", targetVersion),
	}, nil
}

// loadHistory 从 JSON 文件加载更新历史记录
// 如果文件不存在，返回空切片
func (s *sSystem) loadHistory() ([]systemApi.UpdateRecord, error) {
	if !gfile.Exists(updateHistoryFile) {
		return []systemApi.UpdateRecord{}, nil
	}

	data, err := os.ReadFile(updateHistoryFile)
	if err != nil {
		return nil, fmt.Errorf("读取更新历史文件失败: %w", err)
	}

	var records []systemApi.UpdateRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("解析更新历史 JSON 失败: %w", err)
	}

	return records, nil
}

// appendHistory 追加一条更新记录到历史文件
// 自动创建 data 目录（如果不存在）
func (s *sSystem) appendHistory(record systemApi.UpdateRecord) error {
	// 确保 data 目录存在
	dir := filepath.Dir(updateHistoryFile)
	if !gfile.Exists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建 data 目录失败: %w", err)
		}
	}

	// 加载已有记录
	records, err := s.loadHistory()
	if err != nil {
		// 读取失败时从空列表开始
		records = []systemApi.UpdateRecord{}
	}

	// 追加新记录
	records = append(records, record)

	// 序列化并写入文件
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化更新历史失败: %w", err)
	}

	if err := os.WriteFile(updateHistoryFile, data, 0644); err != nil {
		return fmt.Errorf("写入更新历史文件失败: %w", err)
	}

	return nil
}

// compareSemanticVersion 语义版本比较
// 返回值: >0 表示 a > b, =0 表示 a == b, <0 表示 a < b
// 版本格式: major.minor.patch（如 0.3.0）
// 非法版本字符串返回 0（视为相等）
func compareSemanticVersion(a, b string) int {
	aParts := parseVersion(a)
	bParts := parseVersion(b)

	for i := 0; i < 3; i++ {
		if aParts[i] != bParts[i] {
			return aParts[i] - bParts[i]
		}
	}
	return 0
}

// parseVersion 解析版本号为 [major, minor, patch] 数组
// 非法格式返回 [0, 0, 0]
func parseVersion(v string) [3]int {
	var parts [3]int
	segments := strings.Split(v, ".")
	for i := 0; i < len(segments) && i < 3; i++ {
		n, err := strconv.Atoi(segments[i])
		if err != nil {
			return [3]int{}
		}
		parts[i] = n
	}
	return parts
}
