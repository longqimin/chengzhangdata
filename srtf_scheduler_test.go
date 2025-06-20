package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSRTFScheudler(t *testing.T) {

	sched := NewSRTFSchduler()

	sched.AddTask(&Task{ID: 0, Remain: 2})
	sched.AddTask(&Task{ID: 1, Remain: 4})
	sched.AddTask(&Task{ID: 2, Remain: 100})
	sched.AddTask(&Task{ID: 3, Remain: 6})

	scheduled_tasks := sched.Schedule(5)
	assert.Equal(t, 2, len(scheduled_tasks))
	assert.Equal(t, "0,1\t0,1", schedule_log(scheduled_tasks))

	scheduled_tasks = sched.Schedule(5)
	assert.Equal(t, 2, len(scheduled_tasks))
	assert.Equal(t, "1,3\t0,2", schedule_log(scheduled_tasks))

	scheduled_tasks = sched.Schedule(5)
	assert.Equal(t, 2, len(scheduled_tasks))
	assert.Equal(t, "3,2\t0,97", schedule_log(scheduled_tasks))

}
