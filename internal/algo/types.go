package algo

import (
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type Algorithm func(problem *problem.Problem) (bestSolution *solution.Solution, bestCost int)
