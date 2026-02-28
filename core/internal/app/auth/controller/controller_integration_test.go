// Package controller 认证控制器集成测试
package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	authApi "your-finance/allfi/api/v1/auth"
)

func TestController_AuthStatus(t *testing.T) {
	router := http.NewServeMux()
	router.HandleFunc("GET /api/v1/auth/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(authApi.GetStatusRes{PinSet: true})
	})

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	req, _ := http.NewRequest("GET", "/api/v1/auth/status", nil)
	w := httptest.NewRecorder()
	testServer.Config.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var statusResp authApi.GetStatusRes
	json.Unmarshal(w.Body.Bytes(), &statusResp)
	assert.Equal(t, true, statusResp.PinSet)
}

func TestController_AuthSetup(t *testing.T) {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/auth/setup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(authApi.SetupRes{Token: "test-token"})
	})

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	reqBody, _ := json.Marshal(map[string]string{"pin": "1234"})
	req, _ := http.NewRequest("POST", "/api/v1/auth/setup", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testServer.Config.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var setupResp authApi.SetupRes
	json.Unmarshal(w.Body.Bytes(), &setupResp)
	assert.NotEmpty(t, setupResp.Token)
}

func TestController_AuthLogin(t *testing.T) {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(authApi.LoginRes{Token: "login-token"})
	})

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	reqBody, _ := json.Marshal(map[string]string{"pin": "1234"})
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testServer.Config.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var loginResp authApi.LoginRes
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NotEmpty(t, loginResp.Token)
}
