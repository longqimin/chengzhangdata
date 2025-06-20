package main

type FifoSchduler struct {
	ch      chan *Task
	current *Task
}

func NewFifoSchduler() *FifoSchduler {
	return &FifoSchduler{ch: make(chan *Task, task_chan_buffer_size), current: nil}
}

func (fifo *FifoSchduler) AddTask(task *Task) {
	fifo.ch <- task
}

func (fifo *FifoSchduler) Schedule(bandwidth uint64) []*Task {
	schduled_tasks := make([]*Task, 0, 8)

	for bandwidth != 0 {
		if fifo.current == nil {
			select {
			case fifo.current = <-fifo.ch:
			default:
				if len(schduled_tasks) != 0 {
					return schduled_tasks
				}
				fifo.current = <-fifo.ch // block to wait for new task
			}
		}

		schduled_tasks = append(schduled_tasks, fifo.current)
		bandwidth = fifo.current.run(bandwidth)
		if fifo.current.finished() {
			fifo.current = nil
		}
	}

	return schduled_tasks
}
