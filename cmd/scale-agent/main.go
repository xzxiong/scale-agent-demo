package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	// 创建一个 HTTP 服务器
	http.HandleFunc("/", handler)

	// 启动服务器，监听指定端口
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 处理请求的逻辑
	time.Sleep(5 * time.Second)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, sidecar!"))
}
