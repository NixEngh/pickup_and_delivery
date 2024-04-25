package run

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/policy"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
)

func RunAssignment4(problems []*problem.Problem) {
    algorithms := map[string]algo.Algorithm{
        "Equal Probability": EqualProbability(),
        "Moderate": Moderate(),
        "Adventurous": Adventurous(),
        "Intense": Intense(),
        "Extreme": Extreme(),
    }

    Run(algorithms, problems)
}

func EqualProbability() algo.Algorithm {

    policy := policy.NewChooseRandomOperator(
		[]policy.OperatorScore{
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

    policy := policy.NewChooseRandomOperator(
		[]policy.OperatorScore{
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
	policy := policy.NewChooseRandomOperator(
		[]policy.OperatorScore{
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
	policy := policy.NewChooseRandomOperator(
		[]policy.OperatorScore{
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
	policy := policy.NewChooseRandomOperator(
		[]policy.OperatorScore{
			{Operator: operator.PlaceOptimally{}, Probability: 2},
			{Operator: operator.PlaceOptimallyInRandomVehicle{}, Probability: 1},
			{Operator: operator.PlaceRandomly{}, Probability: 1},
			{Operator: operator.PlaceFiveCallsRandomly{}, Probability: 2},
		},
		"Extreme",
    )
	return algo.SimulatedAnnealing(policy)
}
