package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/xzxiong/scale-agent-demo/pkg/containerd"
)

func main() {
	ctx := context.Background()
	// 创建一个 HTTP 服务器
	http.HandleFunc("/", handler)

	// 启动服务器，监听指定端口
	log.Fatal(http.ListenAndServe(":8080", nil))

	// cmd-tools
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 处理请求的逻辑
	time.Sleep(5 * time.Second)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, sidecar!"))
}

func commandList(ctx context.Context) {
	cs, err := containerd.GetAllContainer(ctx)
	if err != nil {
		fmt.Printf("Errro: %s", err.Error())
		os.Exit(1)
	}
	for _, c := range cs {
		c.Labels()
	}
}
