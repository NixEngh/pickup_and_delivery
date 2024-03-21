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

	found := matchAnySlice(solution, expectedOptions)
	if !found {
		t.Errorf("No match found")
	}

}

func TestMoveCallInVehicle(t *testing.T) {
	p := Problem{NumberOfVehicles: 1, NumberOfCalls: 2}

	solution := []int{1, 2, 1, 2, 0}
	expectedOptions := [][]int{
		{2, 1, 1, 2, 0},
		{1, 1, 2, 2, 0},
		{1, 2, 2, 1, 0},
	}
	moveCallInVehicle(&p, solution, []int{0, 2}, []int{4})

	found := matchAnySlice(solution, expectedOptions)
	if !found {
		t.Errorf("No match found")
	}
}

func TestOneReinsert(t *testing.T) {
    p := Problem{NumberOfVehicles: 2, NumberOfCalls: 1}

    solution := []int{0,1,1,0}
    expected := []int{0,0,1,1}

    OneReinsert(&p, solution)
    
    if !matchSlice(solution, expected) {
        t.Errorf("Expected %v, got %v", expected, solution)
    }
}

