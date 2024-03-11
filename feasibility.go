package main

import "fmt"

func (p *Problem) IsFeasible(solution []int) bool {

	for i := 1; i <= p.NumberOfVehicles; i++ {
		if !p.IsVehicleFeasible(solution, i) {
			return false
		}
	}
	return true
}

func (p *Problem) IsVehicleFeasible(solution []int, vehicleIndex int) bool {

	found := make(map[int]struct{})
	tour := getTour(solution, vehicleIndex)
    fmt.Println(tour)

	vehicle := p.Vehicles[vehicleIndex]
	currentTime := vehicle.StartingTime
	currentLoad := 0
	prevNode := vehicle.HomeNode

	for _, call := range tour {
		_, isDelivery := found[call]
		currentCall := p.Calls[call]

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
			fmt.Println("Time after prev: ", timeAfterPrev)
			return false
		}

		// Capacity
		if !isDelivery && currentLoad+currentCall.Size > vehicle.Capacity {
			fmt.Println("Current load: ", currentLoad)
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
