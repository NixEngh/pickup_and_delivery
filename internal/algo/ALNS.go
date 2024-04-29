package algo

import (
	"github.com/NixEngh/pickup_and_delivery/internal/policy"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

func ALNS(operatorPolicy policy.OperatorPolicy, acceptor Acceptor, stopper Stopper) Algorithm {

	return func(problem *problem.Problem) (bestSolution *solution.Solution, bestCost int) {
		bestSolution = solution.GenerateInitialSolution(problem)

		S := bestSolution.Copy()
		var newS *solution.Solution

		for !stopper.CheckStop() {
			newS = S.Copy()
			operatorPolicy.Apply(S)

			if newS.Cost() < bestSolution.Cost() {
				bestSolution = S.Copy()
			}
			if acceptor.Accept(S, newS, bestSolution) {
				S = newS
			}
		}

		return bestSolution, bestSolution.Cost()
	}
}
