package container

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	cri "k8s.io/cri-api/pkg/apis"
	remote "k8s.io/cri-client/pkg"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/features"
)

// defaultTimeout is the default timeout inherited by crictl
const defaultTimeout = 2 * time.Second

func GetContrainerdClient(ctx context.Context) (cri.RuntimeService, error) {

	// cc kubelet.PreInitRuntimeService

	var tp trace.TracerProvider
	if utilfeature.DefaultFeatureGate.Enabled(features.KubeletTracing) {
		tp = trace.NewNoopTracerProvider()
	}

	logger := klog.Background()
	if remoteRuntimeService, err := remote.NewRemoteRuntimeService(defaultContainerdAddress, defaultTimeout, tp, &logger); err != nil {
		return nil, err
	} else {
		return remoteRuntimeService, nil
	}
}
