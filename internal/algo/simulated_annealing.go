package algo

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
)
func SimulatedAnnealing(operatorPolicy operator.OperatorPolicy) Algorithm {
	return func(problem *problem.Problem) (bestSolution *solution.Solution, bestCost int) {
		finalTemperature := 0.1

		bestSolution = solution.GenerateInitialSolution(problem)
		bestCost = bestSolution.Cost()

		incubent := bestSolution.Copy()

		deltas := make([]int, 0)
		fmt.Println("Simulated Annealing for operator policy: ", operatorPolicy.Name())
		fmt.Println("Estimating initial temperature")

		for i := 0; i < 100; i++ {
			utils.PrintLoadingBar(i, 100, 50)
			neighbor := incubent.Copy()
            operatorPolicy.Apply(neighbor)

			neighborCost := neighbor.Cost()
			deltaE := neighborCost - incubent.Cost()

			if !neighbor.Feasible() {
				continue
			}

			if deltaE < 0 {
				incubent = neighbor
				if neighborCost < bestCost {
					bestSolution = neighbor
					bestCost = neighborCost
				}
				continue
			}

			if rand.Float64() < 0.8 {
				incubent = neighbor
			}

			deltas = append(deltas, deltaE)
		}
		fmt.Println("\nStarting annealing")

		deltaSum := 0
		for _, delta := range deltas {
			deltaSum += delta
		}
		deltaAvg := float64(deltaSum) / float64(len(deltas))

		T0 := -deltaAvg / math.Log(0.8)

		alpha := math.Pow(finalTemperature/T0, 1/9900)

		T := T0

		feasibleCount := 0
		for i := 0; i < 9900; i++ {
			utils.PrintLoadingBar(i, 9900, 50)
			neighbor := incubent.Copy()
			operatorPolicy.Apply(neighbor)

			deltaE := neighbor.Cost() - incubent.Cost()

			if !neighbor.Feasible() {
				T *= alpha
				continue
			}
			feasibleCount++

			if deltaE < 0 {
				incubent = neighbor
				if neighbor.Cost() < bestCost {
					bestSolution = neighbor
					bestCost = neighbor.Cost()
				}
			} else if rand.Float64() < math.Exp(-float64(deltaE)/T) {
				incubent = neighbor
			}

			T *= alpha
		}

		fmt.Println()
		fmt.Println("Feasible solutions found:", feasibleCount)

		return bestSolution, bestCost
	}
}
