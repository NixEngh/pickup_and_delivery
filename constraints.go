package main

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
	found := make(map[int]struct{})
	tour := GetTour(s.Solution, vehicleIndex)

	vehicle := s.Problem.Vehicles[vehicleIndex]
	currentTime := vehicle.StartingTime
	currentLoad := 0
	prevNode := vehicle.HomeNode

	s.VehicleCumulativeTimes[vehicleIndex] = make([]int, len(tour))
	s.VehicleCumulativeCapacities[vehicleIndex] = make([]int, len(tour))

	for i, call := range tour {
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
		s.VehicleCumulativeTimes[vehicleIndex][i] = timeAfterPrev

		// Capacity
		if !isDelivery && currentLoad+currentCall.Size > vehicle.Capacity {
			return false
		}

        currCap := vehicle.Capacity - currentLoad+currentCall.Size
		s.VehicleCumulativeCapacities[vehicleIndex][i] = currCap
        if currCap < 0 {
            panic("Negative capacity")
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
