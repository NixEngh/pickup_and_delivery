package main

import (
    "testing"
)

func TestCompareCostFunctions(t *testing.T) {
    p, err := LoadProblem("./Data/Call_7_Vehicle_3.txt")
    if err != nil {
        t.Error(err)
    }

    solution := p.GenerateInitialSolution()

    previousCost := p.CostFunction(solution)

    for i := 0; i < 100; i++ {
        changedIndices := OneReinsert(&p, solution)
        cost := p.CostFunction(solution)
        indexedCost := p.IndexedCostFunction(solution,previousCost, changedIndices)
        if cost != indexedCost {
            t.Errorf("Cost function and indexed cost function do not return the same value: %d != %d", cost, indexedCost)
        }
        previousCost = cost
    }

}


