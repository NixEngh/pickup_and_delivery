package main

import (
	"math/rand"
	"slices"
	"testing"
)

func equal(a, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCalculateTimeSlack(t *testing.T) {
	// Testsolution:
	// [1, 2, 1, 2, 0, 3, 3, 0]
	// Time
	// [1, 3, 4, 6]
	s := Solution{}
	s.VehicleCumulativeTimes = [][]int{
		{1, 3, 4, 6},
		{0, 2},
	}

	tour := []CallNode{
		{TimeWindow: TimeWindow{LowerBound: 2, UpperBound: 3}},
		{TimeWindow: TimeWindow{LowerBound: 2, UpperBound: 5}},
		{TimeWindow: TimeWindow{LowerBound: 2, UpperBound: 5}},
		{TimeWindow: TimeWindow{LowerBound: 7, UpperBound: 9}},
	}
	expected := []int{2, 1, 1, 3}

	vehicleIndex := 0
	startIndex := 0
	result := s.CalulateTimeSlack(tour, vehicleIndex, startIndex)

	if !equal(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func TestFindFeasibleInsertion(t *testing.T) {
	problem, err := LoadProblem("Data/Call_7_Vehicle_3.txt")
	if err != nil {
		t.Errorf("Error loading problem")
	}

	for testIndex := 0; testIndex < 20; testIndex++ {
		var solution *Solution
		for {
			solution = problem.GenerateRandomSolution()
			if solution.Feasible() {
				break
			}
		}
		call := rand.Intn(problem.NumberOfCalls) + 1
        t.Log(call)
		indices := FindIndices(solution.Solution, 0, call)
		zeroIndices := indices[0]
		t.Log(solution.Solution)
		if indices[call][0] < zeroIndices[len(zeroIndices)-1] {
			solution.MoveInSolution(indices[call][1], zeroIndices[len(zeroIndices)-1])
			solution.MoveInSolution(indices[call][0], zeroIndices[len(zeroIndices)-1])
		}
		t.Log(solution.Solution)

        zeroIndices = FindIndices(solution.Solution, 0)[0]

        for vehicleIndex := 1; vehicleIndex < problem.NumberOfVehicles+1; vehicleIndex++ {
			validIndices := solution.GetVehicleInsertionPoints(vehicleIndex, call)
            tour := GetCallNodeTour(&problem, solution.Solution, vehicleIndex)

            callIndices := FindIndices(solution.Solution, call)[call]
            pickupIndex := callIndices[0]
            deliveryIndex := callIndices[1]

            testSolution := solution.copy()
            for i := 0; i<len(tour)+1 ; i++ {
                relative1 := RelativeIndex{
                    VehicleIndex: vehicleIndex,
                    Index: i,
                }
                testSolution.MoveRelativeToVehicle(pickupIndex, relative1)
                pickupIndex = relative1.toAbsolute(FindIndices(testSolution.Solution, 0)[0])
                for j := i+1; j<len(tour)+2 ; j++ {
                    relative2 := RelativeIndex{
                        VehicleIndex: vehicleIndex,
                        Index: j,
                    }
                    testSolution.MoveRelativeToVehicle(deliveryIndex, relative2)
                    deliveryIndex = relative2.toAbsolute(FindIndices(testSolution.Solution, 0)[0])
					comparisonPoint := InsertionPoint{
						pickupIndex:   relative1,
						deliveryIndex: relative2,
					}
					if _, ok := slices.BinarySearchFunc(validIndices, comparisonPoint, func(a, t InsertionPoint) int {

						if t.pickupIndex.Index == a.pickupIndex.Index {
							return t.deliveryIndex.Index - a.deliveryIndex.Index
						}
						return t.pickupIndex.Index - a.pickupIndex.Index
					}); ok {
						if !testSolution.Feasible() {
							t.Error("Infeasible solution marked as feasible")
						}
					} else {
						if testSolution.Feasible() {
                            t.Log(testSolution.Solution)
							t.Error("Feasible solution omitted")
						}
					}
                }
            }
        }
	}
}
