package operator

import "github.com/NixEngh/pickup_and_delivery/internal/solution"

type Operator interface {
	Apply(s *solution.Solution) int
}

type Insert interface {
	Reinsert(s *solution.Solution, calls []int)
}

// Doesn't actually remove, but chooses calls
type Removal interface {
	Choose(s *solution.Solution) []int
}
