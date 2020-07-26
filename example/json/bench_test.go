package testjson

import (
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
	segment "github.com/segmentio/encoding/json"
)

const input = `{"@timestamp":"2020-02-12T14:02:10.849Z","@metadata":{"beat":"filebeat","type":"doc","version":"6.1.1","topic":"k8slog"},"index":"k8slog","offset":908739,"stream":"stdout","topic":"k8slog","k8s_container_name":"shield-rpc","k8s_pod_namespace":"xx-xiaoheiban","stage":"gray","prospector":{"type":"log"},"k8s_node_name":"cn-hangzhou.i-bp15w8irul9hmm3l9mxz","beat":{"name":"log-pilot-7s6qf","hostname":"log-pilot-7s6qf","version":"6.1.1"},"source":"/host/var/lib/docker/containers/4e6dca76f3e38fb8b39631e9bb3a19f9150cc82b1dab84f71d4622a08db20bfb/4e6dca76f3e38fb8b39631e9bb3a19f9150cc82b1dab84f71d4622a08db20bfb-json.log","level":"info","duration":"39.425µs","content":"172.25.5.167:49976 - /remoteshield.Filter/Filter - {\"sentence\":\"王XX2月12日作业\"}","k8s_pod":"shield-rpc-57c9dc6797-55skf","docker_container":"k8s_shield-rpc_shield-rpc-57c9dc6797-55skf_xx-xiaoheiban_a8341ba0-30ee-11ea-8ac4-00163e0fb3ef_0"}`

func BenchmarkStdJsonMarshal(b *testing.B) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		b.FailNow()
	}
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(m); err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkJsonIteratorMarshal(b *testing.B) {
	m := make(map[string]interface{})
	if err := jsoniter.Unmarshal([]byte(input), &m); err != nil {
		b.FailNow()
	}
	for i := 0; i < b.N; i++ {
		if _, err := jsoniter.Marshal(m); err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkSegmentioMarshal(b *testing.B) {
	m := make(map[string]interface{})
	if err := segment.Unmarshal([]byte(input), &m); err != nil {
		b.FailNow()
	}
	for i := 0; i < b.N; i++ {
		if _, err := jsoniter.Marshal(m); err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkStdJsonUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make(map[string]interface{})
		if err := json.Unmarshal([]byte(input), &m); err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkJsonIteratorUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make(map[string]interface{})
		if err := jsoniter.Unmarshal([]byte(input), &m); err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkSegmentioUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make(map[string]interface{})
		if err := segment.Unmarshal([]byte(input), &m); err != nil {
			b.FailNow()
		}
	}
}
