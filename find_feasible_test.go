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
		callIndices := indices[call]
		zeroIndices := indices[0]
		callIndices = solution.MoveCallToOutsource(callIndices, zeroIndices)
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
		//fmt.Println("Call: ", call)
		//fmt.Println("Callindices: ", callIndices)
		//fmt.Println(solution.Solution)
		//fmt.Println("vehicle", vehicle)
		//fmt.Println("Move to", movePickupTo.Index, moveDeliveryTo.Index)

		insertAt := InsertionPoint{
			pickupIndex:   movePickupTo,
			deliveryIndex: moveDeliveryTo,
		}

		zeroIndices = FindIndices(solution.Solution, 0)[0]
		solution.MoveCallToVehicle(callIndices, zeroIndices, insertAt)
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
	generator := rand.New(rand.NewSource(1))

	for testIndex := 0; testIndex < 1; testIndex++ {
		var solution *Solution
		for {
			solution = problem.GenerateRandomSolution()
			if solution.Feasible() {
				break
			}
		}
		call := generator.Intn(problem.NumberOfCalls) + 1
		//t.Log("Call: ", call)
		indices := FindIndices(solution.Solution, 0, call)
		zeroIndices := indices[0]
		//t.Log("Solution: ", solution.Solution)

		if indices[call][0] < zeroIndices[len(zeroIndices)-1] {
			solution.MoveInSolution(indices[call][1], zeroIndices[len(zeroIndices)-1])
			solution.MoveInSolution(indices[call][0], zeroIndices[len(zeroIndices)-1])
			//t.Log("Moved call: ", solution.Solution)
		}

		zeroIndices = FindIndices(solution.Solution, 0)[0]

		for vehicleIndex := 1; vehicleIndex < problem.NumberOfVehicles+1; vehicleIndex++ {
			validIndices := solution.GetVehicleInsertionPoints(vehicleIndex, call)
			tour := GetCallNodeTour(&problem, solution.Solution, vehicleIndex)

			callIndices := FindIndices(solution.Solution, call)[call]
			pickupIndex := callIndices[0]
			deliveryIndex := callIndices[1]

			testSolution := solution.copy()
			//t.Log("VehicleIndex: ", vehicleIndex)
			for i := 0; i < len(tour)+1; i++ {
				relative1 := RelativeIndex{
					VehicleIndex: vehicleIndex,
					Index:        i,
				}
				//t.Log("i: ", i)
				testSolution.MoveRelativeToVehicle(pickupIndex, relative1)
				//t.Log("solution after pickup moved: ", testSolution.Solution)
				pickupIndex = relative1.toAbsolute(FindIndices(testSolution.Solution, 0)[0])
				for j := i + 1; j < len(tour)+2; j++ {
					//t.Log("j: ", j)
					relative2 := RelativeIndex{
						VehicleIndex: vehicleIndex,
						Index:        j,
					}

					testSolution.MoveRelativeToVehicle(deliveryIndex, relative2)
					//t.Log("solution after delivery moved: ", testSolution.Solution)
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
							t.Log("Call: ", call)
							t.Log("Infeasible solution: ", testSolution.Solution)
							t.Log(GetTour(testSolution.Solution, vehicleIndex), " - tour")
							t.Log(testSolution.VehicleCumulativeCapacities[vehicleIndex], " - capacity")
							t.Log(testSolution.VehicleCumulativeTimes[vehicleIndex], " - time")
							SizeString := ""
							for _, call := range GetTour(testSolution.Solution, vehicleIndex) {
								SizeString += fmt.Sprintf("%d, ", problem.Calls[call].Size)
							}
                            vehicle := problem.Vehicles[vehicleIndex]
                            t.Log("Vehiclecapacity: ", vehicle.Capacity)
                            t.Log(SizeString, " - sizes")
							t.Log("infeasibleReason:", testSolution.infeasiblereason)
							t.Log("")
							t.Error("Infeasible solution marked as feasible")
						}
					} else {
						if testSolution.Feasible() {
							t.Log("Feasible solution: ", testSolution.Solution)
							t.Error("Feasible solution omitted")
						}
					}
				}
			}
		}
	}
}
