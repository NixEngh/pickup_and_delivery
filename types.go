package main

type Vehicle struct {
	Index       int
	HomeNode    int
	StartingTime int
	Capacity    int
    TravelTimes [][]int
    TravelCosts [][]int
}

type Call struct {
	Index                 int
	OriginNode            int
    OriginCostForVehicle  []int
    OriginTimeForVehicle  []int
	DestinationNode       int
    DestinationCostForVehicle  []int
    DestinationTimeForVehicle  []int
	Size                  int
	CostOfNotTransporting int
	PickupTimeWindow      TimeWindow
	DeliveryTimeWindow    TimeWindow
}

type TimeWindow struct {
	LowerBound int
	UpperBound int
}

type Problem struct {
	NumberOfNodes int
	NumberOfVehicles int
	Vehicles        []Vehicle
	NumberOfCalls   int
	Calls           []Call
	CallVehicleMap    map[int][]int
}
