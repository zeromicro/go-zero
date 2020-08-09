package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/proc"
	"github.com/tal-tech/go-zero/core/syncx"
	"go.etcd.io/etcd/clientv3"
)

var (
	endpoints []string
	keys      = []string{
		"user.rpc",
		"classroom.rpc",
	}
	vals    = make(map[string]map[string]string)
	barrier syncx.Barrier
)

type listener struct {
	key string
}

func init() {
	cluster := proc.Env("ETCD_CLUSTER")
	if len(cluster) > 0 {
		endpoints = strings.Split(cluster, ",")
	} else {
		endpoints = []string{"localhost:2379"}
	}
}

func (l listener) OnAdd(key, val string) {
	fmt.Printf("add, key: %s, val: %s\n", key, val)
	barrier.Guard(func() {
		if m, ok := vals[l.key]; ok {
			m[key] = val
		} else {
			vals[l.key] = map[string]string{key: val}
		}
	})
}

func (l listener) OnDelete(key string) {
	fmt.Printf("del, key: %s\n", key)
	barrier.Guard(func() {
		if m, ok := vals[l.key]; ok {
			delete(m, key)
		}
	})
}

func load(cli *clientv3.Client, key string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, err
	}

	ret := make(map[string]string)
	for _, ev := range resp.Kvs {
		ret[string(ev.Key)] = string(ev.Value)
	}

	return ret, nil
}

func loadAll(cli *clientv3.Client) (map[string]map[string]string, error) {
	ret := make(map[string]map[string]string)
	for _, key := range keys {
		m, err := load(cli, key)
		if err != nil {
			return nil, err
		}

		ret[key] = m
	}

	return ret, nil
}

func compare(a, b map[string]map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k := range a {
		av := a[k]
		bv := b[k]
		if len(av) != len(bv) {
			return false
		}

		for kk := range av {
			if av[kk] != bv[kk] {
				return false
			}
		}
	}

	return true
}

func serializeMap(m map[string]map[string]string, prefix string) string {
	var builder strings.Builder
	for k, v := range m {
		fmt.Fprintf(&builder, "%s%s:\n", prefix, k)
		for kk, vv := range v {
			fmt.Fprintf(&builder, "%s\t%s: %s\n", prefix, kk, vv)
		}
	}
	return builder.String()
}

func main() {
	registry := discov.NewFacade(endpoints)
	for _, key := range keys {
		registry.Monitor(key, listener{key: key})
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			expect, err := loadAll(registry.Client().(*clientv3.Client))
			if err != nil {
				fmt.Println("[ETCD-test] can't load current keys")
				continue
			}

			check := func() bool {
				var match bool
				barrier.Guard(func() {
					match = compare(expect, vals)
				})
				if match {
					logx.Info("match")
				}
				return match
			}
			if check() {
				continue
			}

			time.AfterFunc(time.Second*5, func() {
				if check() {
					return
				}

				var builder strings.Builder
				builder.WriteString(fmt.Sprintf("expect:\n%s\n", serializeMap(expect, "\t")))
				barrier.Guard(func() {
					builder.WriteString(fmt.Sprintf("actual:\n%s\n", serializeMap(vals, "\t")))
				})
				fmt.Println(builder.String())
			})
		}
	}
}
