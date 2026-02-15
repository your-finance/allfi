package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("AllFi Updater Sidecar 启动，监听 :8081")

	http.HandleFunc("/update", handleUpdate)
	http.HandleFunc("/rollback", handleRollback)
	http.HandleFunc("/status", handleStatus)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
