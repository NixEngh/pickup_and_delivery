package main

import (
	"testing"
)

func TestGenerateInitialSolution(t *testing.T) {
	p := Problem{NumberOfVehicles: 2, NumberOfCalls: 3}
	solution := p.GenerateInitialSolution()
	expected := []int{0, 0, 1, 1, 2, 2, 3, 3}
	for i := range solution {
		if solution[i] != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], solution[i])
		}
	}
}

func TestGenerateRandomSolution(t *testing.T) {
	p := Problem{NumberOfVehicles: 2, NumberOfCalls: 3}
	solution := p.GenerateRandomSolution()

	if len(solution) != 8 {
		t.Errorf("Expected 8, got %d", len(solution))
	}
}

func TestMoveFromOutsource(t *testing.T) {
	p := Problem{NumberOfVehicles: 1, NumberOfCalls: 1, CallVehicleMap: map[int][]int{1: {1}}}

	solution := []int{0, 1, 1}
	expected := []int{1, 1, 0}

	moveFromOutsource(&p, solution, []int{1, 2}, []int{0})

	for i := range solution {
		if solution[i] != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], solution[i])
		}
	}

	p = Problem{NumberOfVehicles: 2, NumberOfCalls: 2, CallVehicleMap: map[int][]int{1: {1, 2}, 2: {2}}}

	solution = p.GenerateInitialSolution()
    expectedOptions := [][]int{
		{0, 1, 1, 0, 2, 2},
		{1, 1, 0, 0, 2, 2},
		{0, 2, 2, 0, 1, 1},
		{2, 2, 0, 0, 1, 1},
	}

    moveFromOutsource(&p, solution, []int{2, 3}, []int{0, 1})

    match := [4]bool{true, true, true, true}

    for i, element := range solution {
        possible := false
        for j := range expectedOptions {
            if element != expectedOptions[j][i] {
                match[j] = false
            } else {
                possible = true
            }
        }
        if !possible {
            t.Errorf("Unexpected element %d", element)
        }
    }
    
    if !match[0] && !match[1] && !match[2] && !match[3] {
        t.Errorf("No match found")
    }

}
