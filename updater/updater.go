package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
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

// doUpdate 执行更新流程：git fetch + checkout tag + docker-compose build + restart
func doUpdate(targetVersion string) {
	setState("updating", 0, 3, "准备中", "开始更新到 v"+targetVersion)
	log.Printf("开始更新到 v%s", targetVersion)

	// 步骤 1：拉取代码
	setState("updating", 1, 3, "拉取代码", "执行 git fetch...")
	if err := runCmd(projectDir, "git", "fetch", "--tags"); err != nil {
		setState("failed", 1, 3, "拉取代码", fmt.Sprintf("git fetch 失败: %v", err))
		log.Printf("git fetch 失败: %v", err)
		return
	}

	// checkout 到目标版本 tag
	tag := "v" + targetVersion
	if err := runCmd(projectDir, "git", "checkout", tag); err != nil {
		setState("failed", 1, 3, "拉取代码", fmt.Sprintf("git checkout %s 失败: %v", tag, err))
		log.Printf("git checkout %s 失败: %v", tag, err)
		return
	}

	// 步骤 2：构建镜像
	setState("updating", 2, 3, "构建镜像", "执行 docker-compose build...")
	if err := runCmd(projectDir, "docker-compose", "build", "backend"); err != nil {
		// 也尝试 docker compose（新版命令）
		if err2 := runCmd(projectDir, "docker", "compose", "build", "backend"); err2 != nil {
			setState("failed", 2, 3, "构建镜像", fmt.Sprintf("docker-compose build 失败: %v", err2))
			log.Printf("docker-compose build 失败: %v", err2)
			return
		}
	}

	// 步骤 3：重启服务（注意：不重启 updater 自身）
	// 先停止并移除旧容器，避免容器名称冲突
	setState("updating", 3, 3, "重启服务", "停止旧容器...")
	_ = runCmd(projectDir, "docker-compose", "rm", "-f", "-s", "backend")
	_ = runCmd(projectDir, "docker", "compose", "rm", "-f", "-s", "backend")

	setState("updating", 3, 3, "重启服务", "启动新容器...")
	if err := runCmd(projectDir, "docker-compose", "up", "-d", "backend"); err != nil {
		if err2 := runCmd(projectDir, "docker", "compose", "up", "-d", "backend"); err2 != nil {
			setState("failed", 3, 3, "重启服务", fmt.Sprintf("重启失败: %v", err2))
			log.Printf("重启服务失败: %v", err2)
			return
		}
	}

	setState("completed", 3, 3, "完成", fmt.Sprintf("已成功更新到 v%s", targetVersion))
	log.Printf("更新到 v%s 完成", targetVersion)
}

// doRollback 执行回滚流程：git checkout tag + docker-compose build + restart
func doRollback(targetVersion string) {
	setState("updating", 0, 3, "准备中", "开始回滚到 v"+targetVersion)
	log.Printf("开始回滚到 v%s", targetVersion)

	// 步骤 1：切换到目标版本
	setState("updating", 1, 3, "切换版本", "执行 git checkout...")
	tag := "v" + targetVersion
	if err := runCmd(projectDir, "git", "checkout", tag); err != nil {
		setState("failed", 1, 3, "切换版本", fmt.Sprintf("git checkout %s 失败: %v", tag, err))
		return
	}

	// 步骤 2：重建镜像
	setState("updating", 2, 3, "构建镜像", "重新构建...")
	if err := runCmd(projectDir, "docker-compose", "build", "backend"); err != nil {
		if err2 := runCmd(projectDir, "docker", "compose", "build", "backend"); err2 != nil {
			setState("failed", 2, 3, "构建镜像", fmt.Sprintf("构建失败: %v", err2))
			return
		}
	}

	// 步骤 3：重启服务
	// 先停止并移除旧容器，避免容器名称冲突
	setState("updating", 3, 3, "重启服务", "停止旧容器...")
	_ = runCmd(projectDir, "docker-compose", "rm", "-f", "-s", "backend")
	_ = runCmd(projectDir, "docker", "compose", "rm", "-f", "-s", "backend")

	setState("updating", 3, 3, "重启服务", "启动新容器...")
	if err := runCmd(projectDir, "docker-compose", "up", "-d", "backend"); err != nil {
		if err2 := runCmd(projectDir, "docker", "compose", "up", "-d", "backend"); err2 != nil {
			setState("failed", 3, 3, "重启服务", fmt.Sprintf("重启失败: %v", err2))
			return
		}
	}

	setState("completed", 3, 3, "完成", fmt.Sprintf("已回滚到 v%s", targetVersion))
	log.Printf("回滚到 v%s 完成", targetVersion)
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
