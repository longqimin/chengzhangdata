package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const (
	bandwidth = 5
)

func schedule_log(tasks []*Task) string {
	ids := make([]string, 0, len(tasks))
	remain_costs := make([]string, 0, len(tasks))
	for _, task := range tasks {
		ids = append(ids, fmt.Sprintf("%d", task.ID))
		remain_costs = append(remain_costs, fmt.Sprintf("%d", task.Remain))
	}

	return fmt.Sprintf("%s\t%s", strings.Join(ids, ","), strings.Join(remain_costs, ","))
}

var task_id uint64 = uint64(0)

// 1<<20 is big enough for tasks queue
var task_chan_buffer_size int = 1 << 20

var default_scheduler string = "fifo"

// assing ID for each input Task without ID
func assign_task_id(cost_ch chan uint64, task_ch chan *Task) {
	for cost := range cost_ch {
		task_ch <- &Task{ID: task_id, Remain: cost}
		task_id++
	}
}

// move Task from task channel to Scheduler
func add_task(s Scheduler, task_ch chan *Task) {
	for task := range task_ch {
		s.AddTask(task)
	}
}

func schdule_loop(s Scheduler) {
	sched_time := 0
	for {
		tasks := s.Schedule(bandwidth)
		if len(tasks) == 0 {
			runtime.Gosched()
		} else {
			sched_time++
			fmt.Printf("%d\t%s\n", sched_time, schedule_log(tasks))
		}
	}
}

func main() {
	scheduler_name := default_scheduler
	if len(os.Args) >= 2 {
		scheduler_name = os.Args[1]
		if scheduler_name == "fifo" || scheduler_name == "srtf" {
			default_scheduler = scheduler_name
		} else {
			fmt.Printf("invalid scheduler\nusage: %s [scheduler] - schceduler can be `fifo` or `srtf`. (default `fifo`)\n", os.Args[0])
			os.Exit(1)
		}
	}

	var sched Scheduler
	if scheduler_name == "fifo" {
		sched = NewFifoSchduler()
	} else {
		sched = NewSRTFSchduler()
	}

	cost_ch := make(chan uint64, task_chan_buffer_size)
	task_ch := make(chan *Task, task_chan_buffer_size)

	go assign_task_id(cost_ch, task_ch)
	go add_task(sched, task_ch)
	go schdule_loop(sched)

	// POST /submit -d'[2,4,100,6]'
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		var task_costs []uint64
		if err := json.NewDecoder(r.Body).Decode(&task_costs); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, cost := range task_costs {
			cost_ch <- cost
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
