package main

import (
	"fmt"
)

func main() {
    problem, _ := load_problem("../Data/Call_7_Vehicle_3.txt")
	//solution := []int{13, 13, 0, 8, 10, 10, 17, 17, 8, 0, 3, 3, 4, 4, 15, 15, 11, 12, 11, 16, 16, 12, 0, 9, 9, 5, 5, 14, 14, 0, 7, 7, 1, 1, 0, 6, 18, 18, 2, 6, 2}

	fmt.Println(problem.GenerateInitialSolution())
    fmt.Println(problem.GenerateRandomSolution())

}
