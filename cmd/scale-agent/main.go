package main

import (
	"log"
	"net/http"
	"os"

	"github.com/xzxiong/scale-agent-demo/cmd"
)

func main() {
	if len(os.Args) > 1 {
		cmd.Execute()
	} else {
		// 创建一个 HTTP 服务器
		http.HandleFunc("/", handler)

		// 启动服务器，监听指定端口
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 处理请求的逻辑
	log.Printf("Received request: %s\n", r.Body)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, sidecar!"))
}
