//go:build !linux
// +build !linux

package kubelet

func buildContainerMgr() error {
	return nil
}

func buildCgroupMgr() error {
	return nil
}

const Mode = "kubelet"
