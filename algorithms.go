package main

import (
	"fmt"
	"math"
	"math/rand"
)

func RandomSearch(problem *Problem) ([]int, int) {
	bestSolution := problem.GenerateInitialSolution()
	bestCost := problem.CostFunction(bestSolution)

	for i := 0; i < 10000; i++ {
		PrintLoadingBar(i, 10000, 50)
		solution := problem.GenerateRandomSolution()
		cost := problem.CostFunction(solution)
		if problem.IsFeasible(solution) && cost < bestCost {
			bestSolution = solution
			bestCost = cost
		}
	}
	fmt.Println()
	return bestSolution, bestCost
}

func LocalSearch(problem *Problem) ([]int, int) {
	bestSolution := problem.GenerateInitialSolution()
	bestCost := problem.CostFunction(bestSolution)

	for i := 0; i < 10000; i++ {
		PrintLoadingBar(i, 10000, 50)
		OneReinsert(problem, bestSolution)

		if problem.IsFeasible(bestSolution) {
			cost := problem.CostFunction(bestSolution)
			if cost < bestCost {
				bestCost = cost
			}
		}
	}
	fmt.Println()

	return bestSolution, bestCost
}

func SimulatedAnnealing(problem *Problem) ([]int, int) {
	finalTemperature := 0.1

	bestSolution := problem.GenerateInitialSolution()
	bestCost := problem.CostFunction(bestSolution)

	incubent := make([]int, len(bestSolution))
	copy(incubent, bestSolution)

	deltas := make([]int, 0)
	fmt.Println("Simulated Annealing")
	fmt.Println("Estimating initial temperature")
	for i := 0; i < 100; i++ {
		PrintLoadingBar(i, 100, 50)
		neighbor := make([]int, len(incubent))
		copy(neighbor, incubent)
		OneReinsert(problem, neighbor)

		neighborCost := problem.CostFunction(neighbor)
		deltaE := neighborCost - bestCost

		if !problem.IsFeasible(neighbor) {
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

	for i := 0; i < 9900; i++ {
		PrintLoadingBar(i, 9900, 50)
		neighbor := make([]int, len(incubent))
        copy(neighbor, incubent)

		OneReinsert(problem, neighbor)

		neighborCost := problem.CostFunction(neighbor)

		deltaE := neighborCost - bestCost
		isFeasible := problem.IsFeasible(neighbor)

		if !isFeasible {
			T *= alpha
			continue
		}

		if deltaE < 0 {
			incubent = neighbor
			if neighborCost < bestCost {
				bestSolution = neighbor
				bestCost = neighborCost
			}
		} else if rand.Float64() < math.Exp(-float64(deltaE)/T) {
			incubent = neighbor
		}

		T *= alpha
	}

	fmt.Println()

	return bestSolution, bestCost
}
