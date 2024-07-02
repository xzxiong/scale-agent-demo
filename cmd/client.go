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
		ctx := context.Background()
		if *cNode {
			node := client.GetNodeName(ctx)
			fmt.Printf("node: %s\n", node)
			return
		}
	},
}

var cNode *bool

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kubeletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kubeletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cNode = clientCmd.Flags().BoolP("node", "n", false, "Show node name")
}
