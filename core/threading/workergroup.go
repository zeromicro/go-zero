package threading

type WorkerGroup struct {
	job     func()
	workers int
}

func NewWorkerGroup(job func(), workers int) WorkerGroup {
	return WorkerGroup{
		job:     job,
		workers: workers,
	}
}

func (wg WorkerGroup) Start() {
	group := NewRoutineGroup()
	for i := 0; i < wg.workers; i++ {
		group.RunSafe(wg.job)
	}
	group.Wait()
}
