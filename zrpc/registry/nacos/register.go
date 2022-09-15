package nacos

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"

	"github.com/zeromicro/go-zero/core/netx"
)

// RegisterService register service to nacos
func RegisterService(opts *Options) error {
	pubListenOn := figureOutListenOn(opts.ListenOn)

	host, ports, err := net.SplitHostPort(pubListenOn)
	if err != nil {
		return fmt.Errorf("failed parsing address error: %v", err)
	}
	port, _ := strconv.ParseUint(ports, 10, 16)

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: opts.ServerConfig,
			ClientConfig:  opts.ClientConfig,
		},
	)
	if err != nil {
		log.Panic(err)
	}

	// service register
	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		ServiceName: opts.ServiceName,
		Ip:          host,
		Port:        port,
		Weight:      opts.Weight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    opts.Metadata,
		ClusterName: opts.Cluster,
		GroupName:   opts.Group,
	})

	if err != nil {
		return err
	}

	// service deregister
	proc.AddShutdownListener(func() {
		_, err := client.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          host,
			Port:        port,
			ServiceName: opts.ServiceName,
			Cluster:     opts.Cluster,
			GroupName:   opts.Group,
			Ephemeral:   true,
		})
		if err != nil {
			logx.Info("deregister service error: ", err.Error())
		} else {
			logx.Info("deregistered service from nacos server.")
		}
	})

	return nil
}

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIP)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}
