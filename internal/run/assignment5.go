package run

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
)

func RunAssignment5(problems []*problem.Problem) {
	algorithms := map[string]algo.Algorithm{}

	Run(algorithms, problems)
}
