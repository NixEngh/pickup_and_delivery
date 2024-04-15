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
	s.vehicleCumulativeTimes = [][]int{
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
	result := s.CalulateTimeSlack(tour, vehicleIndex)

	if !equal(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func CreateRandomFeasible(s *Solution, steps int) {
    generator := rand.New(rand.NewSource(1))

	numberOfMoves := 100
	for i := 0; i < numberOfMoves; i++ {
		vehicleIndex := generator.Intn(s.Problem.NumberOfVehicles) + 1

		callNumber := generator.Intn(s.Problem.NumberOfCalls) + 1
		inds := FindIndices(s.Solution, callNumber, 0)
		inds = s.MoveCallToOutsource(callNumber, inds)

		validIndices := s.GetVehicleInsertionPoints(vehicleIndex, callNumber)
		if len(validIndices) == 0 {
			continue
		}
		s.InsertCall(callNumber, inds, validIndices[generator.Intn(len(validIndices))])
        if !s.Feasible() {
            panic(fmt.Sprintf("Something very wrong:\n  %v,\nvehicle: %d,\ncall: %d \ntour: %v", s.infeasibleReason,vehicleIndex,callNumber, GetTour(s.Solution, vehicleIndex)))
        }
	}
}

func TestFindFeasibleInsertion(t *testing.T) {
	problem, err := LoadProblem("./Data/Call_18_Vehicle_5.txt")
	if err != nil {
		t.Errorf("Error loading problem")
	}

	var correctFeasible, correctInfeasible int
	var falseFeasible, fakeInfeasible int
	for testIndex := 0; testIndex < 200; testIndex++ {
		var solution *Solution = problem.GenerateInitialSolution()
		CreateRandomFeasible(solution, 100)

		call := rand.Intn(problem.NumberOfCalls) + 1
		indices := FindIndices(solution.Solution, 0, call)

		solution.MoveCallToOutsource(call, indices)

		for vehicleIndex := 1; vehicleIndex < problem.NumberOfVehicles+1; vehicleIndex++ {
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
							t.Log(testSolution.VehicleCumulativeCapacities(vehicleIndex), " - capacity")
							SizeString := ""
							for _, call := range GetTour(testSolution.Solution, vehicleIndex) {
								SizeString += fmt.Sprintf("%d, ", problem.Calls[call].Size)
							}
							t.Log(SizeString, " - sizes")
							vehicle := problem.Vehicles[vehicleIndex]
							t.Log("Vehiclecapacity: ", vehicle.Capacity)
							t.Log()
							t.Log("Times: ")
							t.Log(solution.VehicleCumulativeTimes(vehicleIndex), " - originaltimes")
							t.Log(testSolution.VehicleCumulativeTimes(vehicleIndex), " - time")
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
							t.Log(testSolution.VehicleCumulativeCapacities(vehicleIndex), " - capacity")
							SizeString := ""
							for _, call := range GetTour(testSolution.Solution, vehicleIndex) {
								SizeString += fmt.Sprintf("%d, ", problem.Calls[call].Size)
							}
							t.Log(SizeString, " - sizes")
							vehicle := problem.Vehicles[vehicleIndex]
							t.Log("Vehiclecapacity: ", vehicle.Capacity)

							t.Log()
							t.Log("Times: ")
							t.Log(solution.VehicleCumulativeTimes(vehicleIndex), " - originaltimes")
							t.Log(solution.CalulateTimeSlack(tour, vehicleIndex), " - timeslack")

							UpperBounds := ""
							for _, call := range GetCallNodeTour(&problem, solution.Solution, vehicleIndex) {
								UpperBounds += fmt.Sprintf("%d, ", call.TimeWindow.UpperBound)
							}
							t.Log(UpperBounds, " - old upperbounds")
							t.Log(testSolution.VehicleCumulativeTimes(vehicleIndex), " - newTimes")

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

func TestSpecificFindFeasible(t *testing.T) {
	problem, _ := LoadProblem("./Data/Call_18_Vehicle_5.txt")
	solution := problem.GenerateInitialSolution()
	twelveI := InsertionPoint{
		pickupIndex: RelativeIndex{
			2,
			0,
		},
		deliveryIndex: RelativeIndex{
			2,
			0,
		},
	}
	sixteenI := InsertionPoint{
		pickupIndex: RelativeIndex{
			2,
			2,
		},
		deliveryIndex: RelativeIndex{
			2,
			2,
		},
	}
	// tour: 12 12 16 16
	inds := FindIndices(solution.Solution, 0, 12)
	solution.InsertCall(12, inds, twelveI)
	inds = FindIndices(solution.Solution, 0, 16)
	solution.InsertCall(16, inds, sixteenI)
	inds = FindIndices(solution.Solution, 0, 10)
	tour := GetCallNodeTour(&problem, solution.Solution, 2)
	intTour := GetTour(solution.Solution, 2)

	validIndices := solution.GetVehicleInsertionPoints(2, 10)
	t.Log(validIndices)
	for i := 0; i < len(tour)+1; i++ {
		for j := i; j < len(tour)+1; j++ {
			ts := solution.copy()
			insertion := InsertionPoint{
				pickupIndex:   RelativeIndex{2, i},
				deliveryIndex: RelativeIndex{2, j},
			}
			ts.InsertCall(10, inds, insertion)
			shouldBeFeasible := false
			for _, ins := range validIndices {
				if ins.pickupIndex.Index == insertion.pickupIndex.Index &&
					ins.deliveryIndex.Index == insertion.deliveryIndex.Index {
					shouldBeFeasible = true
				}
			}

			if shouldBeFeasible && !ts.Feasible() {
				t.Errorf("Infeasible solution marked as feasible\nReason: %s", ts.infeasibleReason)
				t.Log("Indices: i:", i, "j:", j)
				t.Log("Tour: ", intTour)
				t.Log("Times: ", solution.VehicleCumulativeTimes(2))
				t.Log("Timeslack: ", solution.CalulateTimeSlack(tour, 2))
				t.Log("NewTour: ", GetTour(ts.Solution, 2))
				t.Log("Newtimes: ", ts.VehicleCumulativeTimes(2))
			}
			if !shouldBeFeasible && ts.Feasible() {
				t.Log()
				t.Error("Feasible solution marked as infeasible")
				t.Log("Indices: i:", i, "j:", j)
				t.Log("Tour: ", intTour)
				t.Log("Times: ", solution.VehicleCumulativeTimes(2))
				t.Log("Timeslack: ", solution.CalulateTimeSlack(tour, 2))
				t.Log("NewTour: ", GetTour(ts.Solution, 2))
				t.Log("Newtimes: ", ts.VehicleCumulativeTimes(2))
			}
		}
	}
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
