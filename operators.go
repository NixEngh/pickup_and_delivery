package main

import (
	"math/rand"
)

func (p *Problem) GenerateInitialSolution() []int {
	solution := make([]int, p.NumberOfVehicles)
	for i := 0; i < p.NumberOfVehicles; i++ {
		solution[i] = 0
	}
	for i := 1; i <= p.NumberOfCalls; i++ {
		solution = append(solution, i)
		solution = append(solution, i)
	}
	return solution
}

func (p *Problem) GenerateRandomSolution() []int {
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
	return solution
}

func moveFromOutsource(p *Problem, solution, callInds, zeroInds []int) {
	possible_vehicles := p.CallVehicleMap[solution[callInds[0]]]
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

	MoveElement(solution, callInds[0], position1)
	MoveElement(solution, callInds[1], position2)
}

func moveCallInVehicle(p *Problem, solution, callInds, zeroInds []int) {
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
		if pair.callIndex+pair.delta < 0 || pair.callIndex+pair.delta == len(solution) || solution[pair.callIndex+pair.delta] == 0 {
			continue
		}
		MoveElement(solution, pair.callIndex, pair.callIndex+pair.delta)
	}
}

func OneReinsert(p *Problem, solution []int) {
	call := rand.Intn(p.NumberOfCalls) + 1

	indices := FindIndices(solution, call, 0)
	callInds := indices[call]
	zeroInds := indices[0]

	isOutsourced := callInds[0] > zeroInds[len(zeroInds)-1]

	if isOutsourced {
		moveFromOutsource(p, solution, callInds, zeroInds)
		return
	}

	outsourceProb := 0.3

	if rand.Float64() > outsourceProb {
		moveCallInVehicle(p, solution, callInds, zeroInds)
		return
	}

	MoveElement(solution, callInds[0], zeroInds[len(zeroInds)-1])
	MoveElement(solution, callInds[1], zeroInds[len(zeroInds)-1])
}
