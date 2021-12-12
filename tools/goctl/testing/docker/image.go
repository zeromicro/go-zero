package docker

import (
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type image struct {
	cli *client.Client
}

func newImage(cli *client.Client) *image {
	return &image{cli: cli}
}

func (i *image) build(srcPath string, tags []string) error {
	if ok, _ := i.exist(tag); ok {
		return nil
	}
	tar, err := archive.TarWithOptions(srcPath, &archive.TarOptions{})
	if err != nil {
		return err
	}
	buildRet, err := i.cli.ImageBuild(context.Background(), tar, types.ImageBuildOptions{
		Tags:       tags,
		Dockerfile: "Dockerfile",
	})
	if err != nil {
		return err
	}
	defer buildRet.Body.Close()

	scanner := bufio.NewScanner(buildRet.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return scanner.Err()
}

func (i *image) exist(tag string) (bool, error) {
	list, err := i.cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return false, err
	}
	for _, l := range list {
		for _, t := range l.RepoTags {
			if t == tag {
				return true, nil
			}
		}
	}
	return false, nil
}
