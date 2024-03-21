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
    OutSourceCost int
    VehiclesToCheckCost map[int]bool
    VehiclesToCheckFeasibility map[int]bool
    cost int
    feasible bool
}

type algorithm func (problem *Problem)  (BestSolution *Solution, BestCost int)

