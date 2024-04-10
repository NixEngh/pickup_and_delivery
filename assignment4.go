package main

import (
	"fmt"
	"math"
	"math/rand"
)

func Assignment_4_1(problem *Problem) (bestSolution *Solution, bestCost int) {
	operator := PlaceOptimally{}
	finalTemperature := 0.1

	bestSolution = problem.GenerateInitialSolution()
	bestCost = bestSolution.Cost()

	incubent := bestSolution.copy()

	deltas := make([]int, 0)
	fmt.Println("Simulated Annealing")
	fmt.Println("Estimating initial temperature")

	for i := 0; i < 100; i++ {
		PrintLoadingBar(i, 100, 50)
		neighbor := incubent.copy()
		operator.apply(neighbor)

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
		PrintLoadingBar(i, 9900, 50)
		neighbor := incubent.copy()
		operator.apply(neighbor)

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

