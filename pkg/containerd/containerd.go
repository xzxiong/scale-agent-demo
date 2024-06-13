package containerd

import (
	"context"
	"path"

	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/pkg/namespaces"
	// "github.com/containerd/containerd/v2/pkg/cio"
)

var client *containerd.Client

const defaultAddress = "/run/containerd/containerd.sock"

func InitClient(ctx context.Context, rootDir string) (err error) {
	client, err = containerd.New(path.Join(rootDir, defaultAddress))
	return err
}

// GetContainerClient returns a containerd container client.
// cc https://github.com/containerd/containerd?tab=readme-ov-file#namespaces
func GetContainerClient(ctx context.Context, ns, id string) (containerd.Container, error) {
	c, err := client.NewContainer(namespaces.WithNamespace(ctx, ns), id)
	return c, err
}
