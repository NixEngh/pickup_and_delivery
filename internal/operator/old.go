package operator

import (
	"math"
	"math/rand"
	"slices"
	"sync"

	"github.com/NixEngh/pickup_and_delivery/internal/solution"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
)

type PlaceOptimallyInRandomVehicle struct{}

func (o *PlaceOptimallyInRandomVehicle) Apply(s *solution.Solution) int {
	callIndex := rand.Intn(s.Problem.NumberOfCalls) + 1
	possibleVehicles := s.Problem.CallVehicleMap[callIndex]

	vehicleIndex := possibleVehicles[rand.Intn(len(possibleVehicles))]

	indices := utils.FindIndices(s.Solution, callIndex)
	s.MoveCallToOutsource(callIndex, indices)

	validIndices := s.GetVehicleInsertionPoints(vehicleIndex, callIndex)
	slices.SortFunc(validIndices, func(a, b utils.InsertionPoint) int {
		return b.CostDiff - a.CostDiff
	})

	if len(validIndices) == 0 {
		return math.MaxInt32
	}
	insertionIndex := validIndices[0]
	indices = utils.FindIndices(s.Solution, callIndex)
	s.InsertCall(callIndex, indices, insertionIndex)

	return s.Cost()
}

type PlaceRandomly struct{}

func (o *PlaceRandomly) Apply(s *solution.Solution) int {
	callNumber := rand.Intn(s.Problem.NumberOfCalls) + 1
	s.PlaceCallRandomly(callNumber)
	return s.Cost()
}

type OldOneReinsert struct{}

func (o *OldOneReinsert) Apply(s *solution.Solution) int {
	moveInVehicle := rand.Float64() < 0.5
	call := rand.Intn(s.Problem.NumberOfCalls) + 1

	inds := utils.FindIndices(s.Solution, call)
	if moveInVehicle {
		o.moveCallInVehicle(s, inds[call])
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

	vehicleRangeStart := 0
	vehicleRangeEnd := 0

	if vehicle == 1 {
		vehicleRangeEnd = zeroinds[0]
	} else {
		vehicleRangeStart = zeroinds[vehicle-2]
		vehicleRangeEnd = zeroinds[vehicle-1]
	}
	if vehicleRangeEnd-vehicleRangeStart < 2 {
		s.MoveInSolution(callinds[0], vehicleRangeEnd)
		s.MoveInSolution(callinds[1], vehicleRangeEnd)
		return s.Cost()
	}
	insertIndex := rand.Intn(vehicleRangeEnd-(vehicleRangeStart+1)) + vehicleRangeStart + 1

	s.MoveInSolution(callinds[0], insertIndex)
	s.MoveInSolution(callinds[1], insertIndex)

	return s.Cost()
}

func (o *OldOneReinsert) moveCallInVehicle(s *solution.Solution, callInds []int) bool {
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

type PlaceOptimally struct{}

func (o *PlaceOptimally) ApplyWithoutConc(s *solution.Solution) int {
	callIndex := rand.Intn(s.Problem.NumberOfCalls) + 1
	possibleVehicles := s.Problem.CallVehicleMap[callIndex]

	indices := utils.FindIndices(s.Solution, callIndex)
	indices = s.MoveCallToOutsource(callIndex, indices)
	insertionPoints := make([]utils.InsertionPoint, 0)

	for _, vehicleIndex := range possibleVehicles {
		validIndices := s.GetVehicleInsertionPoints(vehicleIndex, callIndex)
		slices.SortFunc(validIndices, func(a, b utils.InsertionPoint) int {
			return a.CostDiff - b.CostDiff
		})

		if len(validIndices) == 0 {
			continue
		}
		insertionPoints = append(insertionPoints, validIndices[0])
	}

	bestInsertionPoint := utils.InsertionPoint{}

	for _, point := range insertionPoints {
		if point.CostDiff < bestInsertionPoint.CostDiff {
			bestInsertionPoint = point
		}
	}

	if bestInsertionPoint.CostDiff == 0 {
		return math.MaxInt32
	}

	s.InsertCall(callIndex, indices, bestInsertionPoint)

	return s.Cost()
}

// PlaceOptimally picks a call and concurrently checks the best possible location to place it
func (o *PlaceOptimally) Apply(s *solution.Solution) int {
	callIndex := rand.Intn(s.Problem.NumberOfCalls) + 1

	possibleVehicles := s.Problem.CallVehicleMap[callIndex]

	indices := utils.FindIndices(s.Solution, callIndex)
	indices = s.MoveCallToOutsource(callIndex, indices)

	var wg = sync.WaitGroup{}
	insertionPoints := make(chan utils.InsertionPoint)

	for _, vehicleIndex := range possibleVehicles {
		wg.Add(1)
		go func(vehicleIndex int) {
			defer wg.Done()
			validIndices := s.GetVehicleInsertionPoints(vehicleIndex, callIndex)
			slices.SortFunc(validIndices, func(a, b utils.InsertionPoint) int {
				return a.CostDiff - b.CostDiff
			})

			if len(validIndices) == 0 {
				return
			}
			insertionPoints <- validIndices[0]
		}(vehicleIndex)
	}

	go func() {
		wg.Wait()
		close(insertionPoints)
	}()

	bestInsertionPoint := utils.InsertionPoint{}

	for point := range insertionPoints {
		if point.CostDiff < bestInsertionPoint.CostDiff {
			bestInsertionPoint = point
		}
	}

	if bestInsertionPoint.CostDiff == 0 {
		return math.MaxInt32
	}

	s.InsertCall(callIndex, indices, bestInsertionPoint)

	return s.Cost()
}
