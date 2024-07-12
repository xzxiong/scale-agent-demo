/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"

	"github.com/xzxiong/scale-agent-demo/pkg/cmdline/client"
)

// clientCmd represents the kubelet command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cmd 'client' called\n")
		//ctx := context.Background()
	},
}

var clientNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "show all nodes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		nodeName := client.GetNodeName(ctx)
		fmt.Printf("node: %s\n", nodeName)
		// get access ip
		node := client.GetNode(ctx, nodeName)
		for _, addr := range node.Status.Addresses {
			fmt.Printf("node host type: %s\n", addr.Type)
			if addr.Type == corev1.NodeInternalIP {
				fmt.Printf("node host addr: %s [target]\n", addr.Address)
			} else {
				fmt.Printf("node host addr: %s\n", addr.Address)
			}
		}

		return
	},
}

var clientPodCmd = &cobra.Command{
	Use:   "pod",
	Short: "show all pod which belong to node (current pod's node)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		node := *cNode
		if *cNode == "" {
			node = client.GetNodeName(ctx)
		}
		fmt.Printf("node: %s\n", node)
		//pods := client.ListPodsByNode(ctx, node)
		//showAllPods(pods)

		fmt.Printf(">> list (filter by node)\n")
		pods := client.ListPodsByNodeName(ctx, node)
		showAllPods(pods)
		return
	},
}

func showAllPods(pods []*corev1.Pod) {
	nsLength := 10
	for _, pod := range pods {
		if l := len(pod.Namespace); l > nsLength {
			nsLength = l
		}
	}
	formatter := fmt.Sprintf("%%%ds %%s\n", nsLength)
	fmt.Printf(formatter, "Namespace", "Pod")
	for _, pod := range pods {
		fmt.Printf(formatter, pod.Namespace, pod.Name)
	}
	fmt.Printf("cnt: %d\n", len(pods))
}

var cNode *string

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.AddCommand(clientNodeCmd)
	clientCmd.AddCommand(clientPodCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kubeletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kubeletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//
	// part clientCmd
	//cNode = clientCmd.Flags().StringP("nodes", "n", "", "target node name")
	//
	// part clientPodCmd
	cNode = clientPodCmd.Flags().StringP("nodes", "n", "", "target node name")
}
