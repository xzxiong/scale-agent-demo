//go:build !linux
// +build !linux

package kubelet

import (
	"fmt"

	"k8s.io/kubernetes/cmd/kubelet/app/options"
)

var NOTSupported = fmt.Errorf("NOT Supported")

func buildContainerMgr() error { return NOTSupported }
func buildCgroupMgr() error    { return NOTSupported }
func GetKubeletServer() (*options.KubeletServer, error) {
	return nil, NOTSupported
}

const Mode = "kubelet"
