package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	con "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type container struct {
	cli          *client.Client
	name         string
	imageCfg     *ImageCfg
	image        string
	portBindings nat.PortMap
	mounts       []mount.Mount
	env          []string
	cmd          []string
	healthCheck  *con.HealthConfig
	hooks        []*Hooks
}

func newContainer(cli *client.Client, name string, imgCfg *ImageCfg, dir string) (*container, error) {
	port := nat.PortMap{}
	if len(imgCfg.Ports) > 0 {
		for _, p := range imgCfg.Ports {
			arr := strings.Split(p, ":")
			if len(arr) != 2 {
				return nil, fmt.Errorf("wrong port: %v", p)
			}
			port[nat.Port(arr[1]+"/tcp")] = []nat.PortBinding{{
				HostIP:   "127.0.0.1",
				HostPort: arr[0] + "/tcp",
			}}
		}
	}
	var ms []mount.Mount
	mounts := imgCfg.Volumes
	for _, m := range mounts {
		as := strings.Split(m, ":")
		if len(as) != 2 {
			return nil, fmt.Errorf("wrong volumn: %v", m)
		}
		path := filepath.Join(dir, as[0])
		ms = append(ms, mount.Mount{
			Type:   mount.TypeBind,
			Source: path,
			Target: as[1],
		})
	}
	var healthy *con.HealthConfig
	if imgCfg.HealthCheck != nil {
		healthy = &con.HealthConfig{
			Test:     imgCfg.HealthCheck.Test,
			Interval: imgCfg.HealthCheck.Interval,
			Timeout:  imgCfg.HealthCheck.Timeout,
			Retries:  imgCfg.HealthCheck.Retries,
		}
	}
	return &container{
		cli:          cli,
		imageCfg:     imgCfg,
		image:        imgCfg.Image,
		portBindings: port,
		mounts:       ms,
		env:          imgCfg.Environment,
		cmd:          imgCfg.Command,
		healthCheck:  healthy,
		name:         name,
		hooks:        imgCfg.Hooks,
	}, nil
}

func (ct *container) start(c context.Context) error {
	return ct.cli.ContainerStart(c, ct.name, types.ContainerStartOptions{})
}

func (ct *container) exec(c context.Context, cmd []string) error {
	ret, err := ct.cli.ContainerExecCreate(c, ct.name, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	})
	if err != nil {
		return err
	}
	execRet, err := ct.cli.ContainerExecAttach(c, ret.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}
	defer execRet.Close()
	scanner := bufio.NewScanner(execRet.Reader)

	var dir string
	var idx int
	for i, c := range cmd {
		if c == "-dir" {
			idx = i
			break
		}
	}
	if len(cmd) > idx {
		dir = cmd[idx+1]
	}
	// FIXME
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Done") && len(cmd) >= 2 && (cmd[1] == "api" || cmd[1] == "rpc") {
			ct.tree(dir)
		} else {
			fmt.Println(scanner.Text())
		}
	}
	return nil
}

func (ct *container) running(c context.Context) (bool, error) {
	status, err := ct.cli.ContainerInspect(c, ct.name)
	if err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return status.State.Running, nil
}

func (ct *container) exist(c context.Context) (bool, error) {
	if _, err := ct.cli.ContainerInspect(c, ct.name); err != nil {
		if client.IsErrNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ct *container) create(c context.Context) (con.ContainerCreateCreatedBody, error) {
	host := &con.HostConfig{
		LogConfig:    con.LogConfig{},
		PortBindings: ct.portBindings,
		Mounts:       ct.mounts,
	}
	return ct.cli.ContainerCreate(c, &con.Config{
		Image:       ct.longRef(),
		Env:         ct.env,
		Cmd:         ct.cmd,
		Healthcheck: ct.healthCheck,
	}, host, nil, nil, ct.name)
}

func (ct *container) pullIfNotExist(c context.Context) (io.ReadCloser, error) {
	list, err := ct.cli.ImageList(c, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	for _, l := range list {
		for _, t := range l.RepoTags {
			if t == ct.shortRef() {
				return nil, nil
			}
		}
	}
	return ct.pull(c)
}

func (ct *container) pull(c context.Context) (io.ReadCloser, error) {
	return ct.cli.ImagePull(c, ct.longRef(), types.ImagePullOptions{})
}

func (ct *container) healthy(c context.Context) (bool, error) {
	inspect, err := ct.cli.ContainerInspect(c, ct.name)
	if err != nil {
		return false, err
	}
	health := inspect.ContainerJSONBase.State.Health
	if health == nil {
		return true, nil
	}
	return health.Status == types.Healthy, nil
}

func (ct *container) shortRef() string {
	ref := ct.image
	if strings.Contains(ref, ":") {
		return ref
	}
	return ref + ":latest"
}

func (ct *container) longRef() string {
	ref := ct.image
	arr := strings.Split(ref, "/")
	if len(arr) == 1 {
		return "docker.io/library/" + ref
	}
	if len(arr) == 2 {
		return "docker.io/" + ref
	}
	return ref
}

func (ct *container) tree(dir string) {
	ret, err := ct.cli.ContainerExecCreate(context.Background(), ct.name, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"tree", dir},
	})
	if err != nil {
		return
	}
	execRet, err := ct.cli.ContainerExecAttach(context.Background(), ret.ID, types.ExecStartCheck{})
	if err != nil {
		return
	}
	defer execRet.Close()
	scanner := bufio.NewScanner(execRet.Reader)

	fmt.Printf("\033[1;36;40m\n%s\033[0m", "the directory structure")
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
