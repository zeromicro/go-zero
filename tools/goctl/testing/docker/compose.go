package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/docker/docker/client"
)

type compose struct {
	ctx        context.Context
	cli        *client.Client
	cancel     context.CancelFunc
	containers []*container
}

func newCompose(cli *client.Client) *compose {
	ctx, cancel := context.WithCancel(context.Background())
	return &compose{
		ctx:    ctx,
		cli:    cli,
		cancel: cancel,
	}
}

func (c *compose) add(ct *container) {
	c.containers = append(c.containers, ct)
}

func (c *compose) startNotRunning() error {
	for _, ct := range c.containers {
		ct := ct
		if ok, _ := ct.running(c.ctx); !ok {
			if err := c.startContainer(ct); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *compose) startContainer(ct *container) error {
	if ok, _ := ct.exist(c.ctx); !ok {
		if _, err := ct.create(c.ctx); err != nil {
			return err
		}
	}
	return ct.start(c.ctx)
}

func (c *compose) pullIfNotExist() error {
	for _, ct := range c.containers {
		out, err := ct.pullIfNotExist(c.ctx)
		if err != nil {
			return err
		}
		if out == nil {
			continue
		}
		fmt.Printf("%s: pulling image...\n", ct.image)
		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			s := struct {
				Status   string
				Progress string
			}{}
			_ = json.Unmarshal([]byte(scanner.Text()), &s)
			if s.Progress != "" {
				fmt.Printf("%s: Status: %s Progress: %v\n", ct.image, s.Status, s.Progress)
			} else {
				fmt.Printf("%s: Status: %s\n", ct.image, s.Status)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("%s pull scanner err: %v\n", ct.image, err)
		}
	}
	return nil
}

func (c *compose) waitHealthy() error {
	cs := make(map[string]int64)
	for {
		if c.ctx.Err() != nil {
			return c.ctx.Err()
		}
		health := true
		for _, ct := range c.containers {
			if ok, _ := ct.healthy(c.ctx); !ok {
				if cs[ct.image] == 0 {
					cs[ct.image] = time.Now().Unix()
				}
				if time.Now().Unix()-cs[ct.image] > 10 {
					cs[ct.image] = time.Now().UnixNano()
					fmt.Printf("%s healthy check: unhealthy\n", ct.shortRef())
				}
				health = false
				break
			}
		}
		if !health {
			time.Sleep(time.Second)
			continue
		}
		fmt.Printf("\033[1;31;40m%s\033[0m\n", "all images are started")
		return nil
	}
}

func (c *compose) runHooks() error {
	for _, ct := range c.containers {
		if len(ct.hooks) == 0 {
			continue
		}
		for _, hook := range ct.hooks {
			if len(hook.Custom) > 0 {
				if hooks[hook.Custom] == nil {
					return errors.New(fmt.Sprintf("can't find custom hook: %s", hook.Custom))
				}
				fmt.Printf("\033[1;31;40m\nrun custom hook %s\033[0m\n", hook.Custom)
				if err := hooks[hook.Custom](ct); err != nil {
					return err
				}
			}
			if len(hook.Cmd) > 0 {
				fmt.Printf("\033[1;31;40m\n%s exec command %v\033[0m\n", ct.image, hook.Cmd)
				if err := ct.exec(c.ctx, hook.Cmd); err != nil {
					return err
				}
			}
			time.Sleep(time.Second)
		}
	}
	return nil
}

func (c *compose) Run() error {
	if err := c.pullIfNotExist(); err != nil {
		return err
	}
	if err := c.startNotRunning(); err != nil {
		return err
	}
	if err := c.waitHealthy(); err != nil {
		return err
	}
	return c.runHooks()
}
