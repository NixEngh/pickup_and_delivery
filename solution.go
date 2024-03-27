package main

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

	found := make(map[int]struct{})
	tour := GetTour(s.Solution, vehicleIndex)

	vehicle := s.Problem.Vehicles[vehicleIndex]
	currentTime := vehicle.StartingTime
	currentLoad := 0
	prevNode := vehicle.HomeNode

	for _, call := range tour {
		_, isDelivery := found[call]
		currentCall := s.Problem.Calls[call]

		// Checks
		// Time
		var timeWindowToCheck TimeWindow
		if !isDelivery {
			timeWindowToCheck = currentCall.PickupTimeWindow
		} else {
			timeWindowToCheck = currentCall.DeliveryTimeWindow
		}

		var timeAfterPrev int
		if !isDelivery {
			timeAfterPrev = currentTime + vehicle.TravelTimes[prevNode][currentCall.OriginNode]
		} else {
			timeAfterPrev = currentTime + vehicle.TravelTimes[prevNode][currentCall.DestinationNode]
		}

		if timeAfterPrev > timeWindowToCheck.UpperBound {
			return false
		}

		// Capacity
		if !isDelivery && currentLoad+currentCall.Size > vehicle.Capacity {
			return false
		}

		// Prepare for next iteration

		var nodeTime int
		if !isDelivery {
			currentLoad += currentCall.Size
			nodeTime = currentCall.OriginTimeForVehicle[vehicleIndex]
			prevNode = currentCall.OriginNode
		} else {
			currentLoad -= currentCall.Size
			nodeTime = currentCall.DestinationTimeForVehicle[vehicleIndex]
			prevNode = currentCall.DestinationNode
		}

		currentTime = max(timeAfterPrev, timeWindowToCheck.LowerBound) + nodeTime

		found[call] = struct{}{}
	}

	return true
}
