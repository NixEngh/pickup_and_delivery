package main

type Vehicle struct {
	Index        int
	HomeNode     int
	StartingTime int
	Capacity     int
	TravelTimes  [][]int
	TravelCosts  [][]int
}

type Call struct {
	Index                     int
	OriginNode                int
	OriginCostForVehicle      []int
	OriginTimeForVehicle      []int
	DestinationNode           int
	DestinationCostForVehicle []int
	DestinationTimeForVehicle []int
	Size                      int
	CostOfNotTransporting     int
	PickupTimeWindow          TimeWindow
	DeliveryTimeWindow        TimeWindow
}

type TimeWindow struct {
	LowerBound int
	UpperBound int
}

type CallNode struct {
	Node          int
	callIndex     int
	IsDelivery    bool
	TimeWindow    TimeWindow
	OperationTime int
	Cost          int
}

func (c *Call) GetCallNode(isDelivery bool, vehicleIndex int) CallNode {
	if isDelivery {
		return CallNode{
			Node:          c.DestinationNode,
			callIndex:     c.Index,
			IsDelivery:    true,
			TimeWindow:    c.DeliveryTimeWindow,
			OperationTime: c.DestinationTimeForVehicle[vehicleIndex],
			Cost:          c.DestinationCostForVehicle[vehicleIndex],
		}
	} else {
		return CallNode{
			Node:          c.OriginNode,
			callIndex:     c.Index,
			IsDelivery:    false,
			TimeWindow:    c.PickupTimeWindow,
			OperationTime: c.OriginTimeForVehicle[vehicleIndex],
			Cost:          c.OriginCostForVehicle[vehicleIndex],
		}
	}
}

type RelativeIndex struct {
	VehicleIndex int
	Index        int
}

func (r *RelativeIndex) toAbsolute(zeroIndices []int) int {
	var from int = -1
	if r.VehicleIndex > 1 {
		from = zeroIndices[r.VehicleIndex-2]
	}
	return from + r.Index + 1
}

// The indices should not take into account the extension when inserting
type InsertionPoint struct {
	pickupIndex   RelativeIndex
	deliveryIndex RelativeIndex
	costDiff      int
}

type Problem struct {
	Name             string
	NumberOfNodes    int
	NumberOfVehicles int
	Vehicles         []Vehicle
	NumberOfCalls    int
	Calls            []Call
	CallVehicleMap   map[int][]int
}

type Solution struct {
	Problem                *Problem
	Solution               []int
	VehicleCost            []int
	VehicleCumulativeCosts [][]int
	// Contains the leftover capacity after each callnode
	VehicleCumulativeCapacities [][]int
	// For an index, contains the arrival time at that node
	VehicleCumulativeTimes     [][]int
	OutSourceCost              int
	VehiclesToCheckCost        map[int]bool
	VehiclesToCheckFeasibility map[int]bool
	cost                       int
	feasible                   bool
	infeasibleReason           string
}

type algorithm func(problem *Problem) (BestSolution *Solution, BestCost int)
