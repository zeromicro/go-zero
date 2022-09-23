package consul

import (
	"fmt"
	"testing"
)

func TestLoadConf(t *testing.T) {
	conf := &Conf{
		Host:     "127.0.0.1:8500",
		ListenOn: "192.168.5.216:9100",
		Key:      "core.rpc",
		Token:    "",
		Tag:      []string{"core", "rpc"},
		Meta:     map[string]string{"Protocol": "grpc"},
		TTL:      0,
	}

	client, err := conf.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	type User struct {
		Name string `json:"name" yaml:"Name"`
		Age  string `json:"age" yaml:"Age"`
	}

	// data in consul like below
	// Name: Jack
	// Age: 18

	// get config
	var u User
	LoadYAMLConf(client, "core", &u)
	fmt.Println(u)

}
