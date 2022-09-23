package consul

import (
	"fmt"
	"gopkg.in/yaml.v2"
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

	kv := client.KV()

	type User struct {
		Name string `json:"name" yaml:"Name"`
		Age  string `json:"age" yaml:"Age"`
	}

	// data in consul like below
	// Name: Jack
	// Age: 18

	// get config
	pair, _, err := kv.Get("user", nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("KV: %v \n%s\n", pair.Key, pair.Value)

	var u User
	err = yaml.Unmarshal(pair.Value, &u)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(u)
}
