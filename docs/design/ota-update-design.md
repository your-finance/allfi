# AllFi OTA 自动更新设计方案

## 1. 背景与痛点

当前的客户端端程序（特别是宿主机运行模式下），当检测到有新版本并执行 `ApplyUpdate` 时，采用的是执行 `git fetch --tags` 和 `git checkout <tag>` 的方式。
这种更新机制的弊端包括：
- **依赖运行环境包含开发工具链**：如系统中必须安装并配置好了 Git、Go 甚至 pnpm。如果终端用户并没有完整的开发与编译环境，则更新必定失败。
- **本地源代码替换，但进程并无热更新**：即使 `git checkout` 成功了，在当前程序内存里执行的依然是旧代码进程，必须由外部环境去关闭旧进程并使用开发环境工具链（如 `make build` 或者 `go run`）才能真正运行新版本，做不到即点即关即开。

## 2. OTA（Over-The-Air）自升级方案概述

AllFi 已经有了通过 `.github/workflows/release.yml` 自动化跨平台打包的流水线，前端 `webapp` 所编译的 UI 已经被放入到 Go 的单体可执行文件里（借助 `//go:embed` 技术）。这意味着整个应用可以仅通过分发单一的二进制执行文件来进行升级。

因此，宿主机的自我更新流程应当优化为 **二进流原位替换与进程自启（Self-Update）**：
1. **下载构建产物**：系统发起更新时，应用内置服务直接访问 GitHub Releases，获取当前系统平台对应的稳定版预构建压缩包（例如 `allfi-1.0.0-darwin-arm64.tar.gz`）。
2. **二进流替换**：在内存中将 `.tar.gz` 解压出最新的 `allfi` 二进制流，并利用 `go-update` 机制原子级地重写并替换当前正在执行的自身二进制文件。
3. **平滑重启**：由主应用调用操作系统的底层的系统接口命令来挂起/重启当前进程（`exec.Command(exe)` 或 `syscall.Exec`），实现程序的无缝重启，更新过程做到“零配置、零外部依赖”。

## 3. 具体实施步骤

### 3.1 核心依赖引入
引入 `github.com/inconshreveable/go-update` 作为原位无缝替换库。可以以原子操作覆盖操作系统可执行文件并在替换异常时安全回退。

### 3.2 改造 `doHostUpdate` 方法
修改 `core/internal/app/system/logic/system.go` 中的宿主机更新处理逻辑，摒弃所有的 `git...` 逻辑：

1. 判断当前架构平台，拼装对应 Release 包的文件名。
   - 例: `darwin-arm64`, `linux-amd64`
2. 下载远程的 `tar.gz` 压缩包并提取二进制文件。
3. 应用二进制覆盖：
   ```go
   err = update.Apply(binaryReader, update.Options{})
   ```
4. 执行自身重启完成更新：
   ```go
   exe, _ := os.Executable()
   cmd := exec.Command(exe, os.Args[1:]...)
   cmd.Start()
   os.Exit(0)
   ```

### 3.3 Docker 零依赖部署脚本
`deploy/docker-deploy.sh` 改为完全免源码方案：
- 通过 `curl -sSL ... | bash` 可在任意 Linux 服务器上即装即用，无需 `git clone` 源码。
- 脚本自动检测 CPU 架构 (amd64/arm64)，从 GitHub Releases 下载对应平台的预编译二进制包。
- 自动生成 `.env`（含随机安全密钥）+ `Dockerfile` + `docker-compose.yml`。
- 构建轻量级 Alpine 镜像（仅打包二进制，无编译过程），并启动服务。
- 后续版本更新可以直接在前端页面中 OTA 一键完成。

### 3.4 文档及脚本更新
- **部署脚本**（如 `quickstart.sh`）：宿主机一键启动不再推荐纯源码运行，而应优先推荐直接拉取二进制文件运行，因为包含自升级模块，此后一切均为开箱即用的更新模式。
- **README** 更新：声明支持无需编译的一键式跨平台 OTA 自更新（即页面点击更新直接平滑应用最新后端功能和前端界面）。

## 4. 优势总结
1. **Zero-Dependency**：不依赖终端用户的 Go/Git 环境。
2. **Immediately Visible**：真正的"一键更新"，前端点击更新并稍等刷新以后，瞬间感知代码最新逻辑。
3. **Rollback一致**：同样的代码实现一键安全降级回滚。
