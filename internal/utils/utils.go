package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/NixEngh/pickup_and_delivery/internal/problem"
)

func FindIndices[T comparable](slice []T, values ...T) map[T][]int {
	indicesMap := make(map[T][]int)
    for index, element := range slice {
        for _, value := range values {
			if element == value {
				indicesMap[value] = append(indicesMap[value], index)
			}
		}
	}
	return indicesMap
}

func MoveElement[T comparable](s []T, source, destination int) {
	if source < 0 || source >= len(s) {
		return
	}

	elem := s[source]
	if source < destination {
		copy(s[source:], s[source+1:destination+1])
	} else if source > destination {
		copy(s[destination+1:source+1], s[destination:source])
	} else {
		return
	}

	s[destination] = elem
	return
}

// returns the indices such that solution[start_ind:end_ind] gives you the calls in the tour. If the length of the tour is 0, you get end_ind,end_ind
func GetTourIndices(solution []int, vehicleIndex int) []int {
	zeroIndexes := FindIndices(solution, 0)[0]
	start_ind := 0

	if vehicleIndex != 1 {
		start_ind = zeroIndexes[vehicleIndex-2] + 1
	}
	end_ind := zeroIndexes[vehicleIndex-1]
	if end_ind == 0 {
		return []int{end_ind, end_ind}
	}

	return []int{start_ind, end_ind}
}

// Returns a slice of the solution that represents the tour of a vehicle.
// If the slice is modified, the original solution will be modified as well
// Only includes the calls, not the 0's
func GetTour(solution []int, vehicleIndex int) []int {
	indices := GetTourIndices(solution, vehicleIndex)

	return solution[indices[0]:indices[1]]
}

func GetCallNodeTour(p *problem.Problem, solution []int, vehicleIndex int) []CallNode {
	tourIndices := GetTourIndices(solution, vehicleIndex)
	tour := make([]CallNode, 0)

	isDelivery := make(map[int]bool)

	for _, callIndex := range solution[tourIndices[0]:tourIndices[1]] {
		call := p.Calls[callIndex]
		timeWindow := call.PickupTimeWindow
		node := call.OriginNode
		operationTime := call.OriginTimeForVehicle[vehicleIndex]
		cost := call.OriginCostForVehicle[vehicleIndex]
		if isDelivery[callIndex] {
			timeWindow = call.DeliveryTimeWindow
			node = call.DestinationNode
			operationTime = call.DestinationTimeForVehicle[vehicleIndex]
			cost = call.DestinationCostForVehicle[vehicleIndex]
		}
		tour = append(tour, CallNode{
			CallIndex:     callIndex,
			TimeWindow:    timeWindow,
			Node:          node,
			IsDelivery:    isDelivery[callIndex],
			OperationTime: operationTime,
			Cost:          cost,
		})

		isDelivery[callIndex] = true
	}
	return tour
}

func PrintLoadingBar(current int, total int, steps int) {
	percentage := float64(current) / float64(total) * 100
	numberOfEquals := int(percentage / 100 * float64(steps))
	bar := strings.Repeat("=", numberOfEquals) + strings.Repeat(" ", steps-numberOfEquals)
	//only print when the percentage changes

	if current == 0 || int(percentage)%(100/steps) == 0 {
		fmt.Printf("\r[%s] %.0f%%", bar, percentage)
	}
}


type CSVTableRow struct {
	Algorithm    string
	AverageCost  float64
	BestCost     int
	Improvement  float64
	RunningTime  float64
	BestSolution []int
}

func WriteToCSV(directory string, name string, data []CSVTableRow) {

	file, err := os.Create(directory + "/" + name + ".csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// sort data by name index
	sort.Slice(data, func(i, j int) bool {
		i1 := strings.Split(data[i].Algorithm, "_")[0]
		j1 := strings.Split(data[j].Algorithm, "_")[0]
		return i1 < j1

	})

	writer.Write([]string{"Algorithm", "AverageCost", "BestCost", "Improvement", "RunningTime", "BestSolution"})

	for _, row := range data {
		bestSolution := fmt.Sprintf("%v", row.BestSolution)
		writer.Write([]string{row.Algorithm, fmt.Sprintf("%f", row.AverageCost), fmt.Sprintf("%d", row.BestCost), fmt.Sprintf("%f", row.Improvement), fmt.Sprintf("%f", row.RunningTime), bestSolution})
	}
}

func CreateResultsDirectory() string {

	if _, err := os.Stat("./data/results"); os.IsNotExist(err) {
		os.Mkdir("results", 0755)
	}

	t := time.Now()
	directory := fmt.Sprintf("./data/results/%s", t.Format("2006-01-02_15:04:05"))
	os.Mkdir(directory, 0755)

	return directory
}

func RunPythonScript(directory string) {
	file, err := os.OpenFile(directory+"/results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Capture standard error
	stderr := &bytes.Buffer{}

	cmd := exec.Command("python3", "./scripts/printlatest.py")
	cmd.Stdout = file
	cmd.Stderr = stderr // Set standard error to be captured

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Stderr:", stderr.String()) // Print any standard error output
		panic(err)
	}
}
