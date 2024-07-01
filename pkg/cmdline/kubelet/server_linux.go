//go:build linux
// +build linux

package kubelet

import (
	"fmt"
	"strconv"

	cadvisorapi "github.com/google/cadvisor/info/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/sets"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kubelet/app/options"
	"k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/pkg/kubelet"
	"k8s.io/kubernetes/pkg/kubelet/cm"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager/topology"
	"k8s.io/kubernetes/pkg/kubelet/eviction"
	evictionapi "k8s.io/kubernetes/pkg/kubelet/eviction/api"
	"k8s.io/kubernetes/pkg/kubelet/stats/pidlimit"
	"k8s.io/utils/cpuset"
)

// newCgroupManager returns a CgroupManager based on the passed options.
func newCgroupManager(s *options.KubeletServer, kubeDeps *kubelet.Dependencies) (interface{}, error) {

	if s.CgroupsPerQOS && s.CgroupRoot == "" {
		klog.InfoS("--cgroups-per-qos enabled, but --cgroup-root was not specified.  defaulting to /")
		s.CgroupRoot = "/"
	}

	machineInfo, err := kubeDeps.CAdvisorInterface.MachineInfo()
	if err != nil {
		return nil, err
	}
	reservedSystemCPUs, err := getReservedCPUs(machineInfo, s.ReservedSystemCPUs)
	if err != nil {
		return nil, err
	}
	if reservedSystemCPUs.Size() > 0 {
		// at cmd option validation phase it is tested either --system-reserved-cgroup or --kube-reserved-cgroup is specified, so overwrite should be ok
		klog.InfoS("Option --reserved-cpus is specified, it will overwrite the cpu setting in KubeReserved and SystemReserved", "kubeReservedCPUs", s.KubeReserved, "systemReservedCPUs", s.SystemReserved)
		if s.KubeReserved != nil {
			delete(s.KubeReserved, "cpu")
		}
		if s.SystemReserved == nil {
			s.SystemReserved = make(map[string]string)
		}
		s.SystemReserved["cpu"] = strconv.Itoa(reservedSystemCPUs.Size())
		klog.InfoS("After cpu setting is overwritten", "kubeReservedCPUs", s.KubeReserved, "systemReservedCPUs", s.SystemReserved)
	}

	kubeReserved, err := parseResourceList(s.KubeReserved)
	if err != nil {
		return nil, fmt.Errorf("--kube-reserved value failed to parse: %w", err)
	}
	systemReserved, err := parseResourceList(s.SystemReserved)
	if err != nil {
		return nil, fmt.Errorf("--system-reserved value failed to parse: %w", err)
	}
	var hardEvictionThresholds []evictionapi.Threshold
	// If the user requested to ignore eviction thresholds, then do not set valid values for hardEvictionThresholds here.
	if !s.ExperimentalNodeAllocatableIgnoreEvictionThreshold {
		hardEvictionThresholds, err = eviction.ParseThresholdConfig([]string{}, s.EvictionHard, nil, nil, nil)
		if err != nil {
			return nil, err
		}
	}
	experimentalQOSReserved, err := cm.ParseQOSReserved(s.QOSReserved)
	if err != nil {
		return nil, fmt.Errorf("--qos-reserved value failed to parse: %w", err)
	}

	var cpuManagerPolicyOptions map[string]string
	if utilfeature.DefaultFeatureGate.Enabled(features.CPUManagerPolicyOptions) {
		cpuManagerPolicyOptions = s.CPUManagerPolicyOptions
	} else if s.CPUManagerPolicyOptions != nil {
		return nil, fmt.Errorf("CPU Manager policy options %v require feature gates %q, %q enabled",
			s.CPUManagerPolicyOptions, features.CPUManager, features.CPUManagerPolicyOptions)
	}

	var topologyManagerPolicyOptions map[string]string
	if utilfeature.DefaultFeatureGate.Enabled(features.TopologyManagerPolicyOptions) {
		topologyManagerPolicyOptions = s.TopologyManagerPolicyOptions
	} else if s.TopologyManagerPolicyOptions != nil {
		return nil, fmt.Errorf("topology manager policy options %v require feature gates %q enabled",
			s.TopologyManagerPolicyOptions, features.TopologyManagerPolicyOptions)
	}

	nodeConfig := cm.NodeConfig{
		RuntimeCgroupsName:    s.RuntimeCgroups,
		SystemCgroupsName:     s.SystemCgroups,
		KubeletCgroupsName:    s.KubeletCgroups,
		KubeletOOMScoreAdj:    s.OOMScoreAdj,
		CgroupsPerQOS:         s.CgroupsPerQOS,
		CgroupRoot:            s.CgroupRoot,
		CgroupDriver:          s.CgroupDriver,
		KubeletRootDir:        s.RootDirectory,
		ProtectKernelDefaults: s.ProtectKernelDefaults,
		NodeAllocatableConfig: cm.NodeAllocatableConfig{
			KubeReservedCgroupName:   s.KubeReservedCgroup,
			SystemReservedCgroupName: s.SystemReservedCgroup,
			EnforceNodeAllocatable:   sets.NewString(s.EnforceNodeAllocatable...),
			KubeReserved:             kubeReserved,
			SystemReserved:           systemReserved,
			ReservedSystemCPUs:       reservedSystemCPUs,
			HardEvictionThresholds:   hardEvictionThresholds,
		},
		QOSReserved:                             *experimentalQOSReserved,
		CPUManagerPolicy:                        s.CPUManagerPolicy,
		CPUManagerPolicyOptions:                 cpuManagerPolicyOptions,
		CPUManagerReconcilePeriod:               s.CPUManagerReconcilePeriod.Duration,
		ExperimentalMemoryManagerPolicy:         s.MemoryManagerPolicy,
		ExperimentalMemoryManagerReservedMemory: s.ReservedMemory,
		PodPidsLimit:                            s.PodPidsLimit,
		EnforceCPULimits:                        s.CPUCFSQuota,
		CPUCFSQuotaPeriod:                       s.CPUCFSQuotaPeriod.Duration,
		TopologyManagerPolicy:                   s.TopologyManagerPolicy,
		TopologyManagerScope:                    s.TopologyManagerScope,
		TopologyManagerPolicyOptions:            topologyManagerPolicyOptions,
	}
	// END> copy from NewContainerManager

	subsystems, err := cm.GetCgroupSubsystems()
	if err != nil {
		return nil, fmt.Errorf("failed to get mounted cgroup subsystems: %v", err)
	}

	// Turn CgroupRoot from a string (in cgroupfs path format) to internal CgroupName
	cgroupRoot := cm.ParseCgroupfsToCgroupName(nodeConfig.CgroupRoot)
	cgroupManager := cm.NewCgroupManager(subsystems, nodeConfig.CgroupDriver)
	// Check if Cgroup-root actually exists on the node
	if nodeConfig.CgroupsPerQOS {
		// this does default to / when enabled, but this tests against regressions.
		if nodeConfig.CgroupRoot == "" {
			return nil, fmt.Errorf("invalid configuration: cgroups-per-qos was specified and cgroup-root was not specified. To enable the QoS cgroup hierarchy you need to specify a valid cgroup-root")
		}

		// we need to check that the cgroup root actually exists for each subsystem
		// of note, we always use the cgroupfs driver when performing this check since
		// the input is provided in that format.
		// this is important because we do not want any name conversion to occur.
		if err := cgroupManager.Validate(cgroupRoot); err != nil {
			return nil, fmt.Errorf("invalid configuration: %w", err)
		}
		klog.InfoS("Container manager verified user specified cgroup-root exists", "cgroupRoot", cgroupRoot)
		// Include the top level cgroup for enforcing node allocatable into cgroup-root.
		// This way, all sub modules can avoid having to understand the concept of node allocatable.
		cgroupRoot = cm.NewCgroupName(cgroupRoot, defaultNodeAllocatableCgroupName)
	}
	klog.InfoS("Creating Container Manager object based on Node Config", "nodeConfig", nodeConfig)

	return cgroupManager, nil
}

// copy from k8s.io/kubernetes/cmd/kubelet/app/server_linux.go
// ========================================================

func newContainerManager(s *options.KubeletServer, err error, kubeDeps *kubelet.Dependencies) error {
	if s.CgroupsPerQOS && s.CgroupRoot == "" {
		klog.InfoS("--cgroups-per-qos enabled, but --cgroup-root was not specified.  defaulting to /")
		s.CgroupRoot = "/"
	}

	machineInfo, err := kubeDeps.CAdvisorInterface.MachineInfo()
	if err != nil {
		return err
	}
	reservedSystemCPUs, err := getReservedCPUs(machineInfo, s.ReservedSystemCPUs)
	if err != nil {
		return err
	}
	if reservedSystemCPUs.Size() > 0 {
		// at cmd option validation phase it is tested either --system-reserved-cgroup or --kube-reserved-cgroup is specified, so overwrite should be ok
		klog.InfoS("Option --reserved-cpus is specified, it will overwrite the cpu setting in KubeReserved and SystemReserved", "kubeReservedCPUs", s.KubeReserved, "systemReservedCPUs", s.SystemReserved)
		if s.KubeReserved != nil {
			delete(s.KubeReserved, "cpu")
		}
		if s.SystemReserved == nil {
			s.SystemReserved = make(map[string]string)
		}
		s.SystemReserved["cpu"] = strconv.Itoa(reservedSystemCPUs.Size())
		klog.InfoS("After cpu setting is overwritten", "kubeReservedCPUs", s.KubeReserved, "systemReservedCPUs", s.SystemReserved)
	}

	kubeReserved, err := parseResourceList(s.KubeReserved)
	if err != nil {
		return fmt.Errorf("--kube-reserved value failed to parse: %w", err)
	}
	systemReserved, err := parseResourceList(s.SystemReserved)
	if err != nil {
		return fmt.Errorf("--system-reserved value failed to parse: %w", err)
	}
	var hardEvictionThresholds []evictionapi.Threshold
	// If the user requested to ignore eviction thresholds, then do not set valid values for hardEvictionThresholds here.
	if !s.ExperimentalNodeAllocatableIgnoreEvictionThreshold {
		hardEvictionThresholds, err = eviction.ParseThresholdConfig([]string{}, s.EvictionHard, nil, nil, nil)
		if err != nil {
			return err
		}
	}
	experimentalQOSReserved, err := cm.ParseQOSReserved(s.QOSReserved)
	if err != nil {
		return fmt.Errorf("--qos-reserved value failed to parse: %w", err)
	}

	var cpuManagerPolicyOptions map[string]string
	if utilfeature.DefaultFeatureGate.Enabled(features.CPUManagerPolicyOptions) {
		cpuManagerPolicyOptions = s.CPUManagerPolicyOptions
	} else if s.CPUManagerPolicyOptions != nil {
		return fmt.Errorf("CPU Manager policy options %v require feature gates %q, %q enabled",
			s.CPUManagerPolicyOptions, features.CPUManager, features.CPUManagerPolicyOptions)
	}

	var topologyManagerPolicyOptions map[string]string
	if utilfeature.DefaultFeatureGate.Enabled(features.TopologyManagerPolicyOptions) {
		topologyManagerPolicyOptions = s.TopologyManagerPolicyOptions
	} else if s.TopologyManagerPolicyOptions != nil {
		return fmt.Errorf("topology manager policy options %v require feature gates %q enabled",
			s.TopologyManagerPolicyOptions, features.TopologyManagerPolicyOptions)
	}

	kubeDeps.ContainerManager, err = cm.NewContainerManager(
		kubeDeps.Mounter,
		kubeDeps.CAdvisorInterface,
		cm.NodeConfig{
			RuntimeCgroupsName:    s.RuntimeCgroups,
			SystemCgroupsName:     s.SystemCgroups,
			KubeletCgroupsName:    s.KubeletCgroups,
			KubeletOOMScoreAdj:    s.OOMScoreAdj,
			CgroupsPerQOS:         s.CgroupsPerQOS,
			CgroupRoot:            s.CgroupRoot,
			CgroupDriver:          s.CgroupDriver,
			KubeletRootDir:        s.RootDirectory,
			ProtectKernelDefaults: s.ProtectKernelDefaults,
			NodeAllocatableConfig: cm.NodeAllocatableConfig{
				KubeReservedCgroupName:   s.KubeReservedCgroup,
				SystemReservedCgroupName: s.SystemReservedCgroup,
				EnforceNodeAllocatable:   sets.NewString(s.EnforceNodeAllocatable...),
				KubeReserved:             kubeReserved,
				SystemReserved:           systemReserved,
				ReservedSystemCPUs:       reservedSystemCPUs,
				HardEvictionThresholds:   hardEvictionThresholds,
			},
			QOSReserved:                             *experimentalQOSReserved,
			CPUManagerPolicy:                        s.CPUManagerPolicy,
			CPUManagerPolicyOptions:                 cpuManagerPolicyOptions,
			CPUManagerReconcilePeriod:               s.CPUManagerReconcilePeriod.Duration,
			ExperimentalMemoryManagerPolicy:         s.MemoryManagerPolicy,
			ExperimentalMemoryManagerReservedMemory: s.ReservedMemory,
			PodPidsLimit:                            s.PodPidsLimit,
			EnforceCPULimits:                        s.CPUCFSQuota,
			CPUCFSQuotaPeriod:                       s.CPUCFSQuotaPeriod.Duration,
			TopologyManagerPolicy:                   s.TopologyManagerPolicy,
			TopologyManagerScope:                    s.TopologyManagerScope,
			TopologyManagerPolicyOptions:            topologyManagerPolicyOptions,
		},
		s.FailSwapOn,
		kubeDeps.Recorder,
		kubeDeps.KubeClient,
	)

	return nil
}

// parseResourceList parses the given configuration map into an API
// ResourceList or returns an error.
func parseResourceList(m map[string]string) (v1.ResourceList, error) {
	if len(m) == 0 {
		return nil, nil
	}
	rl := make(v1.ResourceList)
	for k, v := range m {
		switch v1.ResourceName(k) {
		// CPU, memory, local storage, and PID resources are supported.
		case v1.ResourceCPU, v1.ResourceMemory, v1.ResourceEphemeralStorage, pidlimit.PIDs:
			q, err := resource.ParseQuantity(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse quantity %q for %q resource: %w", v, k, err)
			}
			if q.Sign() == -1 {
				return nil, fmt.Errorf("resource quantity for %q cannot be negative: %v", k, v)
			}
			rl[v1.ResourceName(k)] = q
		default:
			return nil, fmt.Errorf("cannot reserve %q resource", k)
		}
	}
	return rl, nil
}

func getReservedCPUs(machineInfo *cadvisorapi.MachineInfo, cpus string) (cpuset.CPUSet, error) {
	emptyCPUSet := cpuset.New()

	if cpus == "" {
		return emptyCPUSet, nil
	}

	topo, err := topology.Discover(machineInfo)
	if err != nil {
		return emptyCPUSet, fmt.Errorf("unable to discover CPU topology info: %s", err)
	}
	reservedCPUSet, err := cpuset.Parse(cpus)
	if err != nil {
		return emptyCPUSet, fmt.Errorf("unable to parse reserved-cpus list: %s", err)
	}
	allCPUSet := topo.CPUDetails.CPUs()
	if !reservedCPUSet.IsSubsetOf(allCPUSet) {
		return emptyCPUSet, fmt.Errorf("reserved-cpus: %s is not a subset of online-cpus: %s", cpus, allCPUSet.String())
	}
	return reservedCPUSet, nil
}

// copy from k8s.io/kubernetes/pkg/kubelet/cm/node_container_manager_linux.go
// ========================================================

const (
	defaultNodeAllocatableCgroupName = "kubepods"
)
