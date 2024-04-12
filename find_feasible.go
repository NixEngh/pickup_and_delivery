package main

import (
	"fmt"
)

// Iterate backwards and calculate how much time an insert can take without violating time feasibility. The slack at index i is the maximum time that can be added before index i without violating the time window constraints
func (s *Solution) CalulateTimeSlack(tour []CallNode, vehicleIndex int) []int {
	s.Feasible()
	if len(tour) == 0 {
		return []int{}
	}

	timeSlack := make([]int, len(tour))

	var slack int
	var currentTime int
	var firstIteration bool = true

	for i := len(tour) - 1; i >= 0; i-- {
		currentNode := tour[i]

		currentTime = s.VehicleCumulativeTimes(vehicleIndex)[i]
		constraint := currentNode.TimeWindow.UpperBound - max(currentTime, currentNode.TimeWindow.LowerBound)
		waitTime := max(0, currentNode.TimeWindow.LowerBound-currentTime)

		if firstIteration {
			slack = constraint + waitTime
			firstIteration = false
		} else {
			slack = min(slack, constraint) + waitTime
		}

		timeSlack[i] = slack
	}
	return timeSlack
}

// Assumes to have checked every insertionpoint between pickup and delivery
func (s *Solution) checkCapacityConstraint(call Call, insertAt InsertionPoint) bool {
	s.Feasible()
	vehicleIndex := insertAt.pickupIndex.VehicleIndex
	vehicle := s.Problem.Vehicles[vehicleIndex]

	cumulativeCaps := s.VehicleCumulativeCapacities(vehicleIndex)

	prevCap := vehicle.Capacity
	if insertAt.deliveryIndex.Index > 0 {
		prevCap = cumulativeCaps[insertAt.deliveryIndex.Index-1]
	}
	if call.Size > prevCap {
		return false
	}

	return true
}

func (s *Solution) checkPickupTimeConstraint(callNode CallNode, tour []CallNode, insertAt RelativeIndex, timeSlack []int) (passed bool, arrivalAtNextNode int) {
	s.Feasible()
	vehicle := s.Problem.Vehicles[insertAt.VehicleIndex]
	prevTime := vehicle.StartingTime

	prevNode := vehicle.HomeNode

	if insertAt.Index > 0 {
		prevCallNode := tour[insertAt.Index-1]
		prevNode = prevCallNode.Node
		prevTime = max(prevCallNode.TimeWindow.LowerBound, s.VehicleCumulativeTimes(vehicle.Index)[insertAt.Index-1]) + prevCallNode.OperationTime
	}

	timeAtInsertedNode := prevTime + vehicle.TravelTimes[prevNode][callNode.Node]

	if timeAtInsertedNode > callNode.TimeWindow.UpperBound {
		return false, 0
	}

	var timeAtNextNode int

	if insertAt.Index < len(tour) {
		nextNode := tour[insertAt.Index]
		originalTimeAtNextNode := s.VehicleCumulativeTimes(insertAt.VehicleIndex)[insertAt.Index]

		timeAtNextNode = max(timeAtInsertedNode, callNode.TimeWindow.LowerBound) + callNode.OperationTime + vehicle.TravelTimes[callNode.Node][nextNode.Node]

		delay := timeAtNextNode - originalTimeAtNextNode
		if delay > timeSlack[insertAt.Index] {
			return false, 0
		}

	}
	return true, timeAtNextNode
}

func (s *Solution) handleConsecutiveInsertion(call Call, tour []CallNode, insertAt RelativeIndex, timeSlack []int) (passed bool) {
	vehicle := s.Problem.Vehicles[insertAt.VehicleIndex]

	pickupNode := call.GetCallNode(false, insertAt.VehicleIndex)
	deliveryNode := call.GetCallNode(true, insertAt.VehicleIndex)

	prevNode := vehicle.HomeNode
	// Time finished at previous node
	prevTime := vehicle.StartingTime

	if insertAt.Index > 0 {
		prevNodeCall := tour[insertAt.Index-1]
		prevTime = max(s.VehicleCumulativeTimes(insertAt.VehicleIndex)[insertAt.Index-1], prevNodeCall.TimeWindow.LowerBound)
		prevTime += prevNodeCall.OperationTime

		prevNode = prevNodeCall.Node
	}

	time := prevTime + vehicle.TravelTimes[prevNode][pickupNode.Node]

	if time > pickupNode.TimeWindow.UpperBound {
		return false
	}
	time = max(time, pickupNode.TimeWindow.LowerBound)
	time += pickupNode.OperationTime
	time += vehicle.TravelTimes[pickupNode.Node][deliveryNode.Node]

	if time > deliveryNode.TimeWindow.UpperBound {
		return false
	}
	time = max(time, deliveryNode.TimeWindow.LowerBound)
	time += deliveryNode.OperationTime

	if insertAt.Index < len(tour) {
		nextNode := tour[insertAt.Index]

		time += vehicle.TravelTimes[deliveryNode.Node][nextNode.Node]

		originalTime := s.VehicleCumulativeTimes(vehicle.Index)[insertAt.Index]

		timeDiff := time - originalTime

		if timeDiff > timeSlack[insertAt.Index] {

			return false
		} else {
			return true
		}
	}

	return true
}

func (s *Solution) checkDeliveryTimeConstraint(call Call, tour []CallNode, insertAt InsertionPoint, arrivalAtPreviousNode int, timeSlack []int) (passed bool, arrivalAtNextNode int) {
	s.Feasible()

	vehicle := s.Problem.Vehicles[insertAt.pickupIndex.VehicleIndex]
	prevNode := tour[insertAt.deliveryIndex.Index-1]
	doneAtPrevNode := max(arrivalAtPreviousNode, prevNode.TimeWindow.LowerBound) + prevNode.OperationTime

	deliveryNode := call.GetCallNode(true, vehicle.Index)

	timeAtDelivery := doneAtPrevNode + vehicle.TravelTimes[prevNode.Node][deliveryNode.Node]

	if timeAtDelivery > deliveryNode.TimeWindow.UpperBound {
		return false, 0
	}

	if insertAt.deliveryIndex.Index < len(tour) {
		nextNode := tour[insertAt.deliveryIndex.Index]

		arrivalAtNextNodeWithDelivery := max(timeAtDelivery, deliveryNode.TimeWindow.LowerBound) +
			deliveryNode.OperationTime +
			vehicle.TravelTimes[deliveryNode.Node][nextNode.Node]

		timeDiff := arrivalAtNextNodeWithDelivery - s.VehicleCumulativeTimes(vehicle.Index)[insertAt.deliveryIndex.Index]
		if timeDiff > timeSlack[insertAt.deliveryIndex.Index] {
			return false, 0
		}

		arrivalAtNextNode = doneAtPrevNode + vehicle.TravelTimes[prevNode.Node][nextNode.Node]
	}

	return true, arrivalAtNextNode
}

// For a feasible insertionPoint, update its costImprovement
func (i *InsertionPoint) storeCostDiff(s *Solution, call Call, tour []CallNode) {
	vehicleIndex := i.pickupIndex.VehicleIndex
	vehicle := s.Problem.Vehicles[vehicleIndex]

	prePickupNode := vehicle.HomeNode
	if i.pickupIndex.Index > 0 {
		prePickupNode = tour[i.pickupIndex.Index-1].Node
	}

	if i.pickupIndex.Index == i.deliveryIndex.Index {
		diff := 0
		diff += vehicle.TravelCosts[prePickupNode][call.OriginNode]
		diff += call.OriginCostForVehicle[vehicleIndex]
		diff += vehicle.TravelCosts[call.OriginNode][call.DestinationNode]
		diff += call.DestinationCostForVehicle[vehicleIndex]

		if i.deliveryIndex.Index < len(tour) {
			diff += vehicle.TravelCosts[call.DestinationNode][tour[i.deliveryIndex.Index].Node]
			diff -= vehicle.TravelCosts[prePickupNode][tour[i.deliveryIndex.Index].Node]
		}
		diff -= call.CostOfNotTransporting

		i.costDiff = diff
		return
	}

	diff := 0
	diff += vehicle.TravelCosts[prePickupNode][call.OriginNode]
	diff += call.OriginCostForVehicle[vehicleIndex]
	diff += vehicle.TravelCosts[call.OriginNode][tour[i.pickupIndex.Index].Node]
	diff -= vehicle.TravelCosts[prePickupNode][tour[i.pickupIndex.Index].Node]

	preDeliveryNode := tour[i.deliveryIndex.Index-1].Node
	diff += vehicle.TravelCosts[preDeliveryNode][call.DestinationNode]
	diff += call.DestinationCostForVehicle[vehicleIndex]

	if i.deliveryIndex.Index < len(tour) {
		diff += vehicle.TravelCosts[call.DestinationNode][tour[i.deliveryIndex.Index].Node]
		diff -= vehicle.TravelCosts[preDeliveryNode][tour[i.deliveryIndex.Index].Node]
	}

	diff -= call.CostOfNotTransporting

	i.costDiff = diff
}

// Get indices at which a call can be inserted. The call must be placed in outsource.
func (s *Solution) GetVehicleInsertionPoints(vehicleIndex, callNumber int) []InsertionPoint {
	result := make([]InsertionPoint, 0)

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

	tour := GetCallNodeTour(s.Problem, s.Solution, vehicleIndex)
	timeSlack := s.CalulateTimeSlack(tour, vehicleIndex)

	for i := 0; i < len(tour)+1; i++ {
		pickupPassed, arrivalAtNextNode := s.checkPickupTimeConstraint(pickupNode, tour, RelativeIndex{VehicleIndex: vehicleIndex, Index: i}, timeSlack)

		if !pickupPassed {
			continue
		}

		for j := i; j < len(tour)+1; j++ {
			insertAt := InsertionPoint{
				pickupIndex:   RelativeIndex{VehicleIndex: vehicleIndex, Index: i},
				deliveryIndex: RelativeIndex{VehicleIndex: vehicleIndex, Index: j},
			}

			if !s.checkCapacityConstraint(call, insertAt) {
				break
			}

			var passedTimeCheck bool
			passedTimeCheck, arrivalAtNextNode = s.checkDeliveryTimeConstraint(call, tour, insertAt, arrivalAtNextNode, timeSlack)

			if !passedTimeCheck {
				break
			}
			insertAt.storeCostDiff(s, call, tour)
			result = append(result, insertAt)
		}
	}

	return result
}
