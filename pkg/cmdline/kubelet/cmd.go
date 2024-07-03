//go:build linux
// +build linux

package kubelet

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	libcontainercgroups "github.com/opencontainers/runc/libcontainer/cgroups"
	corev1 "k8s.io/api/core/v1"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kubelet/app"
	"k8s.io/kubernetes/cmd/kubelet/app/options"
	"k8s.io/kubernetes/pkg/kubelet"
	"k8s.io/kubernetes/pkg/kubelet/cm"
)

func buildContainerMgr() (*kubelet.Dependencies, error) {

	kubeletFlags := options.NewKubeletFlags()
	kubeletConfig, err := options.NewKubeletConfiguration()
	if err != nil {
		klog.ErrorS(err, "Failed to create a new kubelet configuration")
		os.Exit(1)
	}

	// construct a KubeletServer from kubeletFlags and kubeletConfig
	kubeletServer := &options.KubeletServer{
		KubeletFlags:         *kubeletFlags,
		KubeletConfiguration: *kubeletConfig, // This is the KEY config, MUST load from config.
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
		return nil, err2
	}
	// err = containerManager.Start()
	klog.Infof("NewContainerManager start: %v", kubeDeps.ContainerManager)

	return kubeDeps, nil
}

func setCgroupCpu(pod *corev1.Pod) error {

	// ref k8s.io/kubernetes@v1.28.4/pkg/kubelet/cm/cgroup_manager_linux.go

	// =======
	// ref cgm.SetCgroupConfig
	// =======
	// cgm, err := buildCgroupMgr()
	// if err != nil {
	// 	panic(err)
	// }

	kubeletServer, err := getKubeletServer()
	if err != nil {
		panic(err)
	}

	kubeDeps, err := buildContainerMgr()
	if err != nil {
		panic(err)
	}

	subSystems, err := cm.GetCgroupSubsystems()
	if err != nil {
		panic(err)
	}

	cmgr := kubeDeps.ContainerManager
	pcm := cmgr.NewPodContainerManager()
	//resCfg, err := pcm.GetPodCgroupConfig(nil, corev1.ResourceStorage)
	podCgroupName, _ := pcm.GetPodContainerName(pod)

	// map: subsystem -> cgroup path
	cgroupPaths := buildCgroupPaths(podCgroupName, kubeletServer.CgroupDriver, subSystems)
	cpuCgroupPath := cgroupPaths[CgroupControllerCpu]

	// set CPU
	cpuPeriod := uint64(CPUPeriodUs)
	cpuMax := int64(CPUPeriodUs * 1.0)
	cpuShared := cm.MilliCPUToShares(cpuMax)
	err = setCgroupv2CpuConfig(cpuCgroupPath, &cm.ResourceConfig{
		CPUPeriod: &cpuPeriod,
		CPUQuota:  &cpuMax,
		CPUShares: &cpuShared,
	})

	return err
}

const CPUPeriodUs = 100000

const CgroupDriverSystemd = "systemd"

const CgroupControllerCpu = string(corev1.ResourceCPU)
const CgroupControllerMemory = string(corev1.ResourceMemory)
const CgroupControllerStorage = string(corev1.ResourceStorage) // this for Volume NOT for cgroup

// buildCgroupPaths ref k8s.io/kubernetes@v1.28.4/pkg/kubelet/cm/cgroup_manager_linux.go
func buildCgroupPaths(name cm.CgroupName, cgroupDriver string, subsystems *cm.CgroupSubsystems) map[string]string {
	// fixme: check
	cgroupFsAdaptedName := name.ToCgroupfs()
	if cgroupDriver == CgroupDriverSystemd {
		cgroupFsAdaptedName = name.ToSystemd()
	}
	cgroupPaths := make(map[string]string, len(subsystems.MountPoints))
	for key, val := range subsystems.MountPoints {
		cgroupPaths[key] = path.Join(val, cgroupFsAdaptedName)
	}
	return cgroupPaths
}

func buildCgroupMgr() (cm.CgroupManager, error) {
	kubeletServer, err := getKubeletServer()
	if err != nil {
		panic(err)
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

	mgr, err := newCgroupManager(s, kubeDeps)
	if err != nil {
		panic(err)
	}
	return mgr, err
}

// setCgroupCpuConfig ref k8s.io/kubernetes@v1.28.4/pkg/kubelet/cm/cgroup_manager_linux.go
func setCgroupCpuConfig(cgroupPath string, resourceConfig *cm.ResourceConfig) error {
	if libcontainercgroups.IsCgroup2UnifiedMode() {
		return setCgroupv2CpuConfig(cgroupPath, resourceConfig)
	} else {
		// return setCgroupv1CpuConfig(cgroupPath, resourceConfig)
		return fmt.Errorf("NOT IMPLEMENTED")
	}
}

// setCgroupv2CpuConfig ref k8s.io/kubernetes@v1.28.4/pkg/kubelet/cm/cgroup_manager_linux.go
func setCgroupv2CpuConfig(cgroupPath string, resourceConfig *cm.ResourceConfig) error {
	if resourceConfig.CPUQuota != nil {
		if resourceConfig.CPUPeriod == nil {
			return fmt.Errorf("CpuPeriod must be specified in order to set CpuLimit")
		}
		cpuLimitStr := cm.Cgroup2MaxCpuLimit
		if *resourceConfig.CPUQuota > -1 {
			cpuLimitStr = strconv.FormatInt(*resourceConfig.CPUQuota, 10)
		}
		cpuPeriodStr := strconv.FormatUint(*resourceConfig.CPUPeriod, 10)
		cpuMaxStr := fmt.Sprintf("%s %s", cpuLimitStr, cpuPeriodStr)
		if err := os.WriteFile(filepath.Join(cgroupPath, "cpu.max"), []byte(cpuMaxStr), 0700); err != nil {
			return fmt.Errorf("failed to write %v to %v: %v", cpuMaxStr, cgroupPath, err)
		}
	}
	if resourceConfig.CPUShares != nil {
		cpuWeight := cm.CpuSharesToCpuWeight(*resourceConfig.CPUShares)
		cpuWeightStr := strconv.FormatUint(cpuWeight, 10)
		if err := os.WriteFile(filepath.Join(cgroupPath, "cpu.weight"), []byte(cpuWeightStr), 0700); err != nil {
			return fmt.Errorf("failed to write %v to %v: %v", cpuWeightStr, cgroupPath, err)
		}
	}
	return nil
}

func getKubeletServer() (*options.KubeletServer, error) {
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
	return kubeletServer, err
}

const Mode = "kubelet"
