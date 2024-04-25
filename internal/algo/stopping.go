package algo

import "time"

type Stopper interface {
	CheckStop() bool
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
	return i.i <= i.iterations
}

type TimeBasedStopper struct {
	createdAt int64
	time      int64
}

func NewTimeBasedStopper(t int64) *TimeBasedStopper {
	return &TimeBasedStopper{createdAt: time.Now().Unix(), time: t}
}

func (t *TimeBasedStopper) CheckStop() bool {
	currentTime := time.Now().Unix()
	return currentTime-t.createdAt <= t.time
}
