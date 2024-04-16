package algo

import (
	"fmt"

	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
)

func LocalSearch(problem *problem.Problem) (bestSolution *solution.Solution, bestCost int) {
	operator := operator.OldOneReinsert{}

	bestSolution = solution.GenerateInitialSolution(problem)
	bestCost = bestSolution.Cost()

	feasibleCount := 0
	for i := 0; i < 10000; i++ {
		utils.PrintLoadingBar(i, 10000, 50)
		operator.Apply(bestSolution)

		if bestSolution.Feasible() {
			feasibleCount++
			if bestSolution.Cost() < bestCost {
				bestCost = bestSolution.Cost()
			}
		}
	}
	fmt.Println()
	fmt.Println("Feasible solutions found:", feasibleCount)

	return bestSolution, bestCost
}
