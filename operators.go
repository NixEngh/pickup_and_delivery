package main

import (
	"math/rand"
)

type operator interface {
    apply(s *Solution)
}

type OneReinsert struct {
    problem *Problem
}

func (o *OneReinsert) apply(s *Solution) {
	call := rand.Intn(s.Problem.NumberOfCalls) + 1

	indices := FindIndices(s.Solution, call, 0)
	callInds := indices[call]
	zeroInds := indices[0]

	isOutsourced := callInds[0] > zeroInds[len(zeroInds)-1]

	if isOutsourced {
		o.moveFromOutsource(s, callInds, zeroInds)
		return
	}

	outsourceProb := 0.5

	if rand.Float64() > outsourceProb {
		if ok := s.moveCallInVehicle(callInds, zeroInds); ok {
			return
		}
	}

	s.MoveInSolution(callInds[1], zeroInds[len(zeroInds)-1])
	s.MoveInSolution(callInds[0], zeroInds[len(zeroInds)-1])
	return
}

func (o *OneReinsert) moveFromOutsource(s *Solution, callInds, zeroInds []int) {
	possible_vehicles := s.Problem.CallVehicleMap[s.Solution[callInds[0]]]
	vehicle := possible_vehicles[rand.Intn(len(possible_vehicles))]

	vehicleRangeEnd := zeroInds[vehicle-1]
	vehicleRangeStart := 0
	if vehicle > 1 {
		vehicleRangeStart = zeroInds[vehicle-2] + 1
	}

	position1, position2 := vehicleRangeEnd, vehicleRangeEnd

	if vehicleRangeStart != vehicleRangeEnd {
		position1 = rand.Intn(vehicleRangeEnd-vehicleRangeStart) + vehicleRangeStart
		position2 = rand.Intn(vehicleRangeEnd+1-vehicleRangeStart) + vehicleRangeStart
	}

	s.MoveInSolution(callInds[0], position1)
	s.MoveInSolution(callInds[1], position2)
	return
}

// Mutates *solution* such
func (s *Solution) moveCallInVehicle(callInds, zeroInds []int) bool {
	solution := s.Solution

	pairs := []struct{ callIndex, delta int }{
		{callInds[0], -1},
		{callInds[0], 1},
		{callInds[1], -1},
		{callInds[1], 1},
	}

	rand.Shuffle(4, func(i, j int) {
		pairs[i], pairs[j] = pairs[j], pairs[i]
	})

	for _, pair := range pairs {
		switch {
		case pair.callIndex+pair.delta < 0:
			continue
		case pair.callIndex+pair.delta == len(solution):
			continue
		case solution[pair.callIndex+pair.delta] == 0:
			continue
		case solution[pair.callIndex] == solution[pair.callIndex+pair.delta]:
			continue
		default:
			s.MoveInSolution(pair.callIndex, pair.callIndex+pair.delta)
			return true
		}
	}

	return false
}

type OldOneReinsert struct {}

func (OldOneReinsert) apply(s *Solution) {
	move_in_vehicle := rand.Float64() < 0.5
	call := rand.Intn(s.Problem.NumberOfCalls) + 1

	inds := FindIndices(s.Solution, call, 0)
	if move_in_vehicle {
		s.moveCallInVehicle(inds[call], inds[0])
		return
	}

	zeroinds := inds[0]
	callinds := inds[call]

	if callinds[1] < zeroinds[len(zeroinds)-1] {
		s.MoveInSolution(callinds[1], zeroinds[len(zeroinds)-1])
		s.MoveInSolution(callinds[0], zeroinds[len(zeroinds)-1])
		return
	}

	vehicle := rand.Intn(s.Problem.NumberOfVehicles) + 1

	vehicle_range_start := 0
	vehicle_range_end := 0

	if vehicle == 1 {
		vehicle_range_end = zeroinds[0]
	} else {
		vehicle_range_start = zeroinds[vehicle-2]
		vehicle_range_end = zeroinds[vehicle-1]
	}
	if vehicle_range_end-vehicle_range_start < 2 {
		s.MoveInSolution(callinds[0], vehicle_range_end)
		s.MoveInSolution(callinds[1], vehicle_range_end)
		return
	}
	insert_index := rand.Intn(vehicle_range_end-(vehicle_range_start+1)) + vehicle_range_start + 1

	s.MoveInSolution(callinds[0], insert_index)
	s.MoveInSolution(callinds[1], insert_index)

	return
}


