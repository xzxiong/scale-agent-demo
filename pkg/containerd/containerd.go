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

// NewContainerClient returns a containerd container client.
// cc https://github.com/containerd/containerd?tab=readme-ov-file#namespaces
func NewContainer(ctx context.Context, ns, id string) (containerd.Container, error) {
	c, err := client.NewContainer(namespaces.WithNamespace(ctx, ns), id)
	return c, err
}

// GetAllContainer returns all containers.
// filters format cc github.com/containerd/containerd/pkg/filters:Parse
// - example: matrixone.cloud/component=cn
func GetAllContainer(ctx context.Context, filters ...string) ([]containerd.Container, error) {
	return client.Containers(ctx, filters...)
}
