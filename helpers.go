package main

func findIndices(slice []int, value int) []int {
    var indices []int
    for index, element := range slice {
        if element == value {
            indices = append(indices, index)
        }
    }
    return indices
}

func getTour(solution []int, vehicleIndex int) []int {
    zeroIndexes := findIndices(solution, 0)
    start_ind := 0

    if vehicleIndex != 1 {
        start_ind = zeroIndexes[vehicleIndex-2] + 1
    }
    end_ind := zeroIndexes[vehicleIndex-1]
    
    return solution[start_ind:end_ind]
}
