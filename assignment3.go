package main

func Assignment3() algorithm {
	policy := &ChooseRandomOperator{
		operators: []OperatorProbability{
			{OldOneReinsert{}, 1},
		},
        name: "Old OneReinsert",
    }
    return SimulatedAnnealing(policy)
}
