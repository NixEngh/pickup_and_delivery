package run

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/policy"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
)

func RunAssignment3(problems []*problem.Problem) {
    algorithms := map[string]algo.Algorithm{
        "Local Search": LocalSearch(),
    }

    Run(algorithms, problems)
}

func LocalSearch() algo.Algorithm {
	policy := policy.NewChooseRandomOperator(
		[]policy.OperatorScore{
			{Operator: operator.OldOneReinsert{}, Probability: 1},
		},
		"Old OneReinsert",
	)
	return algo.SimulatedAnnealing(policy)
}
