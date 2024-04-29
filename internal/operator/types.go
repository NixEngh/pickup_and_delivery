package operator

import "github.com/NixEngh/pickup_and_delivery/internal/solution"

type Operator interface {
	Apply(s *solution.Solution) int
}

type Insert interface {
	Insert(s *solution.Solution, calls []int)
}

type Removal interface {
	Choose(s *solution.Solution, n int) []int
}
