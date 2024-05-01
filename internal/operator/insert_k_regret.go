package operator

import (
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
)

type InsertKRegret struct {
	k int
}

func NewInsertKRegret(k int) *InsertKRegret {
	return &InsertKRegret{k: k}
}

func (i *InsertKRegret) Reinsert(s *solution.Solution, calls []int) {
	for _, c := range calls {
		inds := utils.FindIndices(s.Solution, c)
		s.MoveCallToOutsource(c, inds)
	}

	for iter := 0; iter < len(calls); iter++ {
		bestFeasible := make([]utils.InsertionPoint, len(calls))
		bestCall := -1
		highestRegret := -1000000
		for _, c := range calls {
			feasible := s.GetAllFeasible(c)
			k := min(i.k, len(feasible))
			if len(feasible) == 0 {
				continue
			}
			sum := 0
			for j := range feasible[1:k] {
				sum += feasible[j].CostDiff - feasible[0].CostDiff
			}
			if sum > highestRegret {
				bestCall = c
				highestRegret = sum
				bestFeasible = feasible
			}
		}

		if bestCall == -1 {
			continue
		}

		inds := utils.FindIndices(s.Solution, bestCall)
		s.InsertCall(bestCall, inds, bestFeasible[0])
	}

}
