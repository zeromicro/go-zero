package service

import (
	"log"

	"github.com/tal-tech/go-zero/core/proc"
	"github.com/tal-tech/go-zero/core/syncx"
	"github.com/tal-tech/go-zero/core/threading"
)

type (
	Starter interface {
		Start()
	}

	Stopper interface {
		Stop()
	}

	Service interface {
		Starter
		Stopper
	}

	ServiceGroup struct {
		services []Service
		stopOnce func()
	}
)

func NewServiceGroup() *ServiceGroup {
	sg := new(ServiceGroup)
	sg.stopOnce = syncx.Once(sg.doStop)
	return sg
}

func (sg *ServiceGroup) Add(service Service) {
	sg.services = append(sg.services, service)
}

// There should not be any logic code after calling this method, because this method is a blocking one.
// Also, quitting this method will close the logx output.
func (sg *ServiceGroup) Start() {
	proc.AddShutdownListener(func() {
		log.Println("Shutting down...")
		sg.stopOnce()
	})

	sg.doStart()
}

func (sg *ServiceGroup) Stop() {
	sg.stopOnce()
}

func (sg *ServiceGroup) doStart() {
	routineGroup := threading.NewRoutineGroup()

	for i := range sg.services {
		service := sg.services[i]
		routineGroup.RunSafe(func() {
			service.Start()
		})
	}

	routineGroup.Wait()
}

func (sg *ServiceGroup) doStop() {
	for _, service := range sg.services {
		service.Stop()
	}
}

func WithStart(start func()) Service {
	return startOnlyService{
		start: start,
	}
}

func WithStarter(start Starter) Service {
	return starterOnlyService{
		Starter: start,
	}
}

type (
	stopper struct {
	}

	startOnlyService struct {
		start func()
		stopper
	}

	starterOnlyService struct {
		Starter
		stopper
	}
)

func (s stopper) Stop() {
}

func (s startOnlyService) Start() {
	s.start()
}
