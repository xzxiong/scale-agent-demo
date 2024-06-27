package kubelet

import (
	"fmt"
	"os"

	utilfeature "k8s.io/apiserver/pkg/util/feature"
	klog "k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kubelet/app"
	"k8s.io/kubernetes/cmd/kubelet/app/options"
)

func init() {

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

	//	kubeletDeps, err := app.UnsecuredDependencies(kubeletServer, utilfeature.DefaultFeatureGate)
	_, err = app.UnsecuredDependencies(kubeletServer, utilfeature.DefaultFeatureGate)
	if err != nil {
		panic(fmt.Errorf("failed to construct kubelet dependencies: %w", err))
	}
}
