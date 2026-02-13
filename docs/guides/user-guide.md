# AllFi 用户指南

> **版本**：v2.0 | **更新时间**：2026-02-11

---

## 目录

1. [快速开始](#1-快速开始)
2. [安全配置 API 密钥](#2-安全配置-api-密钥)
3. [功能使用指南](#3-功能使用指南)
4. [常见问题](#4-常见问题)
5. [故障排查](#5-故障排查)

---

## 1. 快速开始

### 1.1 前置要求

**推荐方式：使用 Docker（零配置）**
- Docker 和 Docker Compose
- 512 MB+ 可用内存
- 1 GB 可用磁盘空间

**或者：手动编译运行**
- Go 1.24
- Node.js 18+
- pnpm 9+

### 1.2 部署步骤（Docker 推荐）

```bash
# 下载代码
git clone https://github.com/your-finance/allfi.git
cd allfi

# 配置环境变量
cp .env.example .env
# 编辑 .env，设置 ALLFI_MASTER_KEY（用 openssl rand -base64 32 生成）

# 启动服务
docker-compose up -d
```

访问 http://localhost:3174 即可使用。

> 详细部署说明见 [部署指南](./deployment-guide.md)。

### 1.3 首次使用流程

1. **设置 PIN 码**
   - 首次访问时，系统会引导你设置 4-8 位数字 PIN 码
   - PIN 码用于登录认证，经 bcrypt 加密存储
   - 设置完成后自动获取 JWT Token 并登录

2. **新手引导**
   - 登录后进入引导向导（Onboarding Wizard）
   - 按步骤添加第一个交易所账户、钱包地址或手动资产
   - 完成后自动跳转到仪表盘

3. **查看资产**
   - 仪表盘展示资产总览、分布饼图、趋势图等
   - 等待数据同步完成（通常 5-10 秒）

---

## 2. 安全配置 API 密钥

> AllFi 是自托管应用，API 密钥经 AES-256-GCM 加密后存储在本地 SQLite 数据库，永远不会上传到任何服务器。

### 2.1 获取 API 密钥

#### Binance

1. 登录 Binance → [API 管理页面](https://www.binance.com/zh-CN/my/settings/api-management)
2. 创建 API → 输入名称（如 "AllFi 只读"）→ 完成安全验证
3. **仅勾选「读取」权限**

#### OKX

1. 登录 OKX → [API 申请页面](https://www.okx.com/account/my-api)
2. 创建 API → 输入名称和 Passphrase → 完成安全验证
3. **仅勾选「读取」权限**

#### Coinbase

1. 登录 Coinbase → Settings → API
2. 创建 New API Key → **仅勾选 View 权限**

### 2.2 最小权限原则

**禁止授权的权限**：
- 提现（Withdrawal）
- 交易（Trading / Spot Trading）
- 合约（Enable Futures）
- 资金划转（Transfer）

**只授权读取权限**：
- 读取账户信息（Read Info）
- 查询余额（Query Balance）

> 即使 API Key 泄露，黑客也**无法转走任何资金**。

### 2.3 IP 白名单（强烈推荐）

如果 AllFi 部署在固定 IP 的服务器上，建议在交易所 API 设置中开启 IP 白名单。

### 2.4 在 AllFi 中添加 API Key

1. 进入「账户管理」页面
2. 点击「添加账户」→ 选择交易所（Binance / OKX / Coinbase）
3. 填写 API Key、API Secret（OKX 还需 Passphrase）
4. 点击「测试连接」验证
5. 点击「添加账户」保存

---

## 3. 功能使用指南

### 3.1 仪表盘（Dashboard）

仪表盘是 AllFi 的主页面，一屏展示所有关键数据：

- **总资产卡片**：总价值 + 24h 涨跌幅（支持 USDC/BTC/ETH/CNY 计价切换）
- **资产分布饼图**：CEX / 链上 / 手动资产占比
- **资产趋势图**：7 天 / 30 天 / 90 天 / 1 年折线图
- **Top 10 资产列表**：按价值降序排列
- **健康评分卡片**：多维度资产健康评估
- **目标追踪进度条**：投资目标完成度
- **基准对比**：与 BTC/ETH/S&P500 收益率对比

仪表盘支持自定义布局（DashboardCustomizer），可拖拽调整 Widget 排列。

### 3.2 账户管理（Accounts）

管理三类资产来源：

#### 交易所账户（CEX）

- 支持 **Binance**、**OKX**、**Coinbase**
- 添加、编辑、删除、测试连接
- 查看各账户持仓明细

#### 区块链钱包

- 支持 **Ethereum**、**BSC**、**Polygon**（Etherscan 统一客户端还支持 Arbitrum/Optimism/Base）
- 支持批量导入钱包地址
- 自动获取 ETH + ERC20/BEP20 代币余额

#### 手动资产（传统资产）

适用于银行存款、现金、股票、基金等非加密资产：

| 资产类型 | 说明 |
|---------|------|
| `bank` | 银行存款 |
| `cash` | 现金 |
| `stock` | 股票 |
| `fund` | 基金 |

手动资产需要定期手动更新金额。

### 3.3 历史记录（History）

- **资产趋势折线图**：多时间范围查看资产变化
- **日历热力图**：GitHub 风格，绿色=盈利日，红色=亏损日
- **快照列表**：每日资产快照数据（后端每小时自动快照）
- 支持导出 CSV

### 3.4 数据分析（Analytics）

- **每日盈亏（PnL）**：逐日收益曲线
- **盈亏摘要**：今日 / 本周 / 本月收益
- **资产归因**：分析价格变动 vs 数量变动对收益的贡献
- **趋势预测**：基于历史数据预测达到目标金额的时间
- **费用分析**：交易手续费 + Gas 费 + 提现费统计

### 3.5 资产报告（Reports）

- **每日/每周报告**：自动生成，包含总资产、涨跌、基准对比
- **月度报告**：配置变化、费用摘要、建议
- **年度报告**：年度总结，支持分享

### 3.6 隐私模式

**使用场景**：屏幕共享、直播、公共场合

**开启方式**：点击顶部导航栏的眼睛图标

**效果**：
```
隐私模式 OFF：总资产：$123,456.78 USDC
隐私模式 ON：总资产：**** USDC（百分比仍显示）
```

### 3.7 命令面板（Cmd+K）

快捷键 `Cmd+K`（Mac）/ `Ctrl+K`（Windows）打开命令面板，可快速：
- 跳转到任意页面
- 搜索资产
- 切换主题
- 切换语言
- 刷新数据

### 3.8 多语言支持

支持三种语言，在「设置」页面切换：
- 简体中文（默认）
- 繁体中文
- English

### 3.9 主题切换

内置 4 套主题：
1. **Nexus Pro**（默认深色）
2. **Vestia**（GitHub 风深色）
3. **XChange**（交易所风深色）
4. **Aurora**（浅色）

### 3.10 设置（Settings）

| 设置项 | 说明 |
|-------|------|
| 默认计价币种 | USDC / USD / BTC / ETH / CNY |
| 语言 | 简体中文 / 繁体中文 / English |
| 主题 | 4 套可选 |
| 通知偏好 | 每日摘要、价格预警、资产变动提醒 |
| WebPush 推送 | 浏览器推送通知 |
| 数据导出 | 导出 JSON / CSV |
| 清除缓存 | 清除价格和汇率缓存 |
| 修改 PIN | 更改登录 PIN 码 |

---

## 4. 常见问题

### Q1: 我的 API Key 会上传到云端吗？

**不会。** AllFi 是纯本地运行的自托管软件。API Key 经 AES-256-GCM 加密后存储在本地数据库 `allfi.db` 中，永远不会上传到任何服务器。

### Q2: 如何备份数据？

只需备份 SQLite 数据库文件：

```bash
# Docker 部署
docker cp allfi-backend:/app/data/allfi.db ./backup/

# 手动部署
cp core/data/allfi.db ./backup/
```

详见 [部署指南 - 数据备份](./deployment-guide.md#数据备份与恢复)。

### Q3: 为什么显示的金额和交易所不完全一致？

可能原因：
1. **汇率差异**：AllFi 使用 CoinGecko/Yahoo Finance 汇率，与交易所内部汇率略有不同
2. **缓存延迟**：价格缓存默认 5 分钟
3. **计价单位**：检查是否使用相同的计价单位

**解决**：点击「刷新」按钮手动刷新数据。

### Q4: 支持哪些交易所和区块链？

**交易所**：Binance、OKX、Coinbase

**区块链**：Ethereum、BSC、Polygon（Etherscan 统一客户端还支持 Arbitrum/Optimism/Base 查询）

### Q5: 忘记 PIN 码怎么办？

当前版本需要重置数据库来重新设置 PIN。备份数据后删除 `allfi.db` 文件，重启服务即可重新设置。

### Q6: 数据多久更新一次？

- **资产快照**：后端每小时自动创建一次
- **价格缓存**：5 分钟 TTL
- **手动刷新**：随时可点击刷新按钮

### Q7: 如何更改计价币种？

点击仪表盘顶部的币种按钮（USDC / BTC / ETH / CNY），或在「设置」中修改默认计价币种。

---

## 5. 故障排查

### 5.1 无法访问前端页面

```bash
# 检查容器状态
docker-compose ps

# 检查端口占用
lsof -i :3174
lsof -i :8080

# 重启服务
docker-compose restart
```

### 5.2 无法获取交易所余额

检查清单：
1. API Key 是否正确
2. 是否仅授予了读取权限
3. IP 白名单是否正确
4. 网络是否正常

查看日志：

```bash
docker logs allfi-backend | grep ERROR
```

### 5.3 链上数据为空

原因：缺少 Etherscan/BscScan API 密钥。

解决：在 `.env` 中配置 `ETHERSCAN_API_KEY` 和 `BSCSCAN_API_KEY`，然后重启服务。

### 5.4 数据库损坏

```bash
# 检查数据库完整性
sqlite3 core/data/allfi.db "PRAGMA integrity_check;"

# 如有备份，恢复
cp backup/allfi-latest.db core/data/allfi.db
```

### 5.5 获取帮助

- [GitHub Issues](https://github.com/your-finance/allfi/issues)
- [GitHub Discussions](https://github.com/your-finance/allfi/discussions)

提交 Issue 时请提供：AllFi 版本号、操作系统、完整错误日志、重现步骤。

---

## 相关文档

- [功能全景](../product/feature-overview.md)
- [API 接口文档](../tech/api-reference.md)
- [部署指南](./deployment-guide.md)
- [开发指南](./dev-guide.md)

---

**文档维护者**: @allfi
**最后更新**: 2026-02-11
