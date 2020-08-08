package load

import (
	"io"

	"github.com/tal-tech/go-zero/core/syncx"
)

type ShedderGroup struct {
	options []ShedderOption
	manager *syncx.ResourceManager
}

func NewShedderGroup(opts ...ShedderOption) *ShedderGroup {
	return &ShedderGroup{
		options: opts,
		manager: syncx.NewResourceManager(),
	}
}

func (g *ShedderGroup) GetShedder(key string) Shedder {
	shedder, _ := g.manager.GetResource(key, func() (closer io.Closer, e error) {
		return nopCloser{
			Shedder: NewAdaptiveShedder(g.options...),
		}, nil
	})
	return shedder.(Shedder)
}

type nopCloser struct {
	Shedder
}

func (c nopCloser) Close() error {
	return nil
}
