package discov

import "zero/core/logx"

type (
	Renewer interface {
		Start()
		Stop()
		Pause()
		Resume()
	}

	etcdRenewer struct {
		*Publisher
	}
)

func NewRenewer(endpoints []string, key, value string, renewId int64) Renewer {
	var publisher *Publisher
	if renewId > 0 {
		publisher = NewPublisher(endpoints, key, value, WithId(renewId))
	} else {
		publisher = NewPublisher(endpoints, key, value)
	}

	return &etcdRenewer{
		Publisher: publisher,
	}
}

func (sr *etcdRenewer) Start() {
	if err := sr.KeepAlive(); err != nil {
		logx.Error(err)
	}
}
