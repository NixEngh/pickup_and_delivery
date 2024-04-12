package main

func Assignment3() algorithm {
	policy := &ChooseRandomOperator{
		operators: []OperatorScore{
			{operator: OldOneReinsert{}, probability: 1},
		},
		name: "Old OneReinsert",
	}
	return SimulatedAnnealing(policy)
}
