package operator

import (
	"sort"

	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type RemoveCostly struct {
	n int
}

func NewRemoveCostly(n int) *RemoveCostly {
	return &RemoveCostly{n: n}
}

type valueIndex struct {
	cost  int
	index int
}

func (r *RemoveCostly) Choose(s *solution.Solution) []int {
	callCosts := s.CallCosts()

	valueIndices := make([]valueIndex, len(callCosts))
	for i, v := range callCosts {
		valueIndices[i] = valueIndex{cost: v, index: i}
	}

	sort.Slice(valueIndices, func(i, j int) bool {
		return valueIndices[i].cost > valueIndices[j].cost
	})

	topIndices := make([]int, r.n)
	for i := 0; i < r.n; i++ {
		topIndices[i] = valueIndices[i].index
	}

	return topIndices
}
