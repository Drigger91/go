package main

import (
	"fmt"
	"reflect"
	"time"
)

type TaskHeap []*Task

type Task struct {
	id      string
	runAt   time.Time
	payload string
	handler func()
	index   int
}

func (t TaskHeap) Len() int {
	return len(t)

}
func (t TaskHeap) Less(i, j int) bool {
	return t[i].runAt.Before(t[j].runAt)
}
func (t TaskHeap) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
	t[i].index = i
	t[j].index = j
}
func (t *TaskHeap) Push(x any) {
	var task *Task
	if val, ok := x.(*Task); !ok {
		fmt.Printf("[Error] not a valid Push operation, expected type Task got type : %s", reflect.TypeOf(x))
		return
	} else {
		task = val
	}
	task.index = len(*t)
	*t = append(*t, task)

}
func (t *TaskHeap) Pop() any {
	if t.Len() == 0 {
		return fmt.Errorf("Invalid operation Pop() on empty heap")
	}
	oldQueue := *t
	oldLen := t.Len()
	taskToReturn := oldQueue[oldLen-1]
	taskToReturn.index = -1
	oldQueue[oldLen-1] = nil
	*t = oldQueue[:oldLen-1]
	return taskToReturn
}
