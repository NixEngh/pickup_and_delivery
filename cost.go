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
