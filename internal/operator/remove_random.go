package operator

import (
	"math/rand"

	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type RemoveRandom struct {
	n int
}

func NewRemoveRandom(n int) *RemoveRandom {
	return &RemoveRandom{n: n}
}

func (r *RemoveRandom) Choose(s *solution.Solution) []int {
	calls := rand.Perm(s.Problem.NumberOfCalls)[:r.n]

	res := make([]int, r.n)
	for i, c := range calls {
		res[i] = c + 1
	}
	return res
}
