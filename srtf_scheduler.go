package main

import (
	"container/heap"
	"sync"
)

// Shortest Remaining Time First
type SRTFScheduler struct {
	sync.Mutex
	tasks *SRTFTaskHeap
}

func (srtf *SRTFScheduler) push_task(task *Task) {
	srtf.Lock()
	defer srtf.Unlock()
	heap.Push(srtf.tasks, task)
}

func (srtf *SRTFScheduler) pop_task() *Task {
	srtf.Lock()
	defer srtf.Unlock()

	return heap.Pop(srtf.tasks).(*Task)
}

func NewSRTFSchduler() *SRTFScheduler {
	h := &SRTFTaskHeap{}
	heap.Init(h)

	return &SRTFScheduler{tasks: h}
}

func (srtf *SRTFScheduler) AddTask(task *Task) {
	srtf.push_task(task)
}

func (srtf *SRTFScheduler) Schedule(bandwidth uint64) []*Task {
	schduled_tasks := make([]*Task, 0, 8)

	for bandwidth != 0 && srtf.tasks.Len() != 0 {
		current := srtf.pop_task()
		schduled_tasks = append(schduled_tasks, current)
		bandwidth = current.run(bandwidth)

		if !current.finished() {
			srtf.push_task(current)
		}
	}

	return schduled_tasks
}

type SRTFTaskHeap []*Task

func (h SRTFTaskHeap) Len() int { return len(h) }

// 在所有可调度的任务中，优先选择剩余执⾏时间越短的任务
// 如果剩余时间相同，优先调度索引编号更⼩的任务
func (h SRTFTaskHeap) Less(i, j int) bool {
	if h[i].Remain < h[j].Remain {
		return true
	}

	if h[i].Remain > h[j].Remain {
		return false
	}

	return h[i].ID < h[j].ID
}

func (h SRTFTaskHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *SRTFTaskHeap) Push(x any) {
	*h = append(*h, x.(*Task))
}

func (h *SRTFTaskHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
