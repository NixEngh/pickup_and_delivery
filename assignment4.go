package main

func EqualProbability() algorithm {
	policy := &ChooseRandomOperator{
		operators: []OperatorScore{
			{operator: PlaceOptimally{}, probability: 1},
			{operator: PlaceOptimallyInRandomVehicle{}, probability: 1},
			{operator: PlaceRandomly{}, probability: 1},
			{operator: PlaceFiveCallsRandomly{}, probability: 1},
		},
		name: "Equal Probability",
	}
	return SimulatedAnnealing(policy)
}

func Moderate() algorithm {
	policy := &ChooseRandomOperator{
		operators: []OperatorScore{
			{operator: PlaceOptimally{}, probability: 1},
			{operator: PlaceOptimallyInRandomVehicle{}, probability: 2},
			{operator: PlaceRandomly{}, probability: 2},
			{operator: PlaceFiveCallsRandomly{}, probability: 1},
		},
		name: "Moderate",
	}
	return SimulatedAnnealing(policy)
}

func Adventurous() algorithm {
	policy := &ChooseRandomOperator{
		operators: []OperatorScore{
			{operator: PlaceOptimally{}, probability: 1},
			{operator: PlaceOptimallyInRandomVehicle{}, probability: 1},
			{operator: PlaceRandomly{}, probability: 2},
			{operator: PlaceFiveCallsRandomly{}, probability: 2},
		},
		name: "Adventurous",
	}
	return SimulatedAnnealing(policy)
}

func Intense() algorithm {
	policy := &ChooseRandomOperator{
		operators: []OperatorScore{
			{operator: PlaceOptimally{}, probability: 2},
			{operator: PlaceOptimallyInRandomVehicle{}, probability: 2},
			{operator: PlaceRandomly{}, probability: 1},
			{operator: PlaceFiveCallsRandomly{}, probability: 1},
		},
		name: "Intense",
	}
	return SimulatedAnnealing(policy)
}

func Extreme() algorithm {
	policy := &ChooseRandomOperator{
		operators: []OperatorScore{
			{operator: PlaceOptimally{}, probability: 2},
			{operator: PlaceOptimallyInRandomVehicle{}, probability: 1},
			{operator: PlaceRandomly{}, probability: 1},
			{operator: PlaceFiveCallsRandomly{}, probability: 2},
		},
		name: "Extreme",
	}
	return SimulatedAnnealing(policy)
}
