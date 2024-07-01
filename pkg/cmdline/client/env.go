package kubelet

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GenEnv(key string) string {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 获取当前 Pod 的信息
	pod, err := clientset.CoreV1().Get("your-pod-name", rest.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	val := os.environ.get("NODE_NAME")

}
