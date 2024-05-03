package algo

import (
	"fmt"
	"time"
)

type Stopper interface {
	CheckStop() bool
	Reset()
}

type IterationBasedStopper struct {
	iterations int
	i          int
}

func NewIterationBasedStopper(iterations int) *IterationBasedStopper {
	return &IterationBasedStopper{iterations: iterations}
}

func (i *IterationBasedStopper) CheckStop() bool {
	i.i++
	return i.i > i.iterations
}

func (i *IterationBasedStopper) Reset() {
	i.i = 0
}

type TimeBasedStopper struct {
	createdAt int64
	time      int64
	i         int
}

func NewTimeBasedStopper(t int64) *TimeBasedStopper {
	return &TimeBasedStopper{createdAt: time.Now().Unix(), time: t}
}

func (t *TimeBasedStopper) CheckStop() bool {
	t.i++
	currentTime := time.Now().Unix()
	shouldStop := currentTime-t.createdAt > t.time
	if shouldStop {
		fmt.Println("TimeBasedStopper stopped after ", t.i, " iterations")
	}
	return shouldStop
}

func (t *TimeBasedStopper) Reset() {
	t.createdAt = time.Now().Unix()
}
