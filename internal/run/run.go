package run

import (
	"fmt"
	"time"

	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
)

func Run(algorithms map[string]algo.Algorithm, problems []*problem.Problem) {

	directory := utils.CreateResultsDirectory()
	for _, p := range problems {
		rows := make([]utils.CSVTableRow, 0)
		for name, algorithm := range algorithms {
			row := RunExperiment(p, name, algorithm)
			rows = append(rows, row)
		}
		utils.WriteToCSV(directory, p.Name, rows)
	}
	utils.RunPythonScript(directory)
}

func RunExperiment(problem *problem.Problem, algorithmName string, algorithm algo.Algorithm) utils.CSVTableRow {
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

	initialCost := solution.GenerateInitialSolution(problem).Cost()
	improvement := 100 * (float64(initialCost) - float64(bestCost)) / float64(initialCost)

	return utils.CSVTableRow{
		Algorithm:    algorithmName,
		AverageCost:  averageCost,
		BestCost:     bestCost,
		Improvement:  improvement,
		RunningTime:  averageTime,
		BestSolution: bestSolution,
	}
}
