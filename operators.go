package main

import (
	"math/rand"
	"slices"
)

type operator interface {
	apply(s *Solution)
}

type PlaceOptimallyInRandomVehicle struct{}

func (o PlaceOptimallyInRandomVehicle) apply(s *Solution) {
	vehicleIndex := rand.Intn(s.Problem.NumberOfVehicles) + 1
	callIndex := rand.Intn(s.Problem.NumberOfCalls) + 1

	indices := FindIndices(s.Solution, 0, callIndex)
	s.MoveCallToOutsource(callIndex, indices)

	validIndices := s.GetVehicleInsertionPoints(vehicleIndex, callIndex)
	slices.SortFunc(validIndices, func(a, b InsertionPoint) int {
		return b.costDiff - a.costDiff
	})

	if len(validIndices) == 0 {
		return
	}
	insertionIndex := validIndices[0]
	indices = FindIndices(s.Solution, 0, callIndex)
	s.InsertCall(callIndex, indices, insertionIndex)
}

type PlaceOptimally struct {
	n_v_to_check int
}

// PlaceOptimally picks a call and concurrently checks the best possible location to place it
func (o PlaceOptimally) apply(s *Solution) {
	var n_v_to_check int
	if o.n_v_to_check == 0 {
		n_v_to_check = s.Problem.NumberOfVehicles
	}
	n_v_to_check = min(n_v_to_check, s.Problem.NumberOfVehicles)
	vecs := rand.Perm(s.Problem.NumberOfVehicles)

	callIndex := rand.Intn(s.Problem.NumberOfCalls) + 1

	indices := FindIndices(s.Solution, 0, callIndex)
	indices = s.MoveCallToOutsource(callIndex, indices)

	bestInsertionPoint := InsertionPoint{}
	for _, vehicleIndex := range vecs[:n_v_to_check] {
		vehicleIndex = vehicleIndex + 1

		validIndices := s.GetVehicleInsertionPoints(vehicleIndex, callIndex)
		slices.SortFunc(validIndices, func(a, b InsertionPoint) int {
			return a.costDiff - b.costDiff
		})

		if len(validIndices) == 0 {
			continue
		}
		if validIndices[0].costDiff < bestInsertionPoint.costDiff {
			bestInsertionPoint = validIndices[0]
		}
	}
    if bestInsertionPoint.costDiff == 0 {
        return
    }

	s.InsertCall(callIndex, indices, bestInsertionPoint)
}

type OldOneReinsert struct{}

func (o OldOneReinsert) apply(s *Solution) {
	move_in_vehicle := rand.Float64() < 0.5
	call := rand.Intn(s.Problem.NumberOfCalls) + 1

	inds := FindIndices(s.Solution, call, 0)
	if move_in_vehicle {
		o.moveCallInVehicle(s, inds[call], inds[0])
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

func (o OldOneReinsert) moveCallInVehicle(s *Solution, callInds, zeroInds []int) bool {
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
