/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/xzxiong/scale-agent-demo/pkg/cmdline/client"
	"github.com/xzxiong/scale-agent-demo/pkg/cmdline/kubelet"
)

// kubeletCmd represents the kubelet command
var kubeletCmd = &cobra.Command{
	Use:   "kubelet",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kubelet called")
		if *kShow {
			s, err := kubelet.GetKubeletServer()
			if err != nil {
				fmt.Printf("get kubelet server error: %s\n", err)
				os.Exit(1)
			}
			fmt.Println("CgroupDriver:  ", s.CgroupDriver)
			fmt.Println("CgroupsPerQOS: ", s.CgroupsPerQOS)
			fmt.Println("CgroupRoot:    ", s.CgroupRoot)
			fmt.Println("QOSReserved:   ", s.QOSReserved)
			fmt.Println("Config:        ", s.KubeletConfigFile)
		}
		if *kPid > 0 {
			fmt.Printf("param pid: %d", *kPid)
		}
		if *kCpu {
			fmt.Printf("kubelet cpu profiling for pid %d\n", *kPid)
			fmt.Printf("(not support yet\n")
		}
		if kPod != "" {
			fmt.Printf("kubelet cpu profiling for pid %d\n", *kPid)

		}
	},
}

var kubeletCpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		fmt.Println("kubelet cpu called")
		if kNamespace == "" {
			fmt.Printf("invalid param [namespace]: %s\n", kNamespace)
			os.Exit(1)
		}
		if kPod == "" {
			fmt.Printf("invalid param [pod]: %s\n", kPod)
			os.Exit(1)
		}

		// get cpu max
		pod := client.GetPod(ctx, kNamespace, kPod)
		if pod == nil {
			fmt.Printf("pod not found\n")
			os.Exit(1)
		}
		cfg := kubelet.GetCgroupCpu(pod)
		fmt.Printf(`
CpuQuota : %d
CpuPeriod: %d
CpuShares: %d
`, *cfg.CPUQuota, *cfg.CPUPeriod, *cfg.CPUShares)
	},
}

var (
	kPid  *int
	kCpu  *bool
	kShow *bool

	kNamespace string
	kPod       string
)

func init() {
	rootCmd.AddCommand(kubeletCmd)
	kubeletCmd.AddCommand(kubeletCpuCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kubeletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kubeletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	kPid = kubeletCmd.Flags().IntP("pid", "p", 0, "process id")
	kCpu = kubeletCmd.Flags().BoolP("cpu", "c", false, "Show cpu info")
	kShow = kubeletCmd.Flags().BoolP("show", "s", false, "Show kubelet key config")

	kubeletCpuCmd.Flags().StringVarP(&kNamespace, "namespace", "n", "default", "target pod's namespace")
	kubeletCpuCmd.Flags().StringVarP(&kPod, "pod", "p", "", "target pod")
}
