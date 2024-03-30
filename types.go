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
	IsDelivery    bool
	timeWindow    TimeWindow
	OperationTime int
}

func (c *Call) GetCallNode(isDelivery bool, vehicleIndex int) CallNode {
    if isDelivery {
        return CallNode{Node: c.DestinationNode, IsDelivery: true, timeWindow: c.DeliveryTimeWindow, OperationTime: c.DestinationTimeForVehicle[vehicleIndex]}
    } else {
        return CallNode{Node: c.OriginNode, IsDelivery: false, timeWindow: c.PickupTimeWindow, OperationTime: c.OriginTimeForVehicle[vehicleIndex]}
    }
}

type InsertionPoint struct {
	pickupIndex   int
	deliveryIndex int
	cost          int
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
    Problem *Problem
    Solution []int
    VehicleCost []int
    VehicleCumulativeCosts [][]int
    // todo mae sure these are long enough
    // Contains the leftover capacity after each call
    VehicleCumulativeCapacities [][]int
    // todo
    // For an index, contains the arrival time at that node
    VehicleCumulativeTimes [][]int
    OutSourceCost int
    VehiclesToCheckCost map[int]bool
    VehiclesToCheckFeasibility map[int]bool
    cost int
    feasible bool
}

type algorithm func (problem *Problem)  (BestSolution *Solution, BestCost int)

