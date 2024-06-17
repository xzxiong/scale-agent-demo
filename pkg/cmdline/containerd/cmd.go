package containerd

import (
	"context"
	"fmt"
	"os"

	"github.com/xzxiong/scale-agent-demo/pkg/container"
)

func List(ctx context.Context) error {
	cs, err := container.GetAllContainer(ctx)
	if err != nil {
		fmt.Printf("Errro: %s", err.Error())
		os.Exit(1)
	}
	for _, c := range cs {
		if labels, err := c.Labels(ctx); err == nil {
			fmt.Printf("c %s labels: %v\n", c.ID(), labels)
		}
	}

	//rpcClient, err := container.GetContrainerdClient(ctx)
	//if err != nil {
	//	return err
	//}
	//rpcClient.UpdateContainerResources()
	return nil
}
