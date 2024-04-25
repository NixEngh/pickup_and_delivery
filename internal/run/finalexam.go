package run

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/policy"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
)


func RunFinalExam(problems []*problem.Problem) {
    algorithms := map[string]algo.Algorithm{
        "AlwaysOptimal": AlwaysOptimal(),
    }

    Run(algorithms, problems)

}

func AlwaysOptimal() algo.Algorithm {
    operators := []operator.Operator{
        operator.PlaceOptimally{},
    }
    policy := policy.NewLecturePolicy(50, 0.1, operators)
    return algo.SimulatedAnnealing(policy)
}
