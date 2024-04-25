package operator

import (
	"math"
	"math/rand"
	"slices"
	"sync"

	"github.com/NixEngh/pickup_and_delivery/internal/solution"
	"github.com/NixEngh/pickup_and_delivery/internal/utils"
)
type PlaceOptimally struct{}

func (o PlaceOptimally) ApplyWithoutConc(s *solution.Solution) int {
	callIndex := rand.Intn(s.Problem.NumberOfCalls) + 1
	possibleVehicles := s.Problem.CallVehicleMap[callIndex]

	indices := utils.FindIndices(s.Solution, 0, callIndex)
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
func (o PlaceOptimally) Apply(s *solution.Solution) int {
	callIndex := rand.Intn(s.Problem.NumberOfCalls) + 1

	possibleVehicles := s.Problem.CallVehicleMap[callIndex]

	indices := utils.FindIndices(s.Solution, 0, callIndex)
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
