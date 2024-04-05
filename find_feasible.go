package main

import "fmt"

// Iterate backwards and calculate how much time an insert can take without violating time feasibility. The slack at index i is the maximum time that can be added before index i without violating the time window constraints
func (s *Solution) CalulateTimeSlack(tour []CallNode, vehicleIndex int, startIndex int) []int {
	s.UpdateFeasibility()
	if len(tour) == 0 {
		return []int{}
	}

	timeSlack := make([]int, len(tour))

	var slack int
	var currentTime int

	for i := len(tour) - 1; i >= startIndex; i-- {
		currentNode := tour[i]

		currentTime = s.VehicleCumulativeTimes[vehicleIndex][i]
		constraint := currentNode.TimeWindow.UpperBound - max(currentTime, currentNode.TimeWindow.LowerBound)
		waitTime := max(0, currentNode.TimeWindow.LowerBound-currentTime)

		if slack == 0 {
			slack = constraint + waitTime
		} else {
			slack = min(slack, constraint) + waitTime
		}

		timeSlack[i] = slack
	}
	return timeSlack
}

// Assumes to have checked every insertionpoint between pickup and delivery
func (s *Solution) checkCapacityConstraint(call Call, insertAt InsertionPoint) bool {
    s.UpdateFeasibility()
    vehicleIndex := insertAt.pickupIndex.VehicleIndex
    vehicle := s.Problem.Vehicles[vehicleIndex]

    cumulativeCaps := s.VehicleCumulativeCapacities[vehicleIndex]

    //  |1|2|1|2|
    // |0|1|2|3|4|
    // |a|1|2|1|2|
    // x|0|1|2|3|4

    prevCap := vehicle.Capacity
    if insertAt.pickupIndex.Index > 0 {
        prevCap = cumulativeCaps[insertAt.pickupIndex.Index-1]
    }
    if call.Size>prevCap {
        return false
    }

    prevCap = prevCap - call.Size

    if insertAt.deliveryIndex.Index > insertAt.pickupIndex.Index {
        prevCap = cumulativeCaps[insertAt.deliveryIndex.Index-1]
    }
    if call.Size>prevCap-call.Size {
        return false
    }

    return true
}


func (s *Solution) oldCheckCapacityConstraint(callNode CallNode, vehicleIndex int, inVehicleInsertInd int) bool {
	s.UpdateFeasibility()
	vehicleCumulativeCapacities := s.VehicleCumulativeCapacities[vehicleIndex]
	vehicle := s.Problem.Vehicles[vehicleIndex]

	call := s.Problem.Calls[callNode.callIndex]

	if inVehicleInsertInd == 0 {
		return call.Size <= vehicle.Capacity
	}

	return vehicleCumulativeCapacities[inVehicleInsertInd-1]+call.Size <= vehicle.Capacity
}

func (s *Solution) checkTimeConstraint(callNode CallNode, vehicleIndex int, inVehicleInsertIndex int, tour []CallNode, timeSlack []int) bool {
	// Check time window constraint
	vehicle := s.Problem.Vehicles[vehicleIndex]
	prevTime := vehicle.StartingTime
	prevNode := vehicle.HomeNode
	s.UpdateFeasibility()
	if inVehicleInsertIndex > 0 {
		prevNode = tour[inVehicleInsertIndex-1].Node
		prevTime = max(s.VehicleCumulativeTimes[vehicleIndex][inVehicleInsertIndex-1], tour[inVehicleInsertIndex-1].TimeWindow.LowerBound) + tour[inVehicleInsertIndex-1].OperationTime
	}
	travelTime := vehicle.TravelTimes[prevNode][callNode.Node]
	if prevTime+travelTime > callNode.TimeWindow.UpperBound {
		return false
	}

	if inVehicleInsertIndex == len(tour) {
		return true
	}
	// Check slack constraint
	nextNode := tour[inVehicleInsertIndex].Node
	originalDeltaTime := vehicle.TravelTimes[prevNode][nextNode]

	currentNode := callNode.Node
	newTime := prevTime
	newTime += vehicle.TravelTimes[prevNode][currentNode]
	newTime = max(newTime, callNode.TimeWindow.LowerBound)
	newTime += callNode.OperationTime
	newTime += vehicle.TravelTimes[currentNode][nextNode]

	newDeltaTime := newTime - prevTime

	return newDeltaTime-originalDeltaTime <= timeSlack[inVehicleInsertIndex]
}

// For a feasible insertionPoint, update its costImprovement
func (i *InsertionPoint) storeCostDiff(s *Solution, call Call, tour []CallNode) {
	vehicleIndex := i.pickupIndex.VehicleIndex
	vehicle := s.Problem.Vehicles[vehicleIndex]

	prevNode := vehicle.HomeNode
	pickupIndex := i.pickupIndex.Index
	if pickupIndex > 0 {
		prevNode = tour[pickupIndex-1].Node
	}
	nextNode := tour[pickupIndex].Node
	currentCost := s.VehicleCost[vehicleIndex]

	afterPickupInsertedCost := currentCost - vehicle.TravelCosts[prevNode][nextNode]
	afterPickupInsertedCost += vehicle.TravelCosts[prevNode][call.OriginNode]
	afterPickupInsertedCost += call.OriginCostForVehicle[vehicleIndex]
	afterPickupInsertedCost += vehicle.TravelCosts[call.OriginNode][nextNode]

	prevNode = call.OriginNode
	deliveryIndex := i.deliveryIndex.Index

	if deliveryIndex > pickupIndex {
		prevNode = tour[deliveryIndex-1].Node
	}

	var insertedCost int
	if deliveryIndex >= len(tour) {
		insertedCost = afterPickupInsertedCost
		insertedCost += vehicle.TravelCosts[prevNode][call.DestinationNode]
		insertedCost += call.DestinationCostForVehicle[vehicleIndex]
	} else {
		nextNode = tour[deliveryIndex].Node
		insertedCost = afterPickupInsertedCost - vehicle.TravelCosts[prevNode][nextNode]
		insertedCost += vehicle.TravelCosts[prevNode][call.DestinationNode]
		insertedCost += call.DestinationCostForVehicle[vehicleIndex]
		insertedCost += vehicle.TravelCosts[call.DestinationNode][nextNode]
	}

	i.costDiff = currentCost - insertedCost + call.CostOfNotTransporting
}

// Get indices after which a call can be inserted. The call must be placed in outsource. The insertionpoints work such that the pickup is moved to the index i(such that its index is i). Then the delivery is moved
func (s *Solution) GetVehicleInsertionPoints(vehicleIndex, callNumber int) []InsertionPoint {
	result := make([]InsertionPoint, 0)

	tour := GetCallNodeTour(s.Problem, s.Solution, vehicleIndex)
	timeSlack := s.CalulateTimeSlack(tour, vehicleIndex, 0)

	indices := FindIndices(s.Solution, callNumber, 0)
	callIndices := indices[callNumber]
	zeroIndices := indices[0]

	if callIndices[0] < zeroIndices[len(zeroIndices)-1] {
		fmt.Println("Solution:", s.Solution)
		fmt.Println("Callindices: ", callIndices)
		panic("The call must be outsourced")
	}
	call := s.Problem.Calls[callNumber]

	pickupNode := call.GetCallNode(false, vehicleIndex)
	deliveryNode := call.GetCallNode(true, vehicleIndex)

    for i := 0; i < len(tour)+1; i++ {
        for j := i; j < len(tour)+1; j++ {
            insertAt := InsertionPoint{
                pickupIndex: RelativeIndex{VehicleIndex: vehicleIndex, Index: i},
                deliveryIndex: RelativeIndex{VehicleIndex: vehicleIndex, Index: j},
            }
            if !s.checkCapacityConstraint(call, insertAt) {
                break
            }
        }
    }

	for i := 0; i < len(tour)+1; i++ {
		if !s.checkCapacityConstraint(pickupNode, vehicleIndex, i) {
			continue
		}
		if !s.checkTimeConstraint(pickupNode, vehicleIndex, i, tour, timeSlack) {
			continue
		}

		potentialSolution := s.copy()
		relativeIndex := RelativeIndex{
			VehicleIndex: vehicleIndex,
			Index:        i,
		}

		potentialSolution.MoveRelativeToVehicle(callIndices[0], relativeIndex)
		potentialTour := GetCallNodeTour(potentialSolution.Problem, potentialSolution.Solution, vehicleIndex)
		potentialTimeSlack := potentialSolution.CalulateTimeSlack(potentialTour, vehicleIndex, i)

		for j := i + 1; j < len(potentialTour)+1; j++ {
			if !potentialSolution.checkCapacityConstraint(deliveryNode, vehicleIndex, j) {
				break
			}
			if !potentialSolution.checkTimeConstraint(deliveryNode, vehicleIndex, j, potentialTour, potentialTimeSlack) {
				continue
			}

			insertionPoint := InsertionPoint{
				pickupIndex:   RelativeIndex{VehicleIndex: vehicleIndex, Index: i},
				deliveryIndex: RelativeIndex{VehicleIndex: vehicleIndex, Index: j},
			}

			insertionPoint.storeCostDiff(potentialSolution, call, potentialTour)

			result = append(result, insertionPoint)
		}
	}

	return result
}
