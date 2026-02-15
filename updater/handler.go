package main

import (
	"encoding/json"
	"net/http"
)

// updateRequest 更新/回滚请求体
type updateRequest struct {
	TargetVersion string `json:"target_version"`
}

// statusResponse 状态响应
type statusResponse struct {
	State    string `json:"state"`
	Step     int    `json:"step"`
	Total    int    `json:"total"`
	StepName string `json:"step_name"`
	Message  string `json:"message"`
}

// handleUpdate 处理更新请求，异步执行 git fetch + checkout + docker-compose rebuild
func handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持 POST", http.StatusMethodNotAllowed)
		return
	}
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求体解析失败", http.StatusBadRequest)
		return
	}
	if req.TargetVersion == "" {
		http.Error(w, "target_version 不能为空", http.StatusBadRequest)
		return
	}
	// 异步执行更新
	go doUpdate(req.TargetVersion)
	json.NewEncoder(w).Encode(map[string]string{"status": "started", "message": "更新已启动"})
}

// handleRollback 处理回滚请求，异步执行 git checkout + docker-compose rebuild
func handleRollback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持 POST", http.StatusMethodNotAllowed)
		return
	}
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求体解析失败", http.StatusBadRequest)
		return
	}
	if req.TargetVersion == "" {
		http.Error(w, "target_version 不能为空", http.StatusBadRequest)
		return
	}
	// 异步执行回滚
	go doRollback(req.TargetVersion)
	json.NewEncoder(w).Encode(map[string]string{"status": "started", "message": "回滚已启动"})
}

// handleStatus 返回当前更新/回滚的进度状态
func handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "仅支持 GET", http.StatusMethodNotAllowed)
		return
	}
	json.NewEncoder(w).Encode(getStatus())
}
