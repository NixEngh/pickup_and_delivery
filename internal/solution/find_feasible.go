package solution

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/NixEngh/pickup_and_delivery/internal/problem"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
)

// Iterate backwards and calculate how much time an insert can take without violating time feasibility. The slack at index i is the maximum time that can be added before index i without violating the time window constraints
func (s *Solution) CalulateTimeSlack(tour []utils.CallNode, vehicleIndex int) []int {
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
func (s *Solution) checkCapacityConstraint(call problem.Call, insertAt utils.InsertionPoint) bool {
	s.Feasible()
	vehicleIndex := insertAt.PickupIndex.VehicleIndex
	vehicle := s.Problem.Vehicles[vehicleIndex]

	cumulativeCaps := s.VehicleCumulativeCapacities(vehicleIndex)

	prevCap := vehicle.Capacity
	if insertAt.DeliveryIndex.Index > 0 {
		prevCap = cumulativeCaps[insertAt.DeliveryIndex.Index-1]
	}
	if call.Size > prevCap {
		return false
	}

	return true
}

func (s *Solution) checkPickupTimeConstraint(callNode utils.CallNode, tour []utils.CallNode, insertAt utils.RelativeIndex, timeSlack []int) (passed bool, arrivalAtPickupNode int) {
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

func (s *Solution) checkDeliveryTimeConstraint(call problem.Call, tour []utils.CallNode, insertAt utils.InsertionPoint, arrivalAtPreviousNode int, timeSlack []int) (passed bool, arrivalAtNextNode int) {
	s.Feasible()

	vehicle := s.Problem.Vehicles[insertAt.PickupIndex.VehicleIndex]
	pickupNode := utils.GetCallNode(&call, false, vehicle.Index)
	deliveryNode := utils.GetCallNode(&call, true, vehicle.Index)

	var prevCallNode utils.CallNode = pickupNode
	var prevNode int = call.OriginNode

	if insertAt.PickupIndex.Index != insertAt.DeliveryIndex.Index {
		prevCallNode = tour[insertAt.DeliveryIndex.Index-1]
		prevNode = prevCallNode.Node
	}

	doneAtPrevNode := max(arrivalAtPreviousNode, prevCallNode.TimeWindow.LowerBound) + prevCallNode.OperationTime

	timeAtDelivery := doneAtPrevNode + vehicle.TravelTimes[prevNode][deliveryNode.Node]

	if timeAtDelivery > deliveryNode.TimeWindow.UpperBound {
		return false, 0
	}

	if insertAt.DeliveryIndex.Index < len(tour) {
		nextNode := tour[insertAt.DeliveryIndex.Index]

		arrivalAtNextNodeWithDelivery := max(timeAtDelivery, deliveryNode.TimeWindow.LowerBound) +
			deliveryNode.OperationTime +
			vehicle.TravelTimes[deliveryNode.Node][nextNode.Node]

		timeDiff := arrivalAtNextNodeWithDelivery - s.VehicleCumulativeTimes(vehicle.Index)[insertAt.DeliveryIndex.Index]
		arrivalAtNextNode = doneAtPrevNode + vehicle.TravelTimes[prevNode][nextNode.Node]

		if timeDiff > timeSlack[insertAt.DeliveryIndex.Index] {
			return false, arrivalAtNextNode
		}
	}

	return true, arrivalAtNextNode
}

// For a feasible insertionPoint, update its costImprovement
func storeCostDiff(i *utils.InsertionPoint, s *Solution, call problem.Call, tour []utils.CallNode) {
	vehicleIndex := i.PickupIndex.VehicleIndex
	vehicle := s.Problem.Vehicles[vehicleIndex]

	prePickupNode := vehicle.HomeNode
	if i.PickupIndex.Index > 0 {
		prePickupNode = tour[i.PickupIndex.Index-1].Node
	}

	// Add costs for operating calls
	diff := call.OriginCostForVehicle[vehicleIndex] + call.DestinationCostForVehicle[vehicleIndex]
	diff -= call.CostOfNotTransporting

	// Add new travel costs
	diff += vehicle.TravelCosts[prePickupNode][call.OriginNode]
	var preDelivery int
	if i.PickupIndex.Index == i.DeliveryIndex.Index {
		diff += vehicle.TravelCosts[call.OriginNode][call.DestinationNode]
		preDelivery = prePickupNode
	} else {
		postPickup := tour[i.PickupIndex.Index].Node
		diff += vehicle.TravelCosts[call.OriginNode][postPickup]

		preDelivery = tour[i.DeliveryIndex.Index-1].Node
		diff += vehicle.TravelCosts[preDelivery][call.DestinationNode]
	}

	if i.DeliveryIndex.Index < len(tour) {
		postDelivery := tour[i.DeliveryIndex.Index].Node
		diff += vehicle.TravelCosts[call.DestinationNode][postDelivery]
		diff -= vehicle.TravelCosts[preDelivery][postDelivery]
	}

	i.CostDiff = diff
}

// Get indices at which a call can be inserted. The call must be placed in outsource.
func (s *Solution) GetVehicleInsertionPoints(vehicleIndex, callNumber int) []utils.InsertionPoint {
	result := make([]utils.InsertionPoint, 0)

	indices := utils.FindIndices(s.Solution, callNumber)
	callIndices := indices[callNumber]
	zeroIndices := indices[0]

	if callIndices[0] < zeroIndices[len(zeroIndices)-1] {
		fmt.Println("Solution:", s.Solution)
		fmt.Println("Callindices: ", callIndices)
		panic("The call must be outsourced")
	}
	call := s.Problem.Calls[callNumber]

	pickupNode := utils.GetCallNode(&call, false, vehicleIndex)

	tour := utils.GetCallNodeTour(s.Problem, s.Solution, vehicleIndex)
	timeSlack := s.CalulateTimeSlack(tour, vehicleIndex)

	for i := 0; i < len(tour)+1; i++ {
		pickupPassed, arrivalAtNextNode := s.checkPickupTimeConstraint(pickupNode, tour, utils.RelativeIndex{VehicleIndex: vehicleIndex, Index: i}, timeSlack)

		if !pickupPassed {
			continue
		}

		for j := i; j < len(tour)+1; j++ {
			insertAt := utils.InsertionPoint{
				PickupIndex:   utils.RelativeIndex{VehicleIndex: vehicleIndex, Index: i},
				DeliveryIndex: utils.RelativeIndex{VehicleIndex: vehicleIndex, Index: j},
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
			storeCostDiff(&insertAt, s, call, tour)
			result = append(result, insertAt)
		}
	}

	return result
}

func (s *Solution) GetAllFeasibleNonConcurrent(callNumber int) []utils.InsertionPoint {
	inds := utils.FindIndices(s.Solution, callNumber)

	inds = s.MoveCallToOutsource(callNumber, inds)
	possibleVehicles := s.Problem.CallVehicleMap[callNumber]

	feasibleInsertions := make([]utils.InsertionPoint, 0)

	for _, vehicleIndex := range possibleVehicles {
		currentInsertions := s.GetVehicleInsertionPoints(vehicleIndex, callNumber)
		feasibleInsertions = append(feasibleInsertions, currentInsertions...)
	}

	sort.Slice(feasibleInsertions, func(i, j int) bool {
		return feasibleInsertions[i].CostDiff < feasibleInsertions[j].CostDiff
	})

	return feasibleInsertions
}

func (s *Solution) GetAllFeasible(callNumber int) []utils.InsertionPoint {
	inds := utils.FindIndices(s.Solution, callNumber)

	inds = s.MoveCallToOutsource(callNumber, inds)
	possibleVehicles := s.Problem.CallVehicleMap[callNumber]

	feasibleInsertions := make([]utils.InsertionPoint, 0)
	feasibleInsertionsChan := make(chan []utils.InsertionPoint, len(possibleVehicles))
	wg := sync.WaitGroup{}

	s.Feasible()

	for _, vehicleIndex := range possibleVehicles {
		wg.Add(1)
		go func(vehicleIndex int) {
			defer wg.Done()
			currentInsertions := s.GetVehicleInsertionPoints(vehicleIndex, callNumber)
			feasibleInsertionsChan <- currentInsertions
		}(vehicleIndex)
	}

	go func() {
		wg.Wait()
		close(feasibleInsertionsChan)
	}()

	for insertions := range feasibleInsertionsChan {
		feasibleInsertions = append(feasibleInsertions, insertions...)
	}

	sort.Slice(feasibleInsertions, func(i, j int) bool {
		return feasibleInsertions[i].CostDiff < feasibleInsertions[j].CostDiff
	})

	return feasibleInsertions
}
