// AllFi 后端服务入口
// 全资产聚合平台 - 统一管理 CEX、区块链、传统资产
// 采用 GoFrame 框架 + 模块化分层架构
package main

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/gctx"

	// 导入全局中间件
	"your-finance/allfi/internal/middleware"

	// 导入定时任务管理器
	"your-finance/allfi/internal/cron"

	// 导入版本信息
	"your-finance/allfi/internal/version"

	// 导入所有模块的 logic 包，触发 init() 注册服务
	_ "your-finance/allfi/internal/app/achievement/logic"
	_ "your-finance/allfi/internal/app/asset/logic"
	_ "your-finance/allfi/internal/app/attribution/logic"
	_ "your-finance/allfi/internal/app/auth/logic"
	_ "your-finance/allfi/internal/app/benchmark/logic"
	_ "your-finance/allfi/internal/app/defi/logic"
	_ "your-finance/allfi/internal/app/exchange/logic"
	_ "your-finance/allfi/internal/app/exchange_rate/logic"
	_ "your-finance/allfi/internal/app/fee/logic"
	_ "your-finance/allfi/internal/app/forecast/logic"
	_ "your-finance/allfi/internal/app/goal/logic"
	_ "your-finance/allfi/internal/app/health/logic"
	_ "your-finance/allfi/internal/app/health_score/logic"
	_ "your-finance/allfi/internal/app/manual_asset/logic"
	_ "your-finance/allfi/internal/app/market/logic"
	_ "your-finance/allfi/internal/app/nft/logic"
	_ "your-finance/allfi/internal/app/notification/logic"
	_ "your-finance/allfi/internal/app/pnl/logic"
	_ "your-finance/allfi/internal/app/price_alert/logic"
	_ "your-finance/allfi/internal/app/report/logic"
	_ "your-finance/allfi/internal/app/strategy/logic"
	_ "your-finance/allfi/internal/app/system/logic"
	_ "your-finance/allfi/internal/app/transaction/logic"
	_ "your-finance/allfi/internal/app/user/logic"
	_ "your-finance/allfi/internal/app/wallet/logic"
	_ "your-finance/allfi/internal/app/webpush/logic"

	// 导入控制器注册
	achievementCtrl "your-finance/allfi/internal/app/achievement/controller"
	assetCtrl "your-finance/allfi/internal/app/asset/controller"
	attributionCtrl "your-finance/allfi/internal/app/attribution/controller"
	authCtrl "your-finance/allfi/internal/app/auth/controller"
	benchmarkCtrl "your-finance/allfi/internal/app/benchmark/controller"
	defiCtrl "your-finance/allfi/internal/app/defi/controller"
	exchangeCtrl "your-finance/allfi/internal/app/exchange/controller"
	exchangeRateCtrl "your-finance/allfi/internal/app/exchange_rate/controller"
	feeCtrl "your-finance/allfi/internal/app/fee/controller"
	forecastCtrl "your-finance/allfi/internal/app/forecast/controller"
	goalCtrl "your-finance/allfi/internal/app/goal/controller"
	healthCtrl "your-finance/allfi/internal/app/health/controller"
	healthScoreCtrl "your-finance/allfi/internal/app/health_score/controller"
	manualAssetCtrl "your-finance/allfi/internal/app/manual_asset/controller"
	marketCtrl "your-finance/allfi/internal/app/market/controller"
	nftCtrl "your-finance/allfi/internal/app/nft/controller"
	notificationCtrl "your-finance/allfi/internal/app/notification/controller"
	pnlCtrl "your-finance/allfi/internal/app/pnl/controller"
	priceAlertCtrl "your-finance/allfi/internal/app/price_alert/controller"
	reportCtrl "your-finance/allfi/internal/app/report/controller"
	strategyCtrl "your-finance/allfi/internal/app/strategy/controller"
	systemCtrl "your-finance/allfi/internal/app/system/controller"
	transactionCtrl "your-finance/allfi/internal/app/transaction/controller"
	userCtrl "your-finance/allfi/internal/app/user/controller"
	walletCtrl "your-finance/allfi/internal/app/wallet/controller"
	webpushCtrl "your-finance/allfi/internal/app/webpush/controller"

	// 导入数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
)

func main() {
	// 打印启动横幅
	printBanner()

	ctx := gctx.New()
	s := g.Server()

	// 配置 OpenAPI 文档元数据和 Swagger UI
	enhanceOpenApi(s)

	// 注册全局中间件（CORS → Context → Logger → ErrorHandler → ResponseWrapper）
	middleware.Register(s)

	// 注册路由
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		// ===== 免认证路由 =====
		healthCtrl.Register(group)
		authCtrl.Register(group)

		// ===== 需认证路由 =====
		group.Middleware(middleware.Auth)

		// 核心资产模块
		assetCtrl.Register(group)
		exchangeCtrl.Register(group)
		walletCtrl.Register(group)
		manualAssetCtrl.Register(group)
		exchangeRateCtrl.Register(group)

		// 分析模块
		pnlCtrl.Register(group)
		attributionCtrl.Register(group)
		forecastCtrl.Register(group)
		feeCtrl.Register(group)
		benchmarkCtrl.Register(group)
		healthScoreCtrl.Register(group)

		// 市场数据
		marketCtrl.Register(group)

		// DeFi / NFT
		defiCtrl.Register(group)
		nftCtrl.Register(group)

		// 交易记录
		transactionCtrl.Register(group)

		// 通知系统
		notificationCtrl.Register(group)
		priceAlertCtrl.Register(group)
		webpushCtrl.Register(group)

		// 报告
		reportCtrl.Register(group)

		// 策略和目标
		strategyCtrl.Register(group)
		goalCtrl.Register(group)

		// 用户设置
		userCtrl.Register(group)

		// 成就系统
		achievementCtrl.Register(group)

		// 系统管理
		systemCtrl.Register(group)
	})

	// 启动定时任务（快照/报告/通知/价格预警/策略检查/风险提醒）
	cronManager := cron.NewCronManager()
	cronManager.Start()
	defer cronManager.Stop()

	g.Log().Infof(ctx, "AllFi v%s 启动成功", version.Version)
	s.Run()
}

// enhanceOpenApi 增强 OpenAPI 文档元数据和 Swagger UI 配置
// GoFrame 内置 Swagger 支持，通过 config.yaml 的 openapiPath/swaggerPath 启用
// 此函数补充 Info、SecuritySchemes、Servers 等元数据
func enhanceOpenApi(s *ghttp.Server) {
	openapi := s.GetOpenApi()

	// 设置 API 基本信息
	openapi.Info = goai.Info{
		Title:       "AllFi API",
		Description: "AllFi 全资产聚合平台 API — 统一管理 CEX、区块链、DeFi、NFT、传统资产",
		Version:     version.Version,
	}

	// 设置服务器地址
	openapi.Servers = &goai.Servers{
		{URL: "http://localhost:8080", Description: "本地开发环境"},
	}

	// 设置 JWT Bearer Token 认证方案
	openapi.Components.SecuritySchemes = goai.SecuritySchemes{
		"BearerAuth": goai.SecuritySchemeRef{
			Value: &goai.SecurityScheme{
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
				Description:  "JWT Token 认证，通过 /api/v1/auth/login 获取",
			},
		},
	}

	// 设置全局安全要求（所有接口默认需要认证，免认证接口在 g.Meta 中单独覆盖）
	openapi.Security = &goai.SecurityRequirements{
		goai.SecurityRequirement{"BearerAuth": {}},
	}

	// 使用 SwaggerUI 替代默认的 Redoc（GoFrame v2.10 默认 UI 是 Redoc）
	s.SetSwaggerUITemplate(swaggerUITemplate)
}

// swaggerUITemplate SwaggerUI 自定义模板
// {SwaggerUIDocUrl} 占位符由 GoFrame 自动替换为 openapiPath 路径
const swaggerUITemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>AllFi API 文档</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui.css">
    <style>
        html { box-sizing: border-box; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
        .topbar { display: none; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: '{SwaggerUIDocUrl}',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                defaultModelsExpandDepth: 1,
                docExpansion: "list",
                filter: true,
                tryItOutEnabled: true,
                persistAuthorization: true
            });
        };
    </script>
</body>
</html>`

// printBanner 打印启动横幅
func printBanner() {
	banner := `
    ___    ____   __    ______ _
   /   |  / / /  / /   / ____/(_)
  / /| | / / /  / /   / /_   / /
 / ___ |/ / /  / /   / __/  / /
/_/  |_/_/_/  /_/   /_/    /_/

AllFi - 全资产聚合平台 v%s
GoFrame 模块化架构 | 自托管 | 完全掌控数据
`
	fmt.Printf(banner, version.Version)
	fmt.Println()
}
