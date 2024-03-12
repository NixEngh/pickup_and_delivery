package main

import (
	"fmt"
)

func main() {
	problem, _ := load_problem("../Data/Call_7_Vehicle_3.txt")

	fmt.Println(problem.GenerateInitialSolution())
	fmt.Println(problem.GenerateRandomSolution())
}

func RunExperiment(problem *Problem, algorithm func(*Problem) ([]int, int)) {

	costs := make([]int, 0)
	solutions := make([][]int, 0)
	for i := 0; i < 10; i++ {
		solution, cost := algorithm(problem)
		costs = append(costs, cost)
		solutions = append(solutions, solution)
	}

	average_solution := 0
	for _, cost := range costs {
        average_solution+=cost
	}
    average_solution /= len(costs)

    best_cost := costs[0]
    best_solution := solutions[0]

    for i, cost := range costs[1:] {
        if cost < best_cost {
            best_cost = cost
            best_solution = solutions[i]
        }
    }

    
    
}
