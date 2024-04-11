package main

func EqualProbability() algorithm {
	policy := &ChooseRandomOperator{
		operators: []OperatorProbability{
            {PlaceOptimally{}, 1},
			{PlaceOptimallyInRandomVehicle{}, 1},
            {PlaceRandomly{}, 1},
            {PlaceFiveCallsRandomly{}, 1},
		},
        name: "Equal Probability",
	}
    return SimulatedAnnealing(policy)
}

func Moderate() algorithm {
	policy := &ChooseRandomOperator{
		operators: []OperatorProbability{
            {PlaceOptimally{}, 1},
			{PlaceOptimallyInRandomVehicle{}, 2},
            {PlaceRandomly{}, 2},
            {PlaceFiveCallsRandomly{}, 1},
		},
        name: "Moderate",
	}
    return SimulatedAnnealing(policy)
}

func Adventurous() algorithm {
	policy := &ChooseRandomOperator{
		operators: []OperatorProbability{
            {PlaceOptimally{}, 1},
			{PlaceOptimallyInRandomVehicle{}, 1},
            {PlaceRandomly{}, 2},
            {PlaceFiveCallsRandomly{}, 2},
		},
        name: "Adventurous",
	}
    return SimulatedAnnealing(policy)
}

func Intense() algorithm {
	policy := &ChooseRandomOperator{
		operators: []OperatorProbability{
            {PlaceOptimally{}, 2},
			{PlaceOptimallyInRandomVehicle{}, 2},
            {PlaceRandomly{}, 1},
            {PlaceFiveCallsRandomly{}, 1},
		},
        name: "Intense",
	}
    return SimulatedAnnealing(policy)
}

func Extreme() algorithm {
	policy := &ChooseRandomOperator{
		operators: []OperatorProbability{
            {PlaceOptimally{}, 2},
			{PlaceOptimallyInRandomVehicle{}, 1},
            {PlaceRandomly{}, 1},
            {PlaceFiveCallsRandomly{}, 2},
		},
        name: "Extreme",
	}
    return SimulatedAnnealing(policy)
}
