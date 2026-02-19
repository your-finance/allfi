package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	// mu 保护更新状态的互斥锁
	mu sync.Mutex
	// state 当前状态：idle / updating / completed / failed
	state = "idle"
	// step 当前步骤序号
	step int
	// total 总步骤数
	total int
	// stepName 当前步骤名称
	stepName string
	// message 状态描述信息
	message string
)

// projectDir 项目目录（docker-compose volume 挂载点）
const projectDir = "/app/project"

// getStatus 获取当前更新/回滚进度（线程安全）
func getStatus() statusResponse {
	mu.Lock()
	defer mu.Unlock()
	return statusResponse{
		State:    state,
		Step:     step,
		Total:    total,
		StepName: stepName,
		Message:  message,
	}
}

// setState 设置当前更新/回滚进度（线程安全）
func setState(s string, st int, tot int, name string, msg string) {
	mu.Lock()
	defer mu.Unlock()
	state = s
	step = st
	total = tot
	stepName = name
	message = msg
}

// buildBackendImage 直接使用 docker build 构建后端镜像
// docker-compose.yml 中 build 被注释掉（使用预构建镜像模式），
// 所以必须直接调用 docker build 并打上 allfi-backend:latest 标签
func buildBackendImage(targetVersion string) error {
	buildTime := time.Now().Format("2006-01-02T15:04:05")
	// 获取当前 git commit hash
	gitCommit := "unknown"
	if out, err := exec.Command("git", "-C", projectDir, "rev-parse", "--short", "HEAD").Output(); err == nil {
		gitCommit = strings.TrimSpace(string(out))
	}

	return runCmd(projectDir, "docker", "build",
		"-f", "core/Dockerfile",
		"--build-arg", "APP_VERSION="+targetVersion,
		"--build-arg", "BUILD_TIME="+buildTime,
		"--build-arg", "GIT_COMMIT="+gitCommit,
		"-t", "allfi-backend:latest",
		".",
	)
}

// restartBackend 停止旧容器并启动新容器
func restartBackend() error {
	// 先尝试 docker compose（v2 插件），再回退 docker-compose（v1）
	_ = runCmd(projectDir, "docker", "compose", "rm", "-f", "-s", "backend")
	err := runCmd(projectDir, "docker", "compose", "up", "-d", "backend")
	if err == nil {
		return nil
	}
	_ = runCmd(projectDir, "docker-compose", "rm", "-f", "-s", "backend")
	return runCmd(projectDir, "docker-compose", "up", "-d", "backend")
}

// runCmd 在指定目录执行命令，返回错误信息（包含命令输出）
func runCmd(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	log.Printf("[CMD] %s %v -> OK", name, args)
	return nil
}

// doUpdate 执行更新流程：git fetch + checkout tag + docker build + restart
func doUpdate(targetVersion string) {
	setState("updating", 0, 4, "准备中", "开始更新到 v"+targetVersion)
	log.Printf("开始更新到 v%s", targetVersion)

	// 步骤 1：拉取代码
	setState("updating", 1, 4, "拉取代码", "执行 git fetch...")
	if err := runCmd(projectDir, "git", "fetch", "--tags"); err != nil {
		setState("failed", 1, 4, "拉取代码", fmt.Sprintf("git fetch 失败: %v", err))
		log.Printf("git fetch 失败: %v", err)
		return
	}

	// 步骤 2：切换到目标版本
	tag := "v" + targetVersion
	setState("updating", 2, 4, "切换版本", fmt.Sprintf("checkout %s...", tag))
	if err := runCmd(projectDir, "git", "checkout", tag); err != nil {
		setState("failed", 2, 4, "切换版本", fmt.Sprintf("git checkout %s 失败: %v", tag, err))
		log.Printf("git checkout %s 失败: %v", tag, err)
		return
	}

	// 步骤 3：构建镜像（直接 docker build，不依赖 docker-compose build）
	setState("updating", 3, 4, "构建镜像", "正在构建后端镜像，请耐心等待...")
	if err := buildBackendImage(targetVersion); err != nil {
		setState("failed", 3, 4, "构建镜像", fmt.Sprintf("构建失败: %v", err))
		log.Printf("构建镜像失败: %v", err)
		return
	}

	// 步骤 4：重启服务
	setState("updating", 4, 4, "重启服务", "停止旧容器并启动新容器...")
	if err := restartBackend(); err != nil {
		setState("failed", 4, 4, "重启服务", fmt.Sprintf("重启失败: %v", err))
		log.Printf("重启服务失败: %v", err)
		return
	}

	setState("completed", 4, 4, "完成", fmt.Sprintf("已成功更新到 v%s", targetVersion))
	log.Printf("更新到 v%s 完成", targetVersion)
}

// doRollback 执行回滚流程：git checkout tag + docker build + restart
func doRollback(targetVersion string) {
	setState("updating", 0, 3, "准备中", "开始回滚到 v"+targetVersion)
	log.Printf("开始回滚到 v%s", targetVersion)

	// 步骤 1：切换到目标版本
	tag := "v" + targetVersion
	setState("updating", 1, 3, "切换版本", fmt.Sprintf("checkout %s...", tag))
	if err := runCmd(projectDir, "git", "checkout", tag); err != nil {
		setState("failed", 1, 3, "切换版本", fmt.Sprintf("git checkout %s 失败: %v", tag, err))
		return
	}

	// 步骤 2：构建镜像
	setState("updating", 2, 3, "构建镜像", "正在构建后端镜像，请耐心等待...")
	if err := buildBackendImage(targetVersion); err != nil {
		setState("failed", 2, 3, "构建镜像", fmt.Sprintf("构建失败: %v", err))
		return
	}

	// 步骤 3：重启服务
	setState("updating", 3, 3, "重启服务", "停止旧容器并启动新容器...")
	if err := restartBackend(); err != nil {
		setState("failed", 3, 3, "重启服务", fmt.Sprintf("重启失败: %v", err))
		return
	}

	setState("completed", 3, 3, "完成", fmt.Sprintf("已回滚到 v%s", targetVersion))
	log.Printf("回滚到 v%s 完成", targetVersion)
}
