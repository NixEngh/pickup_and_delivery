package main

import (
	"math/rand"
)

func (p *Problem) GenerateInitialSolution() *Solution {
	var solution Solution

	solutionList := make([]int, p.NumberOfVehicles)
	for i := 0; i < p.NumberOfVehicles; i++ {
		solutionList[i] = 0
	}
	for i := 1; i <= p.NumberOfCalls; i++ {
		solutionList = append(solutionList, i)
		solutionList = append(solutionList, i)
	}

	solution = Solution{
		Problem:                    p,
		Solution:                   solutionList,
		VehicleCost:                make([]int, p.NumberOfVehicles+1),
		OutSourceCost:              0,
		VehiclesToCheckCost:        make(map[int]bool, 0),
		VehiclesToCheckFeasibility: make(map[int]bool, 0),
		cost:                       0,
		feasible:                   true,
	}

	solution.OutSourceCost = solution.OutSourceCostFunction()
	solution.cost = solution.OutSourceCost

	return &solution
}

func (p *Problem) GenerateRandomSolution() *Solution {
	vehicles := make([][]int, p.NumberOfVehicles+1)
	for i := 0; i <= p.NumberOfVehicles; i++ {
		vehicles[i] = make([]int, 0)
	}

	order := rand.Perm(p.NumberOfCalls)
	for _, call := range order {
		call += 1
		vehicle := rand.Intn(p.NumberOfVehicles + 1)
		vehicles[vehicle] = append(vehicles[vehicle], call)
		vehicles[vehicle] = append(vehicles[vehicle], call)
	}

	for i := 1; i <= p.NumberOfVehicles; i++ {
		rand.Shuffle(len(vehicles[i]), func(x, y int) {
			vehicles[i][x], vehicles[i][y] = vehicles[i][y], vehicles[i][x]
		})
		vehicles[i] = append(vehicles[i], 0)
	}

	solution := make([]int, 0)
	for i := 1; i <= p.NumberOfVehicles; i++ {
		solution = append(solution, vehicles[i]...)
	}

	solution = append(solution, vehicles[0]...)

	s := Solution{
		Problem:                    p,
		Solution:                   solution,
		VehicleCost:                make([]int, p.NumberOfVehicles+1),
		OutSourceCost:              0,
		VehiclesToCheckCost:        make(map[int]bool, 0),
		VehiclesToCheckFeasibility: make(map[int]bool, 0),
	}

	for i := 1; i <= p.NumberOfVehicles; i++ {
		if vehicles[i][0] != 0 {
			s.VehiclesToCheckCost[i] = true
			s.VehiclesToCheckFeasibility[i] = true
		}
	}

	s.UpdateFeasibility()

	if s.Feasible() {
		s.UpdateCosts()
	}

	return &s
}

// Move an element in the solution.
func (s *Solution) MoveInSolution(from int, to int) {
	zeroIndices := FindIndices(s.Solution, 0)[0]

	fromVehicle, toVehicle := 0, 0

	for i, zeroIndex := range zeroIndices {
		if fromVehicle == 0 && zeroIndex >= from {
			fromVehicle = i + 1
		}

		if toVehicle == 0 && zeroIndex >= to {
			toVehicle = i + 1
		}
	}

	if fromVehicle != 0 {
		s.VehiclesToCheckCost[fromVehicle] = true
		s.VehiclesToCheckFeasibility[fromVehicle] = true
	}
	if toVehicle != 0 {
		s.VehiclesToCheckCost[toVehicle] = true
		s.VehiclesToCheckFeasibility[toVehicle] = true
	}

	MoveElement(s.Solution, from, to)
}

// Creates a copy of the solution
func (s *Solution) copy() *Solution {
	newSolution := make([]int, len(s.Solution))
	copy(newSolution, s.Solution)

	newVehicleCost := make([]int, len(s.VehicleCost))
	copy(newVehicleCost, s.VehicleCost)

	costVehicles := make(map[int]bool, len(s.VehiclesToCheckCost))
	for vehicle := range s.VehiclesToCheckCost {
		costVehicles[vehicle] = true
	}

	feasVehicles := make(map[int]bool, len(s.VehiclesToCheckFeasibility))
	for vehicle := range s.VehiclesToCheckFeasibility {
		feasVehicles[vehicle] = true
	}

	return &Solution{
		Problem:                    s.Problem,
		Solution:                   newSolution,
		VehicleCost:                newVehicleCost,
		OutSourceCost:              s.OutSourceCost,
		VehiclesToCheckCost:        costVehicles,
		VehiclesToCheckFeasibility: feasVehicles,
		feasible:                   s.feasible,
		cost:                       s.cost,
	}
}

// Iterate backwards and calculate how much time an insert can take without violating time feasibility
func (s *Solution) CalulateTimeSlack(tour []int, vehicleIndex int, startIndex int) []int {
	problem := s.Problem

	if len(tour) == 0 {
		return []int{}
	}

	timeSlack := make([]int, len(tour))
	isPickup := make(map[int]bool)
	var slack int

	var callTimeWindow TimeWindow
	var currentTime int
	for i := len(tour) - 1; i >= startIndex; i-- {
		call := problem.Calls[tour[i]]

		if isPickup[tour[i]] {
			callTimeWindow = call.DeliveryTimeWindow
		} else {
			callTimeWindow = call.PickupTimeWindow
		}

		currentTime = s.VehicleCumulativeTimes[vehicleIndex][i]
		constraint := callTimeWindow.UpperBound - max(currentTime, callTimeWindow.LowerBound)
		waitTime := max(0, callTimeWindow.LowerBound-currentTime)

		if slack == 0 {
			slack = constraint + waitTime
		} else {
			slack = min(slack, constraint) + waitTime
		}

		isPickup[tour[i]] = true
		timeSlack[i] = slack
	}
	return timeSlack
}

// Get indices after which a call can be inserted
func (s *Solution) GetVehicleInsertionPoints(tour []int, vehicleIndex, callIndex int) []InsertionPoint {
	feasibleInsertionPoints := make([]int, len(tour))
	pickupSlack := s.CalulateTimeSlack(tour, vehicleIndex, 0)

	vehicle := s.Problem.Vehicles[vehicleIndex]
	call := s.Problem.Calls[callIndex]
	result := make([]InsertionPoint, 0)

	isDelivery := make(map[int]bool)

	var fromNode, toNode CallNode
	for i := 0; i < len(tour)-1; i++ {
		prevCall := s.Problem.Calls[tour[i]]
		fromNode = prevCall.GetCallNode(isDelivery[tour[i]], vehicleIndex)
		isDelivery[tour[i]] = true

		prevCapacity := s.VehicleCumulativeCapacities[vehicleIndex][i]
		if prevCapacity+prevCall.Size > vehicle.Capacity {
			continue
		}

		nextCall := s.Problem.Calls[tour[i+1]]
		toNode = nextCall.GetCallNode(isDelivery[tour[i+1]], vehicleIndex)

		potentialNode := call.GetCallNode(false, vehicleIndex)

		originalTimeBetweenNodes := vehicle.TravelTimes[fromNode.Node][toNode.Node]

		hypotheticalStartTime := max(s.VehicleCumulativeTimes[vehicleIndex][i], fromNode.timeWindow.LowerBound) + fromNode.OperationTime
        hypotheticalTime := hypotheticalStartTime
        hypotheticalTime += vehicle.TravelTimes[fromNode.Node][callIndex]
		hypotheticalTime = max(hypotheticalTime, call.PickupTimeWindow.LowerBound) + potentialNode.OperationTime
        hypotheticalTime += vehicle.TravelTimes[callIndex][toNode.Node]

        newTimeBetweenNodes := hypotheticalTime - hypotheticalStartTime

		if newTimeBetweenNodes-originalTimeBetweenNodes > pickupSlack[i+1] {
			continue
		}

	}

	return result
}
