package tests

import (
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/suyuan32/simple-admin-tools/plugins/registry/consul"
)

func TestClient(t *testing.T) {
	svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, "round_robin")
	conn, err := grpc.Dial("consul://127.0.0.1:8500/gozero?wait=14s&tag=public", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(svcCfg))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	time.Sleep(29 * time.Second)
}
