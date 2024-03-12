package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func findIndices(slice []int, value int) []int {
	var indices []int
	for index, element := range slice {
		if element == value {
			indices = append(indices, index)
		}
	}
	return indices
}

func getTour(solution []int, vehicleIndex int) []int {
	zeroIndexes := findIndices(solution, 0)
	start_ind := 0

	if vehicleIndex != 1 {
		start_ind = zeroIndexes[vehicleIndex-2] + 1
	}
	end_ind := zeroIndexes[vehicleIndex-1]

	return solution[start_ind:end_ind]
}

type CSVTableRow struct {
	Algorithm    string
	AverageCost  float64
	BestCost     int
	Improvement  float64
	RunningTime  float64
	BestSolution []int
}

func WriteToCSV(filename string, data []CSVTableRow) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Algorithm", "AverageCost", "BestCost", "Improvement", "RunningTime", "BestSolution"})

	// Write data
	for _, row := range data {
		bestSolution := fmt.Sprintf("%v", row.BestSolution)
		writer.Write([]string{row.Algorithm, fmt.Sprintf("%f", row.AverageCost), fmt.Sprintf("%f", row.BestCost), fmt.Sprintf("%f", row.Improvement), fmt.Sprintf("%f", row.RunningTime), bestSolution})
	}
}
