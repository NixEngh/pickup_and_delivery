package main

// Updates (if relevant) and returns the cost of the solution
func (s *Solution) Cost() int {
    if len(s.VehiclesToCheckCost) == 0 {
        return s.cost
    }

    s.UpdateCosts()
    return s.cost
}

// Updates the cost of the solution
func (s *Solution) UpdateCosts() {

	total := 0
    for vehicle, _ := range s.VehiclesToCheckCost {
        total += s.VehicleCostFunction(vehicle)
    }

	s.OutSourceCost = s.OutSourceCostFunction()
	total += s.OutSourceCost
	s.cost = total
}

func (s *Solution) CostFunction() int {
	problem := s.Problem
	total := 0
	for i := 1; i <= problem.NumberOfVehicles; i++ {
		total += s.VehicleCostFunction(i)
	}
	total += s.OutSourceCostFunction()
	return total
}

func (s *Solution) IndexedCostFunction() int {
	return 0
}

func (s *Solution) OutSourceCostFunction() int {
	total := 0
	solution := s.Solution
	problem := s.Problem

	found := make(map[int]struct{})
	for i := len(solution) - 1; solution[i] != 0; i-- {
		if _, ok := found[solution[i]]; ok {
			total += problem.Calls[solution[i]].CostOfNotTransporting
		}
		found[solution[i]] = struct{}{}
	}

	return total
}

func (s *Solution) VehicleCostFunction(vehicleIndex int) int {
	// If outsourced:

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

	for _, call := range tour {
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

		found[call] = struct{}{}
	}

	return total
}
