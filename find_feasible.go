package main

import (
	"fmt"
	"math"
)

// Iterate backwards and calculate how much time an insert can take without violating time feasibility. The slack at index i is the maximum time that can be added before index i without violating the time window constraints
func (s *Solution) CalulateTimeSlack(tour []CallNode, vehicleIndex int) []int {
	if len(tour) == 0 {
		return []int{}
	}

	timeSlack := make([]int, len(tour))

	var slack int = math.MaxInt
	var currentTime int

	for i := len(tour) - 1; i >= 0; i-- {
		currentNode := tour[i]
		currentTime = s.VehicleCumulativeTimes(vehicleIndex)[i]
		constraint := currentNode.TimeWindow.UpperBound - max(currentTime, currentNode.TimeWindow.LowerBound)
		waitTime := max(0, currentNode.TimeWindow.LowerBound-currentTime)

        slack = min(slack, constraint) + waitTime

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

func (s *Solution) checkPickupTimeConstraint(callNode CallNode, tour []CallNode, insertAt RelativeIndex, timeSlack []int) (passed bool, arrivalAtPickupNode int) {
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

	if insertAt.Index < len(tour) {
		nextNode := tour[insertAt.Index]
		originalTimeAtNextNode := s.VehicleCumulativeTimes(insertAt.VehicleIndex)[insertAt.Index]

		timeAtNextNode := max(timeAtInsertedNode, callNode.TimeWindow.LowerBound) + callNode.OperationTime + vehicle.TravelTimes[callNode.Node][nextNode.Node]

		delay := timeAtNextNode - originalTimeAtNextNode
		if delay > timeSlack[insertAt.Index] {
			return false, 0
		}
	}
	return true, timeAtInsertedNode
}

func (s *Solution) checkDeliveryTimeConstraint(call Call, tour []CallNode, insertAt InsertionPoint, arrivalAtPreviousNode int, timeSlack []int) (passed bool, arrivalAtNextNode int) {
	s.Feasible()

	vehicle := s.Problem.Vehicles[insertAt.pickupIndex.VehicleIndex]
	pickupNode := call.GetCallNode(false, vehicle.Index)
	deliveryNode := call.GetCallNode(true, vehicle.Index)

	var prevCallNode CallNode = pickupNode
	var prevNode int = call.OriginNode

	if insertAt.pickupIndex.Index != insertAt.deliveryIndex.Index {
		prevCallNode = tour[insertAt.deliveryIndex.Index-1]
		prevNode = prevCallNode.Node
	}

	doneAtPrevNode := max(arrivalAtPreviousNode, prevCallNode.TimeWindow.LowerBound) + prevCallNode.OperationTime

	timeAtDelivery := doneAtPrevNode + vehicle.TravelTimes[prevNode][deliveryNode.Node]

	if timeAtDelivery > deliveryNode.TimeWindow.UpperBound {
		return false, 0
	}

	if insertAt.deliveryIndex.Index < len(tour) {
		nextNode := tour[insertAt.deliveryIndex.Index]

		arrivalAtNextNodeWithDelivery := max(timeAtDelivery, deliveryNode.TimeWindow.LowerBound) +
			deliveryNode.OperationTime +
			vehicle.TravelTimes[deliveryNode.Node][nextNode.Node]

		timeDiff := arrivalAtNextNodeWithDelivery - s.VehicleCumulativeTimes(vehicle.Index)[insertAt.deliveryIndex.Index]
		arrivalAtNextNode = doneAtPrevNode + vehicle.TravelTimes[prevNode][nextNode.Node]

		if timeDiff > timeSlack[insertAt.deliveryIndex.Index] {
			return false, arrivalAtNextNode
		}
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

	// Add costs for operating calls
	diff := call.OriginCostForVehicle[vehicleIndex] + call.DestinationCostForVehicle[vehicleIndex]
	diff -= call.CostOfNotTransporting

	// Add new travel costs
	diff += vehicle.TravelCosts[prePickupNode][call.OriginNode]
	var preDelivery int
	if i.pickupIndex.Index == i.deliveryIndex.Index {
		diff += vehicle.TravelCosts[call.OriginNode][call.DestinationNode]
		preDelivery = prePickupNode
	} else {
		postPickup := tour[i.pickupIndex.Index].Node
		diff += vehicle.TravelCosts[call.OriginNode][postPickup]

		preDelivery = tour[i.deliveryIndex.Index-1].Node
		diff += vehicle.TravelCosts[preDelivery][call.DestinationNode]
	}

	if i.deliveryIndex.Index < len(tour) {
		postDelivery := tour[i.deliveryIndex.Index].Node
		diff += vehicle.TravelCosts[call.DestinationNode][postDelivery]
		diff -= vehicle.TravelCosts[preDelivery][postDelivery]
	}

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
                if arrivalAtNextNode == 0 {
                    break
                }
				continue
			}
			insertAt.storeCostDiff(s, call, tour)
			result = append(result, insertAt)
		}
	}

	return result
}
