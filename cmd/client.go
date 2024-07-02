/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

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
		node := client.GetNodeName(ctx)
		fmt.Printf("node: %s\n", node)
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
		pods := client.ListPodsByNode(ctx, node)
		for _, pod := range pods {
			fmt.Printf("%10s %20s\n", pod.Namespace, pod.Name)
		}
		fmt.Printf("cnt: %d\n", len(pods))
		return
	},
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
	cNode = clientCmd.Flags().StringP("nodes", "n", "", "target node name")
}
