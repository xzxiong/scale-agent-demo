//go:build !linux
// +build !linux

package kubelet

import (
	"fmt"

	"k8s.io/kubernetes/cmd/kubelet/app/options"
	"k8s.io/kubernetes/pkg/kubelet"
)

func newCgroupManager(s *options.KubeletServer, kubeDeps *kubelet.Dependencies) (interface{}, error) {
	return nil, fmt.Errorf("NOT SUPPORTED")
}

func newContainerManager(s *options.KubeletServer, err error, kubeDeps *kubelet.Dependencies) error {
	return fmt.Errorf("NOT SUPPORTED")
}
