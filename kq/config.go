package kq

import "zero/core/service"

const (
	firstOffset = "first"
	lastOffset  = "last"
)

type KqConf struct {
	service.ServiceConf
	Brokers      []string
	Group        string
	Topic        string
	Offset       string `json:",options=first|last,default=last"`
	NumConns     int    `json:",default=1"`
	NumProducers int    `json:",default=8"`
	NumConsumers int    `json:",default=8"`
	MinBytes     int    `json:",default=10240"`    // 10K
	MaxBytes     int    `json:",default=10485760"` // 10M
}
