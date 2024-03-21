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
    solution.UpdateCosts()

    if solution.Cost() != solution.CostFunction() {
        t.Errorf("Generated solution has different costs")
    }

    var cost int
    var realCost int

    for i := 0; i < 20; i++ {
        solution.OneReinsert()
        if len(solution.VehiclesToCheckCost) != 1 {
            t.Errorf("Expected 1, got %d", len(solution.VehiclesToCheckCost))
        }
        cost = solution.Cost()
        if len(solution.VehiclesToCheckCost) != 0 {
            t.Errorf("Expected 0, got %d", len(solution.VehiclesToCheckCost))
        }
        realCost = solution.CostFunction()

        if cost != realCost {
            t.Errorf("Expected %d, got %d, for i %d", realCost, cost, i)
        }
    }
}

func TestCostsForRandomSolutions(t *testing.T) {
    p, err := LoadProblem("./Data/Call_7_Vehicle_3.txt")
    if err != nil {
        t.Error(err)
    }

    for i := 0; i < 20; i++ {
        solution := p.GenerateRandomSolution()
        cost := solution.Cost()
        realCost := solution.CostFunction()

        if cost != realCost {
            t.Errorf("Expected %d, got %d", realCost, cost)
        }
    }
}

func TestUpdateCosts(t *testing.T) {
    p, err := LoadProblem("./Data/Call_7_Vehicle_3.txt")
    if err != nil {
        t.Error(err)
    }

    solution := p.GenerateInitialSolution()


    cost_before := solution.Cost()
    solution.UpdateCosts()
    cost_after := solution.Cost()

    if cost_before != cost_after {
        t.Errorf("Cost_before %d, Cost_after %d", cost_before, cost_after)
    }

    twoIndices := FindIndices(solution.Solution, 4)[4]
    solution.MoveInSolution(twoIndices[0], 0)
    solution.MoveInSolution(twoIndices[1], 0)
    solution.UpdateCosts()

    if solution.Cost() != solution.CostFunction() {
        t.Errorf("Costs are not equal")
    }

    zeroIndices := FindIndices(solution.Solution, 0)[0]
    solution.MoveInSolution(0, zeroIndices[len(zeroIndices)-1])
    solution.MoveInSolution(0, zeroIndices[len(zeroIndices)-1])
    solution.UpdateCosts()

    if solution.Cost() != solution.CostFunction() {
        t.Errorf("Costs are not equal")
    }
}

func TestCostCheckVehiclesAreUpdated(t *testing.T) {
    p, err := LoadProblem("./Data/Call_7_Vehicle_3.txt")
    if err != nil {
        t.Error(err)
    }

    solution := p.GenerateInitialSolution()

    if len(solution.VehiclesToCheckCost) != 0 {
        t.Error("Mistake in GenerateInitialSolution", len(solution.VehiclesToCheckCost))
    }

    solution.OneReinsert()


    if len(solution.VehiclesToCheckCost) != 1 {
        t.Error("Mistake in OneReinsert", len(solution.VehiclesToCheckCost))
        t.Log(solution.VehiclesToCheckCost)
    }
}

func TestVehicleCostFunction(t *testing.T) {
    p, err := LoadProblem("./Data/Call_7_Vehicle_3.txt")
    if err != nil {
        t.Error(err)
    }

    solution := p.GenerateInitialSolution()
    oneIndices := FindIndices(solution.Solution, 1)[1]

    solution.MoveInSolution(oneIndices[0], 0)
    solution.MoveInSolution(oneIndices[1], 0)

    expectedCost := solution.CostFunction()
    calculatedCost := solution.OutSourceCostFunction() + solution.VehicleCostFunction(1)

    if expectedCost != calculatedCost {
        t.Errorf("Expected %d, got %d", expectedCost, calculatedCost)
    }
}

func TestOutSourceCost(t *testing.T) {
    p, err := LoadProblem("./Data/Call_7_Vehicle_3.txt")
    if err != nil {
        t.Error(err)
    }

    solution := p.GenerateInitialSolution()
    oneIndices := FindIndices(solution.Solution, 1)[1]

    solution.MoveInSolution(oneIndices[0], 0)
    solution.MoveInSolution(oneIndices[1], 0)

    twoIndices := FindIndices(solution.Solution, 2)[2]
    solution.MoveInSolution(twoIndices[0], 0)
    solution.MoveInSolution(twoIndices[1], 0)

    solution.UpdateCosts()

    if solution.OutSourceCost != solution.OutSourceCostFunction() {
        t.Errorf("Expected %d, got %d", solution.OutSourceCostFunction(), solution.OutSourceCost)
    }
}
