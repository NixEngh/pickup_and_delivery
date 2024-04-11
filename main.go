package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func main() {
	problems, err := LoadProblems("./Data/")
	if err != nil {
		fmt.Println(err)
		return
	}

	sort.Slice(problems, func(i, j int) bool {
		name1 := strings.Split(problems[i].Name, "_")[1]
		name2 := strings.Split(problems[j].Name, "_")[1]

		return name1 < name2
	})

	algorithms := map[string]algorithm{
        //"0_Random_Search":       RandomSearch,
        //"1_Local_Search":        LocalSearch,
        //"2_Simulated_Annealing": Assignment3(),
        "3_Equal": EqualProbability(),
        "4_Moderate": Moderate(),
        "5_Adventurous": Adventurous(),
        "6_Intense": Intense(),
    }

	improvements := make(map[string]float64)

	directory := CreateResultsDirectory()
	for _, problem := range problems {
		rows := make([]CSVTableRow, 0)
		for name, algorithm := range algorithms {
			row := RunExperiment(&problem, name, algorithm)
			rows = append(rows, row)
			improvements[name+"-"+problem.Name] += row.Improvement
		}
		WriteToCSV(directory, problem.Name, rows)
	}
	runPythonScript(directory)

	for name, improvement := range improvements {
		fmt.Println(name, improvement)
	}
}

func RunExperiment(problem *Problem, algorithmName string, algorithm algorithm) CSVTableRow {
	costs := make([]int, 0)
	solutions := make([][]int, 0)

	fmt.Println("Running experiment for problem: ", problem.Name)

	start := time.Now()
	for i := 0; i < 10; i++ {
		solution, cost := algorithm(problem)
		costs = append(costs, cost)
		solutions = append(solutions, solution.Solution)
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

	initialCost := problem.GenerateInitialSolution().Cost()
	improvement := 100 * (float64(initialCost) - float64(bestCost)) / float64(initialCost)

	return CSVTableRow{algorithmName, averageCost, bestCost, improvement, averageTime, bestSolution}
}
