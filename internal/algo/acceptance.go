package algo

import (
	"time"

	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type Acceptor interface {
	Accept(s, newS, bestS *solution.Solution) bool
}

type GreedyAcceptor struct{}

func (a *GreedyAcceptor) Accept(s, newS, bestS *solution.Solution) bool {
	return newS.Cost() < s.Cost()
}

type IterationR2RAcceptor struct {
	iterations int
	i          int
}

func NewIterationR2RAcceptor(iterations int) *IterationR2RAcceptor {
	return &IterationR2RAcceptor{iterations: iterations}
}

func (r *IterationR2RAcceptor) Accept(s, newS, bestS *solution.Solution) bool {
	r.i++
	D := 0.2 * (float64(r.iterations-r.i) / float64(r.iterations)) * float64(bestS.Cost())

	return float64(newS.Cost()) < float64(bestS.Cost())+D
}

type TimeR2RAcceptor struct {
	totalTime  int64
	finishTime int64
}

func NewTimeR2RAcceptor(totalTime int64) *TimeR2RAcceptor {
	startTime := time.Now().Unix()

	return &TimeR2RAcceptor{totalTime: totalTime, finishTime: startTime + totalTime}
}

func (r *TimeR2RAcceptor) Accept(s, newS, bestS *solution.Solution) bool {
	currentTime := float64(time.Now().Unix())
	timeLeft := float64(r.finishTime) - currentTime

	D := 0.2 * (timeLeft / float64(r.totalTime)) * float64(bestS.Cost())

	return float64(newS.Cost()) < float64(bestS.Cost())+D
}
