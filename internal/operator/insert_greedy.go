package operator

import (
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
)

type InsertGreedy struct {
}

func NewInsertGreedy() *InsertGreedy {
	return &InsertGreedy{}
}

func (i *InsertGreedy) Reinsert(s *solution.Solution, calls []int) {
	for _, c := range calls {
		inds := utils.FindIndices(s.Solution, c)
		s.MoveCallToOutsource(c, inds)
	}

	for i := len(calls) - 1; i >= 0; i-- {
		feasible := s.GetAllFeasible(calls[i])
		if len(feasible) == 0 {
			continue
		}
		if feasible[0].CostDiff < 0 {
			inds := utils.FindIndices(s.Solution, calls...)
			s.InsertCall(calls[i], inds, feasible[0])
		}
	}
}
