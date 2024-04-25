package operator

import (
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type Operator interface {
	Apply(s *solution.Solution) int
}

