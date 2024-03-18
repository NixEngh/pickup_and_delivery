package main

import (
	"fmt"
	"sort"
	"time"
)

func main() {
	problems, err := LoadProblems("./Data/")
	if err != nil {
		fmt.Println(err)
		return
	}

	sort.Slice(problems, func(i, j int) bool {
		return problems[i].Name < problems[j].Name
	})

    algorithms := map[string]algorithm{"0_Random_Search": RandomSearch, "1_Local_Search": LocalSearch, "2_Simulated_Annealing": SimulatedAnnealing}
    //algorithms := map[string]algorithm{"Random_Search": RandomSearch}

	directory := CreateResultsDirectory()
	for _, problem := range problems {
		rows := make([]CSVTableRow, 0)
		for name, algorithm := range algorithms {
			row := RunExperiment(&problem, name, algorithm)
			rows = append(rows, row)
		}
		WriteToCSV(directory, problem.Name, rows)
	}
    runPythonScript(directory)
}

func RunExperiment(problem *Problem, algorithmName string, algorithm func(*Problem) ([]int, int)) CSVTableRow {
	costs := make([]int, 0)
	solutions := make([][]int, 0)

	fmt.Println("Running experiment for problem: ", problem.Name)

	start := time.Now()
	for i := 0; i < 10; i++ {
		solution, cost := algorithm(problem)
		costs = append(costs, cost)
		solutions = append(solutions, solution)
	}
	elapsed := time.Since(start)
	averageTime := elapsed.Seconds() / 10

	var averageCost float64 = 0
	for _, cost := range costs {
		averageCost += float64(cost)
	}
	averageCost /= float64(len(costs))

	bestCost := costs[0]
	bestSolution := solutions[0]

	for i, cost := range costs[1:] {
		if cost < bestCost {
			bestCost = cost
			bestSolution = solutions[i]
		}
	}

	improvement := 100 * (averageCost - float64(bestCost)) / averageCost

	return CSVTableRow{algorithmName, averageCost, bestCost, improvement, averageTime, bestSolution}
}
