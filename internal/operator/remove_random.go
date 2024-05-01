package operator

import (
	"math"
	"math/rand"

	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type RemoveRandom struct {
	percent int
}

func NewRemoveRandom(percent int) *RemoveRandom {
	return &RemoveRandom{percent: percent}
}

func (r *RemoveRandom) Choose(s *solution.Solution) []int {
	n := int(math.Ceil(float64(s.Problem.NumberOfCalls) * float64(r.percent) / 100))
	calls := rand.Perm(s.Problem.NumberOfCalls)[:n]

	res := make([]int, n)
	for i, c := range calls {
		res[i] = c + 1
	}
	return res
}
