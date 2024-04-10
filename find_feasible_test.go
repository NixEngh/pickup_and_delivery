package main

import (
	"fmt"
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

func TestMoveRelativeToVehicle(t *testing.T) {
	problem, _ := LoadProblem("Data/Call_7_Vehicle_3.txt")
	solution := problem.GenerateInitialSolution()
	for i := 0; i < 100; i++ {
		call := rand.Intn(problem.NumberOfCalls) + 1
		indices := FindIndices(solution.Solution, call, 0)
		indices = solution.MoveCallToOutsource(call, indices)
		vehicle := rand.Intn(problem.NumberOfVehicles) + 1
		tour := GetCallNodeTour(&problem, solution.Solution, vehicle)

		movePickupTo := RelativeIndex{
			VehicleIndex: vehicle,
			Index:        0,
		}
		moveDeliveryTo := RelativeIndex{
			VehicleIndex: vehicle,
			Index:        1,
		}
		if len(tour) > 0 {
			movePickupTo.Index = rand.Intn(len(tour)) + 1
			a := movePickupTo.Index
			moveDeliveryTo.Index = rand.Intn(len(tour)+1-a) + a + 1
		}

		insertAt := InsertionPoint{
			pickupIndex:   movePickupTo,
			deliveryIndex: moveDeliveryTo,
		}

		solution.MoveCallToVehicle(call, indices, insertAt)
		//fmt.Println("After move: ", solution.Solution)

		newTour := GetCallNodeTour(&problem, solution.Solution, vehicle)

		if len(newTour) != len(tour)+2 {
			fmt.Println("len newTour: ", len(newTour))
			fmt.Println("len tour: ", len(tour))
			t.Errorf("Tour length not correct")
			panic("Tour length not correct")
		} else {
			if newTour[movePickupTo.Index].callIndex != call || newTour[moveDeliveryTo.Index].callIndex != call {
				t.Errorf("Call not moved correctly")
			}
		}
	}
}
func TestCalculateTimeSlack(t *testing.T) {
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
	problem, err := LoadProblem("./Data/Call_18_Vehicle_5.txt")
	if err != nil {
		t.Errorf("Error loading problem")
	}

	var correctFeasible, correctInfeasible int
	var falseFeasible, fakeInfeasible int
	for testIndex := 0; testIndex < 20 ; testIndex++ {
		var solution *Solution = problem.GenerateInitialSolution()
		numberOfMoves := 100
		for i := 0; i < numberOfMoves; i++ {
			vehicleIndex := rand.Intn(problem.NumberOfVehicles) + 1

			callNumber := rand.Intn(problem.NumberOfCalls) + 1
			inds := FindIndices(solution.Solution, callNumber, 0)
			inds = solution.MoveCallToOutsource(callNumber, inds)
			validIndices := solution.GetVehicleInsertionPoints(vehicleIndex, callNumber)
			if len(validIndices) == 0 {
				continue
			}
			solution.InsertCall(callNumber, inds, validIndices[rand.Intn(len(validIndices))])
		}

		call := rand.Intn(problem.NumberOfCalls) + 1
        call = 10
		indices := FindIndices(solution.Solution, 0, call)

        solution.MoveCallToOutsource(call, indices)

		for vehicleIndex := 1; vehicleIndex < problem.NumberOfVehicles+1; vehicleIndex++ {
            if vehicleIndex != 2 {continue}

			validIndices := solution.GetVehicleInsertionPoints(vehicleIndex, call)
			tour := GetCallNodeTour(&problem, solution.Solution, vehicleIndex)

			testSolution := solution.copy()
			for i := 0; i < len(tour)+1; i++ {
				for j := i; j < len(tour)+1; j++ {
                    inds := FindIndices(testSolution.Solution, call, 0)
                    inds = testSolution.MoveCallToOutsource(call, inds)

					comparisonPoint := InsertionPoint{
						pickupIndex:   RelativeIndex{VehicleIndex: vehicleIndex, Index: i},
						deliveryIndex: RelativeIndex{VehicleIndex: vehicleIndex, Index: j},
					}
                    testSolution.InsertCall(call, inds, comparisonPoint)

					if _, ok := slices.BinarySearchFunc(validIndices, comparisonPoint, func(a, t InsertionPoint) int {
						if t.pickupIndex.Index == a.pickupIndex.Index {
							return a.deliveryIndex.Index - t.deliveryIndex.Index
						}
						return a.pickupIndex.Index - t.pickupIndex.Index
					}); ok {
						if !testSolution.Feasible() {
							t.Log()
							t.Error("Infeasible solution marked as feasible")
							t.Log("Start solution: ", solution.Solution)
							t.Log("Infeasible solution: ", testSolution.Solution)
							t.Log("Vehicle: ", vehicleIndex)
							t.Log("Call: ", call)
							t.Log("valid indices: ", validIndices)
							t.Log("Current indexes: ", i, j)
							t.Log(GetTour(testSolution.Solution, vehicleIndex), " - tour")

							t.Log("Capacities: ")
							t.Log(testSolution.VehicleCumulativeCapacities[vehicleIndex], " - capacity")
							SizeString := ""
							for _, call := range GetTour(testSolution.Solution, vehicleIndex) {
								SizeString += fmt.Sprintf("%d, ", problem.Calls[call].Size)
							}
							t.Log(SizeString, " - sizes")
							vehicle := problem.Vehicles[vehicleIndex]
							t.Log("Vehiclecapacity: ", vehicle.Capacity)
                            t.Log()
							t.Log("Times: ")
							t.Log(solution.VehicleCumulativeTimes[vehicleIndex], " - originaltimes")
							t.Log(testSolution.VehicleCumulativeTimes[vehicleIndex], " - time")
							t.Log("infeasibleReason:", testSolution.infeasibleReason)
							t.Log("")
							fakeInfeasible++
						} else {
							correctFeasible++
						}
					} else {
						if testSolution.Feasible() {
                            t.Log("################33")
							t.Error("Feasible solution marked as infeasible")
							t.Log("Start solution: ", solution.Solution)
							t.Log("Feasible solution: ", testSolution.Solution)
							t.Log("Vehicle: ", vehicleIndex)
							t.Log("Call: ", call)
							t.Log("valid indices: ", validIndices)
							t.Log("Current index: ", i, j)
							t.Log(GetTour(testSolution.Solution, vehicleIndex), " - tour")

							t.Log("Capacities: ")
							t.Log(testSolution.VehicleCumulativeCapacities[vehicleIndex], " - capacity")
							SizeString := ""
							for _, call := range GetTour(testSolution.Solution, vehicleIndex) {
								SizeString += fmt.Sprintf("%d, ", problem.Calls[call].Size)
							}
							t.Log(SizeString, " - sizes")
							vehicle := problem.Vehicles[vehicleIndex]
							t.Log("Vehiclecapacity: ", vehicle.Capacity)

                            t.Log()
							t.Log("Times: ")
							t.Log(solution.VehicleCumulativeTimes[vehicleIndex], " - originaltimes")
                            t.Log(solution.CalulateTimeSlack(tour, vehicleIndex, 0), " - timeslack")

                            UpperBounds := ""
							for _, call := range GetCallNodeTour(&problem, solution.Solution, vehicleIndex) {
								UpperBounds += fmt.Sprintf("%d, ", call.TimeWindow.UpperBound)
							}
							t.Log(UpperBounds, " - old upperbounds")
							t.Log(testSolution.VehicleCumulativeTimes[vehicleIndex], " - newTimes")

							t.Log("")
							falseFeasible++
						} else {
							correctInfeasible++
						}
					}
				}
			}
		}
	}
	t.Log("Correctly marked as feasible: ", correctFeasible)
	t.Log("Correctly marked as infeasible: ", correctInfeasible)
	t.Log("Incorrectly marked as feasible", fakeInfeasible)
	t.Log("Incorrectly marked as infeasible", falseFeasible)
}

func TestGetCostImprovement(t *testing.T) {
	problem, _ := LoadProblem("./Data/Call_7_Vehicle_3.txt")
	solution := problem.GenerateInitialSolution()
	for i := 0; i < 100; i++ {
		call := rand.Intn(problem.NumberOfCalls) + 1
		indices := FindIndices(solution.Solution, 0, call)
		indices = solution.MoveCallToOutsource(call, indices)

		for vehicleIndex := 1; vehicleIndex <= problem.NumberOfVehicles; vehicleIndex++ {
			validIndices := solution.GetVehicleInsertionPoints(vehicleIndex, call)
			for _, insertionPoint := range validIndices {
				testSolution := solution.copy()
				testSolution.InsertCall(call, indices, insertionPoint)
				if testSolution.Cost()-solution.Cost() != insertionPoint.costDiff {
					t.Log("Originalcost: ", solution.Cost())
					t.Log("Cost of new solution: ", testSolution.Cost())

					t.Log("Real diff: ", solution.Cost()-testSolution.Cost())
					t.Log("Stored Costdiff: ", insertionPoint.costDiff)
					t.Error("Not equal")
				}
			}
		}
	}
}
