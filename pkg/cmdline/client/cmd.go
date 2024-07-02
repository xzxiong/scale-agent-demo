package client

import (
	"context"
	"fmt"
	"os"
	"sync"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

const PodNamespace = "POD_NAMESPACE"
const PodName = "HOSTNAME"

func GetNodeName(ctx context.Context) string {

	clientset := GetK8sClient()

	podName := os.Getenv(PodName)
	posNS := os.Getenv(PodNamespace)
	fmt.Printf("pod name: %s, namespace: %s\n", podName, posNS)

	// 获取当前 Pod 的信息
	pod, err := clientset.CoreV1().Pods(posNS).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return pod.Spec.NodeName
}

func ListPodsByNode(ctx context.Context, nodeName string) (res []*corev1.Pod) {

	clientset := GetK8sClient()
	fmt.Printf("list namespsce\n")
	nsList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ns := range nsList.Items {
		fmt.Printf("list pod in namespace '%s'\n", ns.Name)
		podList, err := clientset.CoreV1().Pods(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			panic(err)
		}

		for _, pod := range podList.Items {
			if pod.Spec.NodeName == nodeName {
				res = append(res, &pod)
			}
		}
	}

	return
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
