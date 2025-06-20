package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFifoScheudler(t *testing.T) {
	sched := NewFifoSchduler()

	sched.AddTask(&Task{ID: 0, Remain: 2})
	sched.AddTask(&Task{ID: 1, Remain: 4})
	sched.AddTask(&Task{ID: 2, Remain: 100})
	sched.AddTask(&Task{ID: 3, Remain: 6})

	scheduled_tasks := sched.Schedule(5)
	assert.Equal(t, 2, len(scheduled_tasks))
	assert.Equal(t, "0,1\t0,1", schedule_log(scheduled_tasks))

	scheduled_tasks = sched.Schedule(5)
	assert.Equal(t, 2, len(scheduled_tasks))
	assert.Equal(t, "1,2\t0,96", schedule_log(scheduled_tasks))

	scheduled_tasks = sched.Schedule(5)
	assert.Equal(t, 1, len(scheduled_tasks))
	assert.Equal(t, "2\t91", schedule_log(scheduled_tasks))

}
