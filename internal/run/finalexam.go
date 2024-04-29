package run

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

func RunFinalExam(problems []*problem.Problem) {
	algorithms := map[string]algo.Algorithm{
		"AlwaysOptimal": AlwaysOptimal(),
	}

	Run(algorithms, problems)

}

func AlwaysOptimal() algo.Algorithm {
	return func(problem *problem.Problem) (bestSolution *solution.Solution, bestCost int) {
		return
	}
}
