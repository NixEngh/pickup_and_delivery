package main

import (
	"fmt"
	"time"
    "os"
)

func main() {
    // Get all text files from ./Data/
    files, err := os.ReadDir("./Data/")
    if err != nil {
        panic(err)
    }

    problems := make([]Problem, 0)

    for _, file := range files {
        if file.IsDir() {
            continue
        }
        problem, err := load_problem("./Data/" + file.Name())
        if err != nil {
            fmt.Println(err)
            continue
        }
        problems = append(problems, problem)
        fmt.Println(problem.Name)
    }

	problem, _ := load_problem("../Data/Call_7_Vehicle_3.txt")

	fmt.Println(problem.GenerateInitialSolution())
	fmt.Println(problem.GenerateRandomSolution())
}

func RunExperiment(problem *Problem, algorithm func(*Problem) ([]int, int)) CSVTableRow {
	costs := make([]int, 0)
	solutions := make([][]int, 0)
	// Time the ten iterations
	start := time.Now()

	for i := 0; i < 10; i++ {
		solution, cost := algorithm(problem)
		costs = append(costs, cost)
		solutions = append(solutions, solution)
	}
	elapsed := time.Since(start)
	averageTime := float64(elapsed / 10)

	var averageCost float64 = 0
	for _, cost := range costs {
		averageCost += float64(cost)
	}
	averageCost /= float64(len(costs))

	bestCost := costs[0]
	best_solution := solutions[0]

	for i, cost := range costs[1:] {
		if cost < bestCost {
			bestCost = cost
			best_solution = solutions[i]
		}
	}

	improvement := 100 * (averageCost - float64(bestCost)) / averageCost

	return CSVTableRow{problem.Name, averageCost, bestCost, improvement, averageTime, best_solution}
}
