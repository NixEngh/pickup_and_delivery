package main

func (p *Problem) CostFunction(solution []int) int {
	total := 0
	for i := 1; i <= p.NumberOfVehicles; i++ {
		total += p.VehicleCostFunction(solution, i)
	}
	found := make(map[int]struct{})
	for i := len(solution) - 1; solution[i] != 0; i-- {
		if _, ok := found[solution[i]]; ok {
			total += p.Calls[solution[i]].CostOfNotTransporting
		}
        found[solution[i]] = struct{}{}
	}
	return total
}

func (p *Problem) VehicleCostFunction(solution []int, vehicleIndex int) int {
	found := make(map[int]struct{})
	tour := getTour(solution, vehicleIndex)

	vehicle := p.Vehicles[vehicleIndex]
	total := 0
	previousNode := vehicle.HomeNode

	for _, call := range tour {
		_, isDelivery := found[call]
		currentCall := p.Calls[call]

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
