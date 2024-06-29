package kubelet

import (
	"fmt"
	"os"

	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kubelet/app"
	"k8s.io/kubernetes/cmd/kubelet/app/options"
)

func buildContainerMgr() error {

	kubeletFlags := options.NewKubeletFlags()
	kubeletConfig, err := options.NewKubeletConfiguration()
	if err != nil {
		klog.ErrorS(err, "Failed to create a new kubelet configuration")
		os.Exit(1)
	}

	// construct a KubeletServer from kubeletFlags and kubeletConfig
	kubeletServer := &options.KubeletServer{
		KubeletFlags:         *kubeletFlags,
		KubeletConfiguration: *kubeletConfig,
	}

	kubeletDeps, err := app.UnsecuredDependencies(kubeletServer, utilfeature.DefaultFeatureGate)
	_, err = app.UnsecuredDependencies(kubeletServer, utilfeature.DefaultFeatureGate)
	if err != nil {
		panic(fmt.Errorf("failed to construct kubelet dependencies: %w", err))
	}
	fmt.Printf("kubeletServer: %s\n", kubeletServer)
	fmt.Printf("kubeletDeps is key for container_manager & cgroup manager: %v\n", kubeletDeps)

	kubeDeps := kubeletDeps
	s := kubeletServer

	err2 := newContainerManager(s, err, kubeDeps)
	if err2 != nil {
		return err2
	}
	// err = containerManager.Start()
	klog.Infof("NewContainerManager start: %v", kubeDeps.ContainerManager)
	return nil
}

func buildCgroupMgr() error {
	kubeletFlags := options.NewKubeletFlags()
	kubeletConfig, err := options.NewKubeletConfiguration()
	if err != nil {
		klog.ErrorS(err, "Failed to create a new kubelet configuration")
		os.Exit(1)
	}

	// construct a KubeletServer from kubeletFlags and kubeletConfig
	kubeletServer := &options.KubeletServer{
		KubeletFlags:         *kubeletFlags,
		KubeletConfiguration: *kubeletConfig,
	}

	kubeletDeps, err := app.UnsecuredDependencies(kubeletServer, utilfeature.DefaultFeatureGate)
	_, err = app.UnsecuredDependencies(kubeletServer, utilfeature.DefaultFeatureGate)
	if err != nil {
		panic(fmt.Errorf("failed to construct kubelet dependencies: %w", err))
	}
	fmt.Printf("kubeletServer: %s\n", kubeletServer)
	fmt.Printf("kubeletDeps is key for container_manager & cgroup manager: %v\n", kubeletDeps)

	kubeDeps := kubeletDeps
	s := kubeletServer

	_, err = newCgroupManager(s, kubeDeps)
	return err
}

const Mode = "kubelet"
