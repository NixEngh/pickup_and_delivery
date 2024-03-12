package main

func RandomSearch(problem *Problem) []int {
    bestSolution := problem.GenerateInitialSolution()
    bestCost := problem.CostFunction(bestSolution)

    for i := 0; i < 10000; i++ {
        solution := problem.GenerateRandomSolution()
        cost := problem.CostFunction(solution)
        if problem.IsFeasible(solution) && cost < bestCost {
            bestSolution = solution
            bestCost = cost
        }
    }

    return bestSolution
}
