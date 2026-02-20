# AllFi 版本历史

> 详细记录每个版本的变更内容。版本号与根目录 `VERSION` 文件和 Git Tag 保持一致。

---







## v0.1.6（2026-02-20）

### 功能新增

- 添加公共 RPC 免费 Gas 查询 + Web 端 API Key 管理

### 问题修复

- 修复更新功能：改用 docker build 直接构建镜像，解决容器名冲突
- 修复升级时容器名称冲突导致重启失败的问题

### 其他变更

- 统一服务端口：移除 SERVER_PORT，合并为 ALLFI_PORT

## v0.1.5（2026-02-17）

### 问题修复

- 修复 Docker 健康检查返回 404 的问题
- 修复前端版本显示问题：优先使用后端 API 返回的版本号
- 修复 Docker 部署内存不足问题
- 修复 updater 更新功能：添加 Git HTTPS 配置

### 其他变更

- 将系统管理 API 移至免认证区域，修复前端无法获取后端版本号的问题
- 统一 Docker 部署脚本，合并 build-local.sh 到 docker-deploy.sh

## v0.1.4（2026-02-17）

### 功能新增

- 添加 CEX 交易所动态列表功能

### 问题修复

- 修复 dashboard 费用分析组件控制台报错
- 修复版本号显示和侧边栏折叠后无法展开的问题

### 文档更新

- docs: 自动更新 VERSION 和 CHANGEHISTORY.md — v0.1.4
- docs: 自动更新 VERSION 和 CHANGEHISTORY.md — v0.1.3

### 其他变更

- 重新设计费用分析省钱建议区域
- 移除策略面板，完善费用分析
- 优化仪表盘和分享卡片功能
- 优化侧边栏交互和版本信息展示

## v0.1.4（2026-02-17）

### 问题修复

- 修复版本号显示和侧边栏折叠后无法展开的问题

### 文档更新

- docs: 自动更新 VERSION 和 CHANGEHISTORY.md — v0.1.3

### 其他变更

- 重新设计费用分析省钱建议区域
- 移除策略面板，完善费用分析
- 优化仪表盘和分享卡片功能
- 优化侧边栏交互和版本信息展示

## v0.1.3（2026-02-17）

### 功能新增

- feat: 侧边栏集成版本检查徽章，删除底部版本号
- feat: 创建 VersionBadge 版本徽章组件

### 文档更新

- docs: 自动更新 VERSION 和 CHANGEHISTORY.md — v0.1.2

### 其他变更

- i18n: 添加 viewReleases 翻译键
- 版本号检查
- 版本号双 v，asset_snapshots 插入失败

## v0.1.2（2026-02-16）

### 功能新增

- feat: release 流水线集成 CHANGEHISTORY.md 自动更新

### 构建与部署

- ci: VERSION 文件随 git tag 自动更新

### 其他变更

- 容器时区现在会与宿主机保持一致

## v0.1.0（Alpha）

### 核心资产管理
- CEX 交易所资产聚合（Binance / OKX / Coinbase）
- Web3 链上资产（Ethereum / BSC / Polygon）
- 传统资产（银行账户 / 现金 / 股票 / 基金）
- 手动资产管理（增删改查）
- 多币种计价（USDC / BTC / ETH / CNY）
- 历史快照 + 趋势图表

### DeFi / NFT
- DeFi 仓位追踪（Lido / RocketPool / Aave / Compound / Uniswap V2+V3 / Curve）
- NFT 资产管理（Alchemy 集成，浏览和估值）

### 交易与分析
- 交易记录聚合（CEX + 链上统一格式）
- 增量同步 + cursor 分页 + 复合索引优化
- 交易日汇总（TransactionDailySummary 预聚合）
- 费用分析（手续费 + Gas 消耗追踪）
- 基准对比（收益率 vs BTC/ETH/S&P500）

### 策略与成就
- 策略引擎（目标配比 + 再平衡建议）
- 成就系统（11 项投资成就解锁）

### 报告与通知
- 资产报告（日报 / 周报 / 月报 / 年报自动生成）
- 通知系统（应用内通知 + WebPush 推送）
- 价格预警

### 前端体验
- 隐私模式增强（顶栏眼睛图标 + themeStore.privacyMode + Ctrl+H 快捷键）
- 智能资产分组（弹药库/核心持仓/风险资产，可折叠分组视图）
- Gas 加油站小组件（顶栏 Gas 价格胶囊，30s 刷新，颜色标识拥堵等级）
- 4 套专业金融主题（深色/浅色）
- 3 语言国际化（zh-CN / zh-TW / en）
- PWA 支持（离线可用）
- 命令面板（Command Palette）
- 仪表盘自定义布局
- 响应式设计 + 移动端适配
- Toast 反馈系统
- 交易同步设置面板（Settings.vue）

### 体验增强
- 主播模式增强（隐私模式 + Ctrl+H 快捷键）
- 智能资产分组（弹药库/核心持仓/风险资产）
- 多链 Gas 加油站（ETH/BSC/Polygon）
- Docker Compose 一键部署（只读容器 + 非特权）
- PIN 认证（bcrypt + JWT + 锁定保护）
- 移动端优化（底部导航栏 + 下拉刷新）
- 风险提醒服务（集中度/平台/波动/缓冲 4 类规则）
- 趋势预测 API（线性回归 + R² 置信度）
- 每日盈亏 + 资产归因分析 + 资产健康评分
- 目标追踪后端持久化
- OpenAPI 3.0 文档 + Swagger UI（GoFrame 内置，自动从 `g.Meta` 生成）

### 前后端对接
- 前后端 API 全量对接（Mock → Real API 切换）
- API 服务层响应适配器（后端格式 → 前端 Store 期望格式）
- 资产总览聚合适配器（5 接口并发聚合）
- CEX 账户/钱包地址/手动资产字段映射
- 汇率/通知/价格预警/交易记录/分析模块对接
- DeFi/NFT/策略引擎/成就系统/报告系统对接
- 用户设置/目标追踪/健康评分对接
- 环境变量切换支持（.env / .env.mock）
- goalStore 从 localStorage 迁移到后端 CRUD API
- WebPush 前端服务（VAPID 获取 + 订阅/取消订阅）
- 认证中间件修复（JWT 密钥来源与 auth logic 统一为 DB）
- 通知偏好查询修复（处理首次使用时的空记录）
- 实体模型清理（gorm.DeletedAt → *time.Time，移除 GORM 依赖）
- 59 个后端端点全部完成前端 API 对接（100% 覆盖）

### 开发体验优化
- Makefile 统一命令入口（make dev / make setup / make docker 等）
- 快速启动脚本（scripts/quickstart.sh 依赖检测 + 一键启动）
- 前端 Mock 快速体验模式（pnpm dev:mock 无需后端）
- 代理配置指南（Go/pnpm/Docker/外部 API 全覆盖）
- Docker Compose 代理变量传递
- 端口统一修复（开发端口 3174）

### 版本管理与在线更新
- 统一版本管理（根目录 `VERSION` 文件 + `ldflags` 构建注入）
- 后端 system 模块（版本信息 / 更新检测 / 一键更新 / 回滚 / 更新状态 / 更新历史 6 个 API）
- GitHub Releases API 集成（自动检测最新版本 + 语义版本比较）
- Updater Sidecar 容器（Docker 模式下独立更新执行器，挂载 docker.sock）
- 宿主机更新模式（直接执行 git checkout + go build）
- 前端 systemStore + systemService（Pinia 状态管理 + API 调用）
- 前端「关于 AllFi」页面（Settings 页新增区域）
- 三语言国际化支持（zh-CN / zh-TW / en-US）
- 每小时自动检查新版本
