package main

type Task struct {
	ID     uint64
	Remain uint64
}

func (task *Task) run(bandwidth uint64) (remain_bandwidth uint64) {
	if task.Remain >= bandwidth {
		task.Remain -= bandwidth
		return 0
	}

	remain_bandwidth = bandwidth - task.Remain
	task.Remain = 0
	return remain_bandwidth
}

func (task *Task) finished() bool {
	return task.Remain == 0
}

type Scheduler interface {
	AddTask(*Task)
	Schedule(bandwidth uint64) []*Task
}
