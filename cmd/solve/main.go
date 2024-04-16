package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/assignment"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
)

func main() {
    algorithms := map[string]algo.Algorithm{
        "1_Adaptive": assignment.Adaptive(),
    }

    Run(algorithms)
}

func Run(algorithms map[string]algo.Algorithm) {
	problems, err := problem.LoadProblems("./data/input/")
	if err != nil {
		fmt.Println(err)
		return
	}

	sort.Slice(problems, func(i, j int) bool {
		name1 := strings.Split(problems[i].Name, "_")[1]
		name2 := strings.Split(problems[j].Name, "_")[1]

		return name1 < name2
	})

	improvements := make(map[string]float64)

	directory := utils.CreateResultsDirectory()
	for _, problem := range problems {
		rows := make([]utils.CSVTableRow, 0)
		for name, algorithm := range algorithms {
			row := RunExperiment(&problem, name, algorithm)
			rows = append(rows, row)
			improvements[name+"-"+problem.Name] += row.Improvement
		}
		utils.WriteToCSV(directory, problem.Name, rows)
	}
	utils.RunPythonScript(directory)

	for name, improvement := range improvements {
		fmt.Println(name, improvement)
	}

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
        Algorithm:  algorithmName,
        AverageCost:  averageCost,
        BestCost:    bestCost,
        Improvement: improvement,
        RunningTime: averageTime,
        BestSolution: bestSolution,
    }
}
