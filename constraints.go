package main

import "fmt"

// Updates (if relevant) and returns the cost of the solution
func (s *Solution) Cost() int {
	if len(s.VehiclesToCheckCost) == 0 {
		return s.cost
	}

	s.UpdateCosts()
	return s.cost
}

// Updates the cost of the solution using VehiclesToCheckCost
func (s *Solution) UpdateCosts() {
	newCost := s.cost
	newCost -= s.OutSourceCost

	for vehicle := range s.VehiclesToCheckCost {
		newCost -= s.VehicleCost[vehicle]
		s.VehicleCost[vehicle] = s.VehicleCostFunction(vehicle)
		newCost += s.VehicleCost[vehicle]
	}

	s.OutSourceCost = s.OutSourceCostFunction()

	newCost += s.OutSourceCost
	s.cost = newCost

	s.VehiclesToCheckCost = make(map[int]bool, 0)
}

// Calculates the total cost of the solution
func (s *Solution) CostFunction() int {
	problem := s.Problem
	total := 0
	for i := 1; i <= problem.NumberOfVehicles; i++ {
		total += s.VehicleCostFunction(i)
	}
	total += s.OutSourceCostFunction()
	return total
}

// Total cost of outsourced vehicles
func (s *Solution) OutSourceCostFunction() int {
	total := 0

	found := make(map[int]struct{})
	for i := len(s.Solution) - 1; s.Solution[i] != 0; i-- {
		if _, ok := found[s.Solution[i]]; ok {
			total += s.Problem.Calls[s.Solution[i]].CostOfNotTransporting
		}
		found[s.Solution[i]] = struct{}{}
	}

	return total
}

// Calculate the cost of one vehicle
func (s *Solution) VehicleCostFunction(vehicleIndex int) int {
	problem := s.Problem
	solution := s.Solution

	if vehicleIndex == problem.NumberOfVehicles+1 {
		return s.OutSourceCostFunction()
	}

	found := make(map[int]struct{})
	tour := GetTour(solution, vehicleIndex)

	vehicle := problem.Vehicles[vehicleIndex]
	total := 0
	previousNode := vehicle.HomeNode

	s.VehicleCumulativeCosts[vehicleIndex] = make([]int, len(tour))

	for i, call := range tour {
		_, isDelivery := found[call]
		currentCall := problem.Calls[call]

		if !isDelivery {
			total += vehicle.TravelCosts[previousNode][currentCall.OriginNode]
			total += currentCall.OriginCostForVehicle[vehicleIndex]
			previousNode = currentCall.OriginNode
		} else {
			total += vehicle.TravelCosts[previousNode][currentCall.DestinationNode]
			total += currentCall.DestinationCostForVehicle[vehicleIndex]
			previousNode = currentCall.DestinationNode
		}
		s.VehicleCumulativeCosts[vehicleIndex][i] = total
		found[call] = struct{}{}
	}

	return total
}

// Returns true if the solution is feasible
func (s *Solution) Feasible() bool {
	if len(s.VehiclesToCheckFeasibility) == 0 {
		return s.feasible
	}

	s.UpdateFeasibility()
	return s.feasible
}

// Checks every unchecked vehicle, updates feasibility and resets the list of unchecked vehicles
func (s *Solution) UpdateFeasibility() {
	s.feasible = true
	for vehicle := range s.VehiclesToCheckFeasibility {
		if !s.IsVehicleFeasible(vehicle) {
			s.feasible = false
			break
		}
	}

	s.VehiclesToCheckFeasibility = make(map[int]bool, 0)
}

// Returns true if the vehicle is feasible
func (s *Solution) IsVehicleFeasible(vehicleIndex int) bool {
	tour := GetCallNodeTour(s.Problem, s.Solution, vehicleIndex)

	vehicle := s.Problem.Vehicles[vehicleIndex]
	currentTime := vehicle.StartingTime
	currentLoad := 0
	prevNode := vehicle.HomeNode

	s.VehicleCumulativeTimes[vehicleIndex] = make([]int, len(tour))
	s.VehicleCumulativeCapacities[vehicleIndex] = make([]int, len(tour))
	openCount := 0

	for i, callNode := range tour {
		// Checks
		// Time
		timeAtCallNode := currentTime + vehicle.TravelTimes[prevNode][callNode.Node]

		s.VehicleCumulativeTimes[vehicleIndex][i] = timeAtCallNode

		if timeAtCallNode > callNode.TimeWindow.UpperBound {
			s.infeasibleReason = fmt.Sprintf("The time %d at index %d was too high for call %d with upperbound %d\ncumulative times: %v", timeAtCallNode, i, callNode.callIndex, callNode.TimeWindow.UpperBound, s.VehicleCumulativeTimes[vehicleIndex])
			return false
		}

		size := s.Problem.Calls[callNode.callIndex].Size
		if callNode.IsDelivery {
			size = -size
		}

		currCap := vehicle.Capacity - (currentLoad + size)
		s.VehicleCumulativeCapacities[vehicleIndex][i] = currCap

		// Capacity
		if currentLoad+size > vehicle.Capacity {
			s.infeasibleReason = fmt.Sprintf("The capacity %d at index %d was too low for call %d with size %d", currCap, i, callNode.callIndex, s.Problem.Calls[callNode.callIndex].Size)
			return false
		}

		// Prepare for next iteration
		currentLoad += size
		prevNode = callNode.Node

		currentTime = max(timeAtCallNode, callNode.TimeWindow.LowerBound) + callNode.OperationTime
		if !callNode.IsDelivery {
			openCount += 1
		} else {
			openCount -= 1
		}
	}

	if openCount != 0 {
		s.infeasibleReason = fmt.Sprintf("All calls weren't closed, openCount: %d", openCount)
		return false
	}

	return true
}
