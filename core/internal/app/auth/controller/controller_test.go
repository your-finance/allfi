package controller_test

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"

	"fmt"

	"your-finance/allfi/internal/app/auth/controller"
	_ "your-finance/allfi/internal/app/auth/logic" // Register service
	"your-finance/allfi/internal/database"
)

func init() {
	ctx := context.Background()
	if err := database.Initialize(ctx); err != nil {
		g.Log().Fatalf(ctx, "Failed to initialize test DB: %v", err)
	}
}

func TestController_AuthFlow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()

		// Clear DB for a fresh test
		_, _ = g.DB().Exec(ctx, "DELETE FROM system_config;")

		s := g.Server(g.Map{"LogPath": ""})
		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			group.Middleware(ghttp.MiddlewareHandlerResponse)
			controller.Register(group)
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond) // wait for server to start

		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// 1. Get Status (Unset)
		resp, err := client.Get(ctx, "/api/v1/auth/status")
		t.AssertNil(err)
		t.Assert(resp.StatusCode, 200)
		t.Assert(resp.ReadAllString(), `{"code":0,"message":"OK","data":{"pin_set":false,"two_fa_enabled":false,"password_type":"pin"}}`)
		resp.Close()

		// 2. Setup PIN
		resp, err = client.Post(ctx, "/api/v1/auth/setup", g.Map{"pin": "123456"})
		t.AssertNil(err)
		t.Assert(resp.StatusCode, 200)
		// Should return token
		content := resp.ReadAllString()
		t.AssertNE(content, "")
		resp.Close()

		// 3. Get Status (Set)
		resp, err = client.Get(ctx, "/api/v1/auth/status")
		t.AssertNil(err)
		t.Assert(resp.StatusCode, 200)
		t.Assert(resp.ReadAllString(), `{"code":0,"message":"OK","data":{"pin_set":true,"two_fa_enabled":false,"password_type":"pin"}}`)
		resp.Close()

		// 4. Login PIN
		resp, err = client.Post(ctx, "/api/v1/auth/login", g.Map{"pin": "123456"})
		t.AssertNil(err)
		t.Assert(resp.StatusCode, 200)
		content = resp.ReadAllString()
		t.AssertNE(content, "")
		resp.Close()

		// 5. Login wrong PIN
		resp, err = client.Post(ctx, "/api/v1/auth/login", g.Map{"pin": "654321"})
		t.AssertNil(err)
		t.Assert(resp.StatusCode, 200) // GoFrame returns 200 by default for business errors if not configured otherwise
		// Check that it contains error message
		content = resp.ReadAllString()
		t.AssertIN("PIN 错误", content)
		resp.Close()
	})
}
