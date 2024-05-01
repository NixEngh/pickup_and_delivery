package operator

import (
	"math"
	"sort"

	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type RemoveCostly struct {
	percent int
}

func NewRemoveCostly(percent int) *RemoveCostly {
	return &RemoveCostly{percent: percent}
}

type valueIndex struct {
	cost  int
	index int
}

func (r *RemoveCostly) Choose(s *solution.Solution) []int {
	n := int(math.Ceil(float64(s.Problem.NumberOfCalls) * float64(r.percent) / 100))
	callCosts := s.CallCosts()

	valueIndices := make([]valueIndex, len(callCosts))
	for i, v := range callCosts {
		valueIndices[i] = valueIndex{cost: v, index: i}
	}

	sort.Slice(valueIndices, func(i, j int) bool {
		return valueIndices[i].cost > valueIndices[j].cost
	})

	topIndices := make([]int, n)
	for i := 0; i < n; i++ {
		topIndices[i] = valueIndices[i].index
	}

	return topIndices
}
