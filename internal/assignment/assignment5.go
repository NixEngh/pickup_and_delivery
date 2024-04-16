package assignment

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
)

func Adaptive() algo.Algorithm {
	operators := []operator.Operator{
		operator.PlaceOptimally{},
		operator.PlaceOptimallyInRandomVehicle{},
		operator.PlaceRandomly{},
		operator.PlaceFiveCallsRandomly{},
	}
	policy := operator.NewLecturePolicy(50, 0.1, operators)
	return algo.SimulatedAnnealing(policy)
}
