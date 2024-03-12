package main

import (
	"math/rand"
)

func (p *Problem) GenerateInitialSolution() []int {
	solution := make([]int, p.NumberOfVehicles)
	for i := 0; i < p.NumberOfVehicles; i++ {
		solution[i] = 0
	}
    for i := 1; i <= p.NumberOfCalls; i++ {
        solution = append(solution, i)
        solution = append(solution, i)
    }

	return solution
}

func (p *Problem) GenerateRandomSolution() []int {
	vehicles := make([][]int, p.NumberOfVehicles+1)
	for i := 0; i <= p.NumberOfVehicles; i++ {
		vehicles[i] = make([]int, 0)
	}

	order := rand.Perm(p.NumberOfCalls)
	for _, call := range order {
		call += 1
		vehicle := rand.Intn(p.NumberOfVehicles + 1)
		vehicles[vehicle] = append(vehicles[vehicle], call)
		vehicles[vehicle] = append(vehicles[vehicle], call)
	}

	for i := 1; i <= p.NumberOfVehicles; i++ {
		rand.Shuffle(len(vehicles[i]), func(x, y int) {
			vehicles[i][x], vehicles[i][y] = vehicles[i][y], vehicles[i][x]
		})
		vehicles[i] = append(vehicles[i], 0)
	}

	solution := make([]int, 0)
	for i := 1; i <= p.NumberOfVehicles; i++ {
		solution = append(solution, vehicles[i]...)
	}

	solution = append(solution, vehicles[0]...)
	return solution
}

func (p *Problem) one_reinsert()
