package client

import (
	"context"
	"fmt"
	"os"
	"sync"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/xzxiong/scale-agent-demo/pkg/util"
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

func ListPodsByNodeName(ctx context.Context, nodeName string) (res []*corev1.Pod) {
	var err error

	cli := GetK8sManagerClient(ctx)
	podList := &corev1.PodList{}
	err = cli.List(ctx, podList,
		client.MatchingFields{
			fieldNodeName: nodeName,
		},
		//	&client.ListOptions{
		//		Raw: &metav1.ListOptions{
		//			FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
		//		},
		//	},
		//	// 看起来只支持 namespace
		//client.MatchingFieldsSelector{
		//	Selector: fields.OneTermEqualSelector(fieldNodeName, nodeName),
		//},
	)
	util.NoErrOrDie(err)

	for _, pod := range podList.Items {
		res = append(res, &pod)
	}
	return
}

func GetNode(ctx context.Context, nodeName string) *corev1.Node {

	clientset := GetK8sClient()
	fmt.Printf("Get node: %s\n", nodeName)
	node, err := clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	return node
}

func GetPod(ctx context.Context, ns, podName string) *corev1.Pod {
	clientset := GetK8sClient()
	pod, err := clientset.CoreV1().Pods(ns).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return pod
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

var getMgrOnce sync.Once
var mgr manager.Manager

const fieldNodeName = "spec.nodeName"

func GetK8sManagerClient(ctx context.Context) client.Client {
	var err error

	getMgrOnce.Do(func() {
		scheme := runtime.NewScheme()
		utilruntime.Must(corev1.AddToScheme(scheme))

		mgr, err = ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
			Scheme: scheme,
			Metrics: metricsserver.Options{
				BindAddress: "0", // close
			},
			PprofBindAddress:       ":8183",
			HealthProbeBindAddress: ":8182",
			LeaderElection:         false,
		})
		util.NoErrOrDie(err)

		go func() {
			err = mgr.Start(ctx)
			util.NoErrOrDie(err)
		}()

		// init self-defined. indexer
		indexer := mgr.GetFieldIndexer()
		indexer.IndexField(ctx, &corev1.Pod{}, fieldNodeName, func(o client.Object) []string {
			nodeName := o.(*corev1.Pod).Spec.NodeName
			if nodeName != "" {
				return []string{nodeName}
			}
			return nil
		})

		fmt.Println("wait ctlMgr.Elected")
		<-mgr.Elected()
		fmt.Println("wait ctlMgr.Elected: done")
	})

	return mgr.GetClient()
}
