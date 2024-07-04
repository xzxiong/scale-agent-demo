//go:build !linux
// +build !linux

package kubelet

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/cmd/kubelet/app/options"
	"k8s.io/kubernetes/pkg/kubelet/cm"
)

var NOTSupported = fmt.Errorf("NOT Supported")

func GetCgroupCpu(pod *corev1.Pod) *cm.ResourceConfig   { return nil }
func GetKubeletServer() (*options.KubeletServer, error) { return nil, NOTSupported }

const Mode = "kubelet"
