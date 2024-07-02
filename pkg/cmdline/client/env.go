package client

import (
	"context"
	"fmt"
	"os"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

const POD_NAMESPACE = "POD_NAMESPACE"
const POD_NAME = "HOSTNAME"

func GetNodeName(ctx context.Context) string {

	clientset := GetK8sClient()

	podName := os.Getenv(POD_NAME)
	posNS := os.Getenv(POD_NAMESPACE)
	fmt.Printf("pod name: %s, namespace: %s\n", podName, posNS)

	// 获取当前 Pod 的信息
	pod, err := clientset.CoreV1().Pods(posNS).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return pod.Spec.NodeName
}

var getClinetOnce sync.Once
var gClientset *kubernetes.Clientset

func GetK8sClient() *kubernetes.Clientset {
	getClinetOnce.Do(func() {
		config := ctrl.GetConfigOrDie()
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		gClientset = clientset
	})
	return gClientset
}
