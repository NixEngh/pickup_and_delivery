package main

import (
	"math/rand"
)

func (p *Problem) GenerateInitialSolution() Solution {
	var solution Solution

	solutionList := make([]int, p.NumberOfVehicles)
	for i := 0; i < p.NumberOfVehicles; i++ {
		solutionList[i] = 0
	}
	for i := 1; i <= p.NumberOfCalls; i++ {
		solutionList = append(solutionList, i)
		solutionList = append(solutionList, i)
	}

	outsourceCost := 0

	for i := 1; i <= p.NumberOfCalls; i++ {
		outsourceCost += p.Calls[i].CostOfNotTransporting
	}

	solution = Solution{
		Problem:           p,
		Solution:          solutionList,
		VehicleCost:       make([]int, p.NumberOfVehicles),
		OutSourceCost:     outsourceCost,
		Feasible:          true,
		UncheckedVehicles: make([]int, 0),
	}

	return solution
}

func (p *Problem) GenerateRandomSolution() Solution {
	vehicles := make([][]int, p.NumberOfVehicles+1)
	for i := 0; i <= p.NumberOfVehicles; i++ {
		vehicles[i] = make([]int, 0)
	}

	order := rand.Perm(p.NumberOfCalls)
	for _, call := range order {
		call += 1
		vehicle := rand.Intn(p.NumberOfVehicles + 1)
		vehicles[vehicle] = append(vehicles[vehicle], call)
		vehicles[vehicle] = append(vehicles[vehicle], call)
	}

	for i := 1; i <= p.NumberOfVehicles; i++ {
		rand.Shuffle(len(vehicles[i]), func(x, y int) {
			vehicles[i][x], vehicles[i][y] = vehicles[i][y], vehicles[i][x]
		})
		vehicles[i] = append(vehicles[i], 0)
	}

	solution := make([]int, 0)
	for i := 1; i <= p.NumberOfVehicles; i++ {
		solution = append(solution, vehicles[i]...)
	}

	solution = append(solution, vehicles[0]...)

	s := Solution{
		Problem:           p,
		Solution:          solution,
		VehicleCost:       make([]int, p.NumberOfVehicles),
		OutSourceCost:     0,
		Feasible:          false,
		UncheckedVehicles: make([]int, 0),
	}

	for i := 1; i <= p.NumberOfVehicles; i++ {
		s.UncheckedVehicles = append(s.UncheckedVehicles, i)
	}

	s.UpdateFeasibility()

	if s.Feasible {
		s.UpdateCosts()
	}

	return s
}

func (s *Solution) moveFromOutsource(callInds, zeroInds []int) {
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

	s.moveCall(callInds[0], position1)
	s.moveCall(callInds[1], position2)
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
			s.moveCall(pair.callIndex, pair.callIndex+pair.delta)
			return true
		}
	}

	return false
}

func (s *Solution) OneReinsert() {
	call := rand.Intn(s.Problem.NumberOfCalls) + 1

	indices := FindIndices(s.Solution, call, 0)
	callInds := indices[call]
	zeroInds := indices[0]

	isOutsourced := callInds[0] > zeroInds[len(zeroInds)-1]

	if isOutsourced {
		s.moveFromOutsource(callInds, zeroInds)
		return
	}

	outsourceProb := 0.5

	if rand.Float64() > outsourceProb {
		if ok := s.moveCallInVehicle(callInds, zeroInds); ok {
			return
		}
	}

	s.moveCall(callInds[1], zeroInds[len(zeroInds)-1])
	s.moveCall(callInds[0], zeroInds[len(zeroInds)-1])
	return
}
