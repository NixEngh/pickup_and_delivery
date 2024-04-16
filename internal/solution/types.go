package solution

import "github.com/NixEngh/pickup_and_delivery/internal/problem"
type Solution struct {
	Problem                *problem.Problem
	Solution               []int
	VehicleCost            []int
	vehicleCumulativeCosts [][]int
	// Contains the leftover capacity after each callnode
	vehicleCumulativeCapacities [][]int
	// For an index, contains the arrival time at that node
	vehicleCumulativeTimes     [][]int
	OutSourceCost              int
	VehiclesToCheckCost        map[int]bool
	VehiclesToCheckFeasibility map[int]bool
	cost                       int
	feasible                   bool
	infeasibleReason           string
}

