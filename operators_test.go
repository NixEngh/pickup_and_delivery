package main

import (
	"testing"
)

func TestGenerateInitialSolution(t *testing.T) {
	p, err := LoadProblem("./Data/Call_7_Vehicle_3.txt")
	if err != nil {
		t.Error(err)
	}

	s := p.GenerateInitialSolution()
	solution := s.Solution

	expected := []int{0, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7}
	for i := range solution {
		if solution[i] != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], solution[i])
		}
	}
}

func TestGenerateRandomSolution(t *testing.T) {
	p, err := LoadProblem("./Data/Call_7_Vehicle_3.txt")
	if err != nil {
		t.Error(err)
	}

	s := p.GenerateRandomSolution()

	for vehicle := 1; vehicle <= p.NumberOfVehicles; vehicle++ {
		if len(GetTour(s.Solution, vehicle))%2 != 0 {
			t.Error(s.Solution)
		}
	}
}

func TestMoveCallInVehicle(t *testing.T) {
	p, err := LoadProblem("./Data/Call_7_Vehicle_3.txt")
	if err != nil {
		t.Error(err)
	}

	s := p.GenerateInitialSolution()
	indices := FindIndices[int](s.Solution, 0, 1, 2)
	s.MoveInSolution(indices[2][0], indices[0][0])
	s.MoveInSolution(indices[2][1], indices[0][0])
	s.MoveInSolution(indices[1][0], indices[0][0])
	s.MoveInSolution(indices[1][1], indices[0][0])

	expected := []int{1, 2, 1, 2}
	tour := GetTour(s.Solution, 1)

	for i := range tour {
		if tour[i] != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], tour[i])
		}
	}
}
