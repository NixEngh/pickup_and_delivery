package main

import (
	"fmt"
	"math/rand"
)

func (p *Problem) GenerateInitialSolution() *Solution {
	var solution Solution

	solutionList := make([]int, p.NumberOfVehicles)
	for i := 0; i < p.NumberOfVehicles; i++ {
		solutionList[i] = 0
	}
	for i := 1; i <= p.NumberOfCalls; i++ {
		solutionList = append(solutionList, i)
		solutionList = append(solutionList, i)
	}

	solution = Solution{
		Problem:                     p,
		Solution:                    solutionList,
		VehicleCost:                 make([]int, p.NumberOfVehicles+1),
		OutSourceCost:               0,
		VehiclesToCheckCost:         make(map[int]bool, 0),
		VehiclesToCheckFeasibility:  make(map[int]bool, 0),
		VehicleCumulativeCosts:      make([][]int, p.NumberOfVehicles+1),
		VehicleCumulativeCapacities: make([][]int, p.NumberOfVehicles+1),
		VehicleCumulativeTimes:      make([][]int, p.NumberOfVehicles+1),
		cost:                        0,
		feasible:                    true,
	}

	solution.OutSourceCost = solution.OutSourceCostFunction()
	solution.cost = solution.OutSourceCost

	return &solution
}

func (p *Problem) GenerateRandomSolution() *Solution {
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
		Problem:                     p,
		Solution:                    solution,
		VehicleCost:                 make([]int, p.NumberOfVehicles+1),
		OutSourceCost:               0,
		VehiclesToCheckCost:         make(map[int]bool, 0),
		VehiclesToCheckFeasibility:  make(map[int]bool, 0),
		VehicleCumulativeCosts:      make([][]int, p.NumberOfVehicles+1),
		VehicleCumulativeCapacities: make([][]int, p.NumberOfVehicles+1),
		VehicleCumulativeTimes:      make([][]int, p.NumberOfVehicles+1),
	}

	for i := 1; i <= p.NumberOfVehicles; i++ {
		if vehicles[i][0] != 0 {
			s.VehiclesToCheckCost[i] = true
			s.VehiclesToCheckFeasibility[i] = true
		}
	}

	s.UpdateFeasibility()

	if s.Feasible() {
		s.UpdateCosts()
	}

	return &s
}

// Move an element in the solution
func (s *Solution) MoveInSolution(from int, to int) {
	zeroIndices := FindIndices(s.Solution, 0)[0]

	fromVehicle, toVehicle := 0, 0

	for i, zeroIndex := range zeroIndices {
		if fromVehicle == 0 && from < zeroIndex {
			fromVehicle = i + 1
		}

		if toVehicle == 0 && (to < zeroIndex || (from > zeroIndex && to == zeroIndex)) {
			toVehicle = i + 1
		}
	}

	if fromVehicle != 0 {
		s.VehiclesToCheckCost[fromVehicle] = true
		s.VehiclesToCheckFeasibility[fromVehicle] = true
	}
	if toVehicle != 0 {
		s.VehiclesToCheckCost[toVehicle] = true
		s.VehiclesToCheckFeasibility[toVehicle] = true
	}

	MoveElement(s.Solution, from, to)
}

// Move to a position in a vehicle from [0, len(tour)+1>
func (s *Solution) MoveRelativeToVehicle(from int, newIndex RelativeIndex) {
	tourIndices := GetTourIndices(s.Solution, newIndex.VehicleIndex)
    tour := GetTour(s.Solution, newIndex.VehicleIndex)

	if tourIndices[0] == tourIndices[1] {
		if newIndex.Index != 0 {
            fmt.Println(tour)
			panic("Index out of bounds")
		}
		s.MoveInSolution(from, tourIndices[0])
		return
	}

	zeroIndices := FindIndices(s.Solution, 0)[0]
	absoluteIndex := newIndex.toAbsolute(zeroIndices)

	if from < tourIndices[0] {
		s.MoveInSolution(from, absoluteIndex-1)
		return
	}

	s.MoveInSolution(from, absoluteIndex)
}

func (s *Solution) InsertCall(callNumber int, inds map[int][]int, insertionPoint InsertionPoint) {
	callInds := inds[callNumber]
	deliveryInd := callInds[1]
	s.MoveRelativeToVehicle(callInds[0], insertionPoint.pickupIndex)

	if insertionPoint.pickupIndex.toAbsolute(inds[0]) > callInds[1] {
		deliveryInd -= 1
	}
	insertionPoint.deliveryIndex.Index += 1
	s.MoveRelativeToVehicle(deliveryInd, insertionPoint.deliveryIndex)
    
}

func (s *Solution) MoveCallToVehicle(callNumber int, inds map[int][]int, insertionPoint InsertionPoint) {
	newInds := s.MoveCallToOutsource(callNumber, inds)[callNumber]

	s.MoveRelativeToVehicle(newInds[0], insertionPoint.pickupIndex)
	s.MoveRelativeToVehicle(newInds[1], insertionPoint.deliveryIndex)
}

func (s *Solution) MoveCallToOutsource(callNumber int, inds map[int][]int) (newInds map[int][]int) {
	callInds := inds[callNumber]
	zeroInds := inds[0]

	if callInds[0] > zeroInds[len(zeroInds)-1] {
		return inds
	}

	moveTo := zeroInds[len(zeroInds)-1]
	s.MoveInSolution(callInds[1], moveTo)
	s.MoveInSolution(callInds[0], moveTo)

	newInds = FindIndices(s.Solution, 0, callNumber)
	return newInds
}

// Creates a copy of the solution
func (s *Solution) copy() *Solution {
	newSolution := make([]int, len(s.Solution))
	copy(newSolution, s.Solution)

	newVehicleCost := make([]int, len(s.VehicleCost))
	copy(newVehicleCost, s.VehicleCost)

	costVehicles := make(map[int]bool, len(s.VehiclesToCheckCost))
	for vehicle := range s.VehiclesToCheckCost {
		costVehicles[vehicle] = true
	}

	feasVehicles := make(map[int]bool, len(s.VehiclesToCheckFeasibility))
	for vehicle := range s.VehiclesToCheckFeasibility {
		feasVehicles[vehicle] = true
	}

	copyVehicleCumulativeCosts := make([][]int, len(s.VehicleCumulativeCosts))
	for i := range s.VehicleCumulativeCosts {
		copyVehicleCumulativeCosts[i] = make([]int, len(s.VehicleCumulativeCosts[i]))
		copy(copyVehicleCumulativeCosts[i], s.VehicleCumulativeCosts[i])
	}

	copyVehicleCumulativeCapacities := make([][]int, len(s.VehicleCumulativeCapacities))
	for i := range s.VehicleCumulativeCapacities {
		copyVehicleCumulativeCapacities[i] = make([]int, len(s.VehicleCumulativeCapacities[i]))
		copy(copyVehicleCumulativeCapacities[i], s.VehicleCumulativeCapacities[i])
	}

	copyVehicleCumulativeTimes := make([][]int, len(s.VehicleCumulativeTimes))
	for i := range s.VehicleCumulativeTimes {
		copyVehicleCumulativeTimes[i] = make([]int, len(s.VehicleCumulativeTimes[i]))
		copy(copyVehicleCumulativeTimes[i], s.VehicleCumulativeTimes[i])
	}

	return &Solution{
		Problem:                     s.Problem,
		Solution:                    newSolution,
		VehicleCost:                 newVehicleCost,
		OutSourceCost:               s.OutSourceCost,
		VehiclesToCheckCost:         costVehicles,
		VehiclesToCheckFeasibility:  feasVehicles,
		VehicleCumulativeCosts:      copyVehicleCumulativeCosts,
		VehicleCumulativeCapacities: copyVehicleCumulativeCapacities,
		VehicleCumulativeTimes:      copyVehicleCumulativeTimes,
		feasible:                    s.feasible,
		cost:                        s.cost,
	}
}
