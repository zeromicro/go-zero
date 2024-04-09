package fx

import "github.com/zeromicro/go-zero/core/threading"

// Parallel runs fns parallelly and waits for done.
func Parallel(fns ...func()) {
	group := threading.NewRoutineGroup()
	for _, fn := range fns {
		group.RunSafe(fn)
	}
	group.Wait()
}

// ParallelFnErr Execute a set of functions in parallel, each returning the error.
func ParallelFnErr(fns ...func() error) error {
	group := threading.NewRoutineErrGroup()
	group.SetLimit(len(fns))

	for _, fn := range fns {
		group.RunSafe(fn)
	}

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}
