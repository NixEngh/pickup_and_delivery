package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func FindIndices[T comparable](slice []T, values ...T) map[T][]int {
	indicesMap := make(map[T][]int)
	for _, value := range values {
		for index, element := range slice {
			if element == value {
				indicesMap[value] = append(indicesMap[value], index)
			}
		}
	}
	return indicesMap
}

func MoveElement(s []int, source, destination int) []int {
	if source < 0 || source >= len(s) {
		return s
	}

	elem := s[source]
	if source < destination {
		copy(s[source:], s[source+1:destination+1])
	} else if source > destination {
		copy(s[destination+1:source+1], s[destination:source])
	} else {
		return s
	}

	s[destination] = elem
	return s
}

// Returns a slice of the solution that represents the tour of a vehicle.
// If the slice is modified, the original solution will be modified as well
func GetTour(solution []int, vehicleIndex int) []int {
	zeroIndexes := FindIndices(solution, 0)[0]
	start_ind := 0

	if vehicleIndex != 1 {
		start_ind = zeroIndexes[vehicleIndex-2] + 1
	}
	end_ind := zeroIndexes[vehicleIndex-1]

	return solution[start_ind:end_ind]
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

func LoadProblems(directory string) ([]Problem, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	problems := make([]Problem, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		problem, err := LoadProblem(directory + file.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		problems = append(problems, problem)
	}

	return problems, nil
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

	if _, err := os.Stat("results"); os.IsNotExist(err) {
		os.Mkdir("results", 0755)
	}

	t := time.Now()
	directory := fmt.Sprintf("results/%s", t.Format("2006-01-02_15:04:05"))
	os.Mkdir(directory, 0755)

	return directory
}

func runPythonScript(directory string) {
	file, err := os.OpenFile(directory+"/results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Capture standard error
	stderr := &bytes.Buffer{}

	cmd := exec.Command("python3", "printlatest.py")
	cmd.Stdout = file
	cmd.Stderr = stderr // Set standard error to be captured

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Stderr:", stderr.String()) // Print any standard error output
		panic(err)
	}
}
