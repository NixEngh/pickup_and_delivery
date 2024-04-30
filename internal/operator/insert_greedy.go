package operator

import (
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
)

type InsertGreedy struct {
}

func (i *InsertGreedy) Insert(s *solution.Solution, calls []int) {
	for _, c := range calls {
		inds := utils.FindIndices(s.Solution, c)
		s.MoveCallToOutsource(c, inds)
	}

	inds := utils.FindIndices(s.Solution, calls...)
	for i := len(calls) - 1; i <= 0; i-- {
		feasible := s.GetAllFeasible(calls[i])

		if feasible[0].CostDiff < 0 {
			s.InsertCall(calls[i], inds, feasible[0])
		}
	}
}
