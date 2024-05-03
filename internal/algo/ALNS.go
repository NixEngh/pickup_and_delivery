package algo

import (
	"fmt"
	"math"

	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/policy"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

func ALNS(operatorPolicy policy.OperatorPolicy, acceptor Acceptor, stopper Stopper) Algorithm {

	escape := operator.NewCombineOperator(operator.NewRemoveRandom(5), operator.NewInsertGreedy(), "Escape")

	return func(problem *problem.Problem) (bestSolution *solution.Solution, bestCost int) {
		bestSolution = solution.GenerateInitialSolution(problem)

		S := bestSolution.Copy()
		var newS *solution.Solution

		var iterationsSinceNewBest int
		var allowedIterations int = 500
		//var numberOfEscapes int = int(math.Round(math.Sqrt(float64(problem.NumberOfCalls))))
		var numberOfEscapes int = int(math.Sqrt(float64(problem.NumberOfCalls)))

		fmt.Println("ALNS for operator policy: ", operatorPolicy.Name())

		for !stopper.CheckStop() {

			newS = S.Copy()
			operatorPolicy.Apply(S)

			if newS.Cost() < bestSolution.Cost() {
				bestSolution = S.Copy()
				iterationsSinceNewBest = 0
			}
			if acceptor.Accept(S, newS, bestSolution) {
				S = newS
			}
			if iterationsSinceNewBest > allowedIterations {
				S = bestSolution.Copy()
				for i := 0; i < numberOfEscapes; i++ {
					escape.Apply(S)
				}
				iterationsSinceNewBest = 0
			}
		}
		stopper.Reset()
		fmt.Println()
		fmt.Println("ALNS finished for operator policy: ", problem.Name)

		return bestSolution, bestSolution.Cost()
	}
}
