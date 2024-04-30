package solution

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/NixEngh/pickup_and_delivery/internal/problem"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
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
	problem, _ := problem.LoadProblem("Data/Call_7_Vehicle_3.txt")
	solution := GenerateInitialSolution(problem)
	for i := 0; i < 100; i++ {
		call := rand.Intn(problem.NumberOfCalls) + 1
		indices := utils.FindIndices(solution.Solution, call)
		indices = solution.MoveCallToOutsource(call, indices)
		vehicle := rand.Intn(problem.NumberOfVehicles) + 1
		tour := utils.GetCallNodeTour(problem, solution.Solution, vehicle)

		movePickupTo := utils.RelativeIndex{
			VehicleIndex: vehicle,
			Index:        0,
		}
		moveDeliveryTo := utils.RelativeIndex{
			VehicleIndex: vehicle,
			Index:        1,
		}
		if len(tour) > 0 {
			movePickupTo.Index = rand.Intn(len(tour)) + 1
			a := movePickupTo.Index
			moveDeliveryTo.Index = rand.Intn(len(tour)+1-a) + a + 1
		}

		insertAt := utils.InsertionPoint{
			PickupIndex:   movePickupTo,
			DeliveryIndex: moveDeliveryTo,
		}

		solution.MoveCallToVehicle(call, indices, insertAt)
		//fmt.Println("After move: ", solution.Solution)

		newTour := utils.GetCallNodeTour(problem, solution.Solution, vehicle)

		if len(newTour) != len(tour)+2 {
			fmt.Println("len newTour: ", len(newTour))
			fmt.Println("len tour: ", len(tour))
			t.Errorf("Tour length not correct")
			panic("Tour length not correct")
		} else {
			if newTour[movePickupTo.Index].CallIndex != call || newTour[moveDeliveryTo.Index].CallIndex != call {
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

	tour := []utils.CallNode{
		{TimeWindow: problem.TimeWindow{LowerBound: 2, UpperBound: 3}},
		{TimeWindow: problem.TimeWindow{LowerBound: 2, UpperBound: 5}},
		{TimeWindow: problem.TimeWindow{LowerBound: 2, UpperBound: 5}},
		{TimeWindow: problem.TimeWindow{LowerBound: 7, UpperBound: 9}},
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
		inds := utils.FindIndices(s.Solution, callNumber)
		inds = s.MoveCallToOutsource(callNumber, inds)

		validIndices := s.GetVehicleInsertionPoints(vehicleIndex, callNumber)
		if len(validIndices) == 0 {
			continue
		}
		s.InsertCall(callNumber, inds, validIndices[generator.Intn(len(validIndices))])
		if !s.Feasible() {
			panic(fmt.Sprintf("Something very wrong:\n  %v,\nvehicle: %d,\ncall: %d \ntour: %v", s.infeasibleReason, vehicleIndex, callNumber, utils.GetTour(s.Solution, vehicleIndex)))
		}
	}
}

func TestFindFeasibleInsertion(t *testing.T) {
	problem, err := problem.LoadProblem("./Data/Call_18_Vehicle_5.txt")
	if err != nil {
		t.Errorf("Error loading problem")
	}

	var correctFeasible, correctInfeasible int
	var falseFeasible, fakeInfeasible int
	for testIndex := 0; testIndex < 200; testIndex++ {
		var solution *Solution = GenerateInitialSolution(problem)
		CreateRandomFeasible(solution, 100)

		call := rand.Intn(problem.NumberOfCalls) + 1
		indices := utils.FindIndices(solution.Solution, call)

		solution.MoveCallToOutsource(call, indices)

		for vehicleIndex := 1; vehicleIndex < problem.NumberOfVehicles+1; vehicleIndex++ {
			validIndices := solution.GetVehicleInsertionPoints(vehicleIndex, call)

			tour := utils.GetCallNodeTour(problem, solution.Solution, vehicleIndex)

			testSolution := solution.Copy()
			for i := 0; i < len(tour)+1; i++ {
				for j := i; j < len(tour)+1; j++ {
					inds := utils.FindIndices(testSolution.Solution, call)
					inds = testSolution.MoveCallToOutsource(call, inds)

					comparisonPoint := utils.InsertionPoint{
						PickupIndex:   utils.RelativeIndex{VehicleIndex: vehicleIndex, Index: i},
						DeliveryIndex: utils.RelativeIndex{VehicleIndex: vehicleIndex, Index: j},
					}
					testSolution.InsertCall(call, inds, comparisonPoint)

					if _, ok := slices.BinarySearchFunc(validIndices, comparisonPoint, func(a, t utils.InsertionPoint) int {
						if t.PickupIndex.Index == a.PickupIndex.Index {
							return a.DeliveryIndex.Index - t.DeliveryIndex.Index
						}
						return a.PickupIndex.Index - t.PickupIndex.Index
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
							t.Log(utils.GetTour(testSolution.Solution, vehicleIndex), " - tour")

							t.Log("Capacities: ")
							t.Log(testSolution.VehicleCumulativeCapacities(vehicleIndex), " - capacity")
							SizeString := ""
							for _, call := range utils.GetTour(testSolution.Solution, vehicleIndex) {
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
							t.Log(utils.GetTour(testSolution.Solution, vehicleIndex), " - tour")

							t.Log("Capacities: ")
							t.Log(testSolution.VehicleCumulativeCapacities(vehicleIndex), " - capacity")
							SizeString := ""
							for _, call := range utils.GetTour(testSolution.Solution, vehicleIndex) {
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
							for _, call := range utils.GetCallNodeTour(problem, solution.Solution, vehicleIndex) {
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
	problem, _ := problem.LoadProblem("./Data/Call_18_Vehicle_5.txt")
	solution := GenerateInitialSolution(problem)
	twelveI := utils.InsertionPoint{
		PickupIndex: utils.RelativeIndex{
			VehicleIndex: 2,
			Index:        0,
		},
		DeliveryIndex: utils.RelativeIndex{
			VehicleIndex: 2,
			Index:        0,
		},
	}
	sixteenI := utils.InsertionPoint{
		PickupIndex: utils.RelativeIndex{
			VehicleIndex: 2,
			Index:        2,
		},
		DeliveryIndex: utils.RelativeIndex{
			VehicleIndex: 2,
			Index:        2,
		},
	}
	// tour: 12 12 16 16
	inds := utils.FindIndices(solution.Solution, 12)
	solution.InsertCall(12, inds, twelveI)
	inds = utils.FindIndices(solution.Solution, 16)
	solution.InsertCall(16, inds, sixteenI)
	inds = utils.FindIndices(solution.Solution, 10)
	tour := utils.GetCallNodeTour(problem, solution.Solution, 2)
	intTour := utils.GetTour(solution.Solution, 2)

	validIndices := solution.GetVehicleInsertionPoints(2, 10)
	t.Log(validIndices)
	for i := 0; i < len(tour)+1; i++ {
		for j := i; j < len(tour)+1; j++ {
			ts := solution.Copy()
			insertion := utils.InsertionPoint{
				PickupIndex:   utils.RelativeIndex{VehicleIndex: 2, Index: i},
				DeliveryIndex: utils.RelativeIndex{VehicleIndex: 2, Index: j},
			}
			ts.InsertCall(10, inds, insertion)
			shouldBeFeasible := false
			for _, ins := range validIndices {
				if ins.PickupIndex.Index == insertion.PickupIndex.Index &&
					ins.DeliveryIndex.Index == insertion.DeliveryIndex.Index {
					shouldBeFeasible = true
				}
			}

			if shouldBeFeasible && !ts.Feasible() {
				t.Errorf("Infeasible solution marked as feasible\nReason: %s", ts.infeasibleReason)
				t.Log("Indices: i:", i, "j:", j)
				t.Log("Tour: ", intTour)
				t.Log("Times: ", solution.VehicleCumulativeTimes(2))
				t.Log("Timeslack: ", solution.CalulateTimeSlack(tour, 2))
				t.Log("NewTour: ", utils.GetTour(ts.Solution, 2))
				t.Log("Newtimes: ", ts.VehicleCumulativeTimes(2))
			}
			if !shouldBeFeasible && ts.Feasible() {
				t.Log()
				t.Error("Feasible solution marked as infeasible")
				t.Log("Indices: i:", i, "j:", j)
				t.Log("Tour: ", intTour)
				t.Log("Times: ", solution.VehicleCumulativeTimes(2))
				t.Log("Timeslack: ", solution.CalulateTimeSlack(tour, 2))
				t.Log("NewTour: ", utils.GetTour(ts.Solution, 2))
				t.Log("Newtimes: ", ts.VehicleCumulativeTimes(2))
			}
		}
	}
}

func TestGetCostImprovement(t *testing.T) {
	problem, _ := problem.LoadProblem("./Data/Call_7_Vehicle_3.txt")
	solution := GenerateInitialSolution(problem)
	for i := 0; i < 100; i++ {
		call := rand.Intn(problem.NumberOfCalls) + 1
		indices := utils.FindIndices(solution.Solution, call)
		indices = solution.MoveCallToOutsource(call, indices)

		for vehicleIndex := 1; vehicleIndex <= problem.NumberOfVehicles; vehicleIndex++ {
			validIndices := solution.GetVehicleInsertionPoints(vehicleIndex, call)
			for _, insertionPoint := range validIndices {
				testSolution := solution.Copy()
				testSolution.InsertCall(call, indices, insertionPoint)
				if testSolution.Cost()-solution.Cost() != insertionPoint.CostDiff {
					t.Log("Originalcost: ", solution.Cost())
					t.Log("Cost of new solution: ", testSolution.Cost())

					t.Log("Real diff: ", solution.Cost()-testSolution.Cost())
					t.Log("Stored Costdiff: ", insertionPoint.CostDiff)
					t.Error("Not equal")
				}
			}
		}
	}
}
