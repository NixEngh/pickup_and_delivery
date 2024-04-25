package run

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/policy"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
)

func RunAssignment5(problems []*problem.Problem) {
    algorithms := map[string]algo.Algorithm{
        "Adaptive": Adaptive(),
    }

    Run(algorithms, problems)
}

func Adaptive() algo.Algorithm {
	operators := []operator.Operator{
		operator.PlaceOptimally{},
		operator.PlaceOptimallyInRandomVehicle{},
		operator.PlaceRandomly{},
		operator.PlaceFiveCallsRandomly{},
	}
	policy := policy.NewLecturePolicy(50, 0.1, operators)
	return algo.SimulatedAnnealing(policy)
}
