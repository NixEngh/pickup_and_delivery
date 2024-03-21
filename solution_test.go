package main

import (
    "testing"
)

func TestCompareMoveVehicle(t *testing.T) {
    p, err := LoadProblem("./Data/Call_7_Vehicle_3.txt")
    if err != nil {
        t.Error(err)
    }

    solution := p.GenerateInitialSolution()

    oneIndices := FindIndices(solution.Solution, 1)

    solution.MoveInSolution(oneIndices[1][0],0)
    solution.MoveInSolution(oneIndices[1][1],1)
    if !solution.VehiclesToCheckCost[1] {
        t.Errorf("Expected true, got false")
        t.Log(solution.VehiclesToCheckCost)
        t.Log("Solution: ", solution.Solution)
    }
}
