package assignment

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
)

func Assignment3() algo.Algorithm {
	policy := operator.NewChooseRandomOperator(
		[]operator.OperatorScore{
			{Operator: operator.PlaceOptimally{}, Probability: 1},
		},
		"Old OneReinsert",
	)
	return algo.SimulatedAnnealing(policy)
}
