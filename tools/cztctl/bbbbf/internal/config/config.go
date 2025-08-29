package config

import (
	"github.com/lerity-yao/go-mq/rabbitmq"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	GBoxCommonRabbitmqConf  rabbitmq.RabbitListenerConf
	GBoxCommon1RabbitmqConf rabbitmq.RabbitListenerConf
}
