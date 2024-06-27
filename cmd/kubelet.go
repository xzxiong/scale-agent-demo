/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

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
		fmt.Printf("kubelet called mode: %s\n", kubelet.Mode)
		if *kPid <= 0 {
			fmt.Printf("[Error] invalid pid: %d", *kPid)
			os.Exit(1)
		}
		if *kCpu {
			fmt.Printf("kubelet cpu profiling for pid %d\n", *kPid)
			fmt.Printf("(not support yet\n")
		}
	},
}

var kPid *int
var kCpu *bool

func init() {
	rootCmd.AddCommand(kubeletCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kubeletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kubeletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	kPid = kubeletCmd.Flags().IntP("pid", "p", 0, "process id")
	kCpu = kubeletCmd.Flags().BoolP("cpu", "c", false, "Show cpu info")
}
