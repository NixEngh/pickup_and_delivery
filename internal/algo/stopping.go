package algo

import (
	"time"

	"github.com/NixEngh/pickup_and_delivery/internal/utils"
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
	utils.PrintLoadingBar(i.i, i.iterations, 50)
	return i.i > i.iterations
}

func (i *IterationBasedStopper) Reset() {
	i.i = 0
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
