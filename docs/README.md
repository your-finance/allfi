# AllFi 项目文档

> 文档已于 2026-03-31 按当前仓库实现重新审计。

---

## 优先阅读

| 文档 | 说明 |
|------|------|
| [项目现状](project-status.md) | 当前仓库状态、部署形态、测试结果、文档审计结论 |
| [部署指南](guides/deployment-guide.md) | 真实可执行的部署方式与端口说明 |
| [开发指南](guides/dev-guide.md) | 本地开发、调试、测试、构建命令 |
| [技术基线](tech/tech-baseline.md) | 当前前后端架构、模块边界、运行方式 |

> 在线文档入口：后端启动后访问 `http://localhost:8080/swagger/`，OpenAPI JSON 为 `http://localhost:8080/api.json`。

## 产品文档

| 文档 | 说明 |
|------|------|
| [业务概览](product/biz-overview.md) | 产品定位、目标用户、价值主张 |
| [功能全景](product/feature-overview.md) | 当前功能模块总览与能力边界 |

## 技术文档

| 文档 | 说明 |
|------|------|
| [技术基线](tech/tech-baseline.md) | 当前实现的技术栈、架构、部署模式 |
| [API 接口文档](tech/api-reference.md) | 按模块整理的接口说明，权威路径以 Swagger 为准 |

## 规格文档

| 文档 | 说明 |
|------|------|
| [前端规格](specs/frontend-spec.md) | 前端页面、状态管理、接口层与认证流转 |
| [后端规格](specs/backend-spec.md) | 后端模块、路由分层、数据存储与更新机制 |

## 设计文档

| 文档 | 说明 |
|------|------|
| [UI/UX 设计规范](design/ui-ux-standards.md) | 主题系统、视觉语言、组件分组 |
| [多语言系统](design/i18n.md) | 语言切换机制与翻译入口 |

## 使用指南

| 文档 | 说明 |
|------|------|
| [部署指南](guides/deployment-guide.md) | 最终用户部署、维护者部署、备份恢复 |
| [开发指南](guides/dev-guide.md) | 开发环境、测试命令、目录说明 |
| [代理配置指南](guides/proxy-guide.md) | 国内网络环境代理与镜像源配置 |
| [用户指南](guides/user-guide.md) | 首次使用、账户接入、安全建议 |

---

## 文档约定

- `product/` 关注产品价值与功能边界。
- `tech/` 和 `specs/` 以“当前代码实现”为准，不再写成理想化方案。
- `design/` 只保留长期有效的设计原则和系统说明。
- 已完成且明显过期的计划文档已从主文档区清理。
