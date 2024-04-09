package threading

import (
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/rescue"
)

// A RoutineGroup is used to group goroutines together and all wait all goroutines to be done.
type RoutineGroup struct {
	waitGroup sync.WaitGroup
}

// NewRoutineGroup returns a RoutineGroup.
func NewRoutineGroup() *RoutineGroup {
	return new(RoutineGroup)
}

// Run runs the given fn in RoutineGroup.
// Don't reference the variables from outside,
// because outside variables can be changed by other goroutines
func (g *RoutineGroup) Run(fn func()) {
	g.waitGroup.Add(1)

	go func() {
		defer g.waitGroup.Done()
		fn()
	}()
}

// RunSafe runs the given fn in RoutineGroup, and avoid panics.
// Don't reference the variables from outside,
// because outside variables can be changed by other goroutines
func (g *RoutineGroup) RunSafe(fn func()) {
	g.waitGroup.Add(1)

	GoSafe(func() {
		defer g.waitGroup.Done()
		fn()
	})
}

// Wait waits all running functions to be done.
func (g *RoutineGroup) Wait() {
	g.waitGroup.Wait()
}

// A RoutineErrGroup is used to group goroutines together and all wait all goroutines to be done.
// but adds handling of tasks returning errors.
type RoutineErrGroup struct {
	errGroup   errgroup.Group
	batchError errorx.BatchError
}

// NewRoutineErrGroup returns a RoutineErrGroup.
func NewRoutineErrGroup() *RoutineErrGroup {
	return new(RoutineErrGroup)
}

// SetLimit limits the number of active goroutines in this group to at most n.
// A negative value indicates no limit.
func (g *RoutineErrGroup) SetLimit(n int) {
	g.errGroup.SetLimit(n)
}

// Run Execute the given function in a concurrent error group.
// If an error occurs during function execution, the error is added to the error batch.
func (g *RoutineErrGroup) Run(fn func() error) {
	g.errGroup.Go(func() error {
		if err := fn(); err != nil {
			g.batchError.Add(err)
			return err
		}
		return nil
	})
}

// RunSafe Executes the given function safely in the group of concurrent errors and avoid panics.
// If an error occurs during function execution, the error is added to the error batch.
func (g *RoutineErrGroup) RunSafe(fn func() error) {
	g.errGroup.Go(func() error {
		defer rescue.Recover()

		if err := fn(); err != nil {
			g.batchError.Add(err)
			return err
		}
		return nil
	})
}

// Wait waits all running functions to be done.
func (g *RoutineErrGroup) Wait() error {
	if err := g.errGroup.Wait(); err != nil {
		return g.batchError.Err()
	}
	return nil
}
