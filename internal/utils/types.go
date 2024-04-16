package utils

import "github.com/NixEngh/pickup_and_delivery/internal/problem"
type CallNode struct {
	Node          int
	CallIndex     int
	IsDelivery    bool
	TimeWindow    problem.TimeWindow
	OperationTime int
	Cost          int
}

func GetCallNode(c *problem.Call, isDelivery bool, vehicleIndex int) CallNode {
	if isDelivery {
		return CallNode{
			Node:          c.DestinationNode,
			CallIndex:     c.Index,
			IsDelivery:    true,
			TimeWindow:    c.DeliveryTimeWindow,
			OperationTime: c.DestinationTimeForVehicle[vehicleIndex],
			Cost:          c.DestinationCostForVehicle[vehicleIndex],
		}
	} else {
		return CallNode{
			Node:          c.OriginNode,
			CallIndex:     c.Index,
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
func (r *RelativeIndex) ToAbsolute(zeroIndices []int) int {
	var from int = -1
	if r.VehicleIndex > 1 {
		from = zeroIndices[r.VehicleIndex-2]
	}
	return from + r.Index + 1
}

// The indices should not take into account the extension when inserting
type InsertionPoint struct {
	PickupIndex   RelativeIndex
	DeliveryIndex RelativeIndex
    // Negative values are good
	CostDiff      int
}

