package main

import (
	"math/rand"
	"slices"
	"sync"
)

type Operator interface {
	apply(s *Solution)
}

type PlaceOptimallyInRandomVehicle struct{}

func (o PlaceOptimallyInRandomVehicle) apply(s *Solution) {
	callIndex := rand.Intn(s.Problem.NumberOfCalls) + 1
	possibleVehicles := s.Problem.CallVehicleMap[callIndex]

	vehicleIndex := possibleVehicles[rand.Intn(len(possibleVehicles))]

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

type PlaceOptimally struct{}

// PlaceOptimally picks a call and concurrently checks the best possible location to place it
func (o PlaceOptimally) apply(s *Solution) {
	callIndex := rand.Intn(s.Problem.NumberOfCalls) + 1

	possibleVehicles := s.Problem.CallVehicleMap[callIndex]

	indices := FindIndices(s.Solution, 0, callIndex)
	indices = s.MoveCallToOutsource(callIndex, indices)

	var wg = sync.WaitGroup{}
	insertionPoints := make(chan InsertionPoint)

	for _, vehicleIndex := range possibleVehicles {
		wg.Add(1)
		go func(vehicleIndex int) {
			defer wg.Done()
			validIndices := s.GetVehicleInsertionPoints(vehicleIndex, callIndex)
			slices.SortFunc(validIndices, func(a, b InsertionPoint) int {
				return a.costDiff - b.costDiff
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

	bestInsertionPoint := InsertionPoint{}

	for point := range insertionPoints {
		if point.costDiff < bestInsertionPoint.costDiff {
			bestInsertionPoint = point
		}
	}

	if bestInsertionPoint.costDiff == 0 {
		return
	}

	s.InsertCall(callIndex, indices, bestInsertionPoint)
}

type PlaceRandomly struct{}

func (o PlaceRandomly) apply(s *Solution) {
	callNumber := rand.Intn(s.Problem.NumberOfCalls) + 1
	s.placeCallRandomly(callNumber)
}

type PlaceFiveCallsRandomly struct{}

func (o PlaceFiveCallsRandomly) apply(s *Solution) {
	callsToMove := rand.Perm(s.Problem.NumberOfCalls)
	count := 0
	for _, callToMove := range callsToMove {
		if ok := s.placeCallRandomly(callToMove + 1); ok {
			count += 1
		}

		if count == 5 {
			break
		}
	}
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

type OperatorPolicy interface {
	ChooseOperator() Operator
	Name() string
}

type OperatorProbability struct {
	operator    Operator
	probability float64
}

type ChooseRandomOperator struct {
	operators []OperatorProbability
	name      string
}

func (c *ChooseRandomOperator) ChooseOperator() Operator {
	var total float64
	for _, op := range c.operators {
		total += op.probability
	}

	r := rand.Float64() * total
	for _, op := range c.operators {
		if r -= op.probability; r < 0 {
			return op.operator
		}
	}

	panic("Should not reach here")
}

func (c *ChooseRandomOperator) Name() string {
    return c.name
}

