package operator

import (
	"math"
	"math/rand"
	"slices"

	"github.com/NixEngh/pickup_and_delivery/internal/solution"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
)
type PlaceOptimallyInRandomVehicle struct{}

func (o PlaceOptimallyInRandomVehicle) Apply(s *solution.Solution) int {
	callIndex := rand.Intn(s.Problem.NumberOfCalls) + 1
	possibleVehicles := s.Problem.CallVehicleMap[callIndex]

	vehicleIndex := possibleVehicles[rand.Intn(len(possibleVehicles))]

	indices := utils.FindIndices(s.Solution, 0, callIndex)
	s.MoveCallToOutsource(callIndex, indices)

	validIndices := s.GetVehicleInsertionPoints(vehicleIndex, callIndex)
	slices.SortFunc(validIndices, func(a, b utils.InsertionPoint) int {
		return b.CostDiff - a.CostDiff
	})

	if len(validIndices) == 0 {
		return math.MaxInt32
	}
	insertionIndex := validIndices[0]
	indices = utils.FindIndices(s.Solution, 0, callIndex)
	s.InsertCall(callIndex, indices, insertionIndex)

	return s.Cost()
}


type PlaceRandomly struct{}

func (o PlaceRandomly) Apply(s *solution.Solution) int {
	callNumber := rand.Intn(s.Problem.NumberOfCalls) + 1
	s.PlaceCallRandomly(callNumber)
	return s.Cost()
}


type OldOneReinsert struct{}

func (o OldOneReinsert) Apply(s *solution.Solution) int {
	move_in_vehicle := rand.Float64() < 0.5
	call := rand.Intn(s.Problem.NumberOfCalls) + 1

	inds := utils.FindIndices(s.Solution, call, 0)
	if move_in_vehicle {
		o.moveCallInVehicle(s, inds[call], inds[0])
		return s.Cost()
	}

	zeroinds := inds[0]
	callinds := inds[call]

	if callinds[1] < zeroinds[len(zeroinds)-1] {
		s.MoveInSolution(callinds[1], zeroinds[len(zeroinds)-1])
		s.MoveInSolution(callinds[0], zeroinds[len(zeroinds)-1])
		return s.Cost()
	}

	possibleVehicles := s.Problem.CallVehicleMap[call]

	vehicle := possibleVehicles[rand.Intn(len(possibleVehicles))]

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
		return s.Cost()
	}
	insert_index := rand.Intn(vehicle_range_end-(vehicle_range_start+1)) + vehicle_range_start + 1

	s.MoveInSolution(callinds[0], insert_index)
	s.MoveInSolution(callinds[1], insert_index)

	return s.Cost()
}

func (o OldOneReinsert) moveCallInVehicle(s *solution.Solution, callInds, zeroInds []int) bool {
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
