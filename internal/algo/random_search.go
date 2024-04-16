package algo

import (
	"fmt"

	"github.com/NixEngh/pickup_and_delivery/internal/problem"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
)

func RandomSearch(problem *problem.Problem) (bestSolution *solution.Solution, bestCost int) {
	bestSolution = solution.GenerateInitialSolution(problem)
	bestCost = bestSolution.Cost()

	for i := 0; i < 10000; i++ {
		utils.PrintLoadingBar(i, 10000, 50)
		solution := solution.GenerateInitialSolution(problem)

		if solution.Feasible() && solution.Cost() < bestCost {
			bestSolution = solution
			bestCost = solution.Cost()
		}
	}
	fmt.Println()
	return bestSolution, bestCost
}
