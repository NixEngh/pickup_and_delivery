package assignment

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
)

func EqualProbability() algo.Algorithm {

    policy := operator.NewChooseRandomOperator(
		[]operator.OperatorScore{
			{Operator: operator.PlaceOptimally{}, Probability: 1},
			{Operator: operator.PlaceOptimallyInRandomVehicle{}, Probability: 1},
			{Operator: operator.PlaceRandomly{}, Probability: 1},
			{Operator: operator.PlaceFiveCallsRandomly{}, Probability: 1},
		},
		"Equal Probability",
	)

	return algo.SimulatedAnnealing(policy)
}

func Moderate() algo.Algorithm {

    policy := operator.NewChooseRandomOperator(
		[]operator.OperatorScore{
			{Operator: operator.PlaceOptimally{}, Probability: 1},
			{Operator: operator.PlaceOptimallyInRandomVehicle{}, Probability: 2},
			{Operator: operator.PlaceRandomly{}, Probability: 2},
			{Operator: operator.PlaceFiveCallsRandomly{}, Probability: 1},
		},
		"Moderate",
	)
	return algo.SimulatedAnnealing(policy)
}

func Adventurous() algo.Algorithm {
	policy := operator.NewChooseRandomOperator(
		[]operator.OperatorScore{
			{Operator: operator.PlaceOptimally{}, Probability: 1},
			{Operator: operator.PlaceOptimallyInRandomVehicle{}, Probability: 1},
			{Operator: operator.PlaceRandomly{}, Probability: 2},
			{Operator: operator.PlaceFiveCallsRandomly{}, Probability: 2},
		},
		"Adventurous",
    )
	return algo.SimulatedAnnealing(policy)
}

func Intense() algo.Algorithm {
	policy := operator.NewChooseRandomOperator(
		[]operator.OperatorScore{
			{Operator: operator.PlaceOptimally{}, Probability: 2},
			{Operator: operator.PlaceOptimallyInRandomVehicle{}, Probability: 2},
			{Operator: operator.PlaceRandomly{}, Probability: 1},
			{Operator: operator.PlaceFiveCallsRandomly{}, Probability: 1},
		},
		"Intense",
    )
	return algo.SimulatedAnnealing(policy)
}

func Extreme() algo.Algorithm {
	policy := operator.NewChooseRandomOperator(
		[]operator.OperatorScore{
			{Operator: operator.PlaceOptimally{}, Probability: 2},
			{Operator: operator.PlaceOptimallyInRandomVehicle{}, Probability: 1},
			{Operator: operator.PlaceRandomly{}, Probability: 1},
			{Operator: operator.PlaceFiveCallsRandomly{}, Probability: 2},
		},
		"Extreme",
    )
	return algo.SimulatedAnnealing(policy)
}
