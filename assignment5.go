package main

func Adaptive() algorithm {
    operators := []Operator{
        PlaceOptimally{},
        PlaceOptimallyInRandomVehicle{},
        PlaceRandomly{},
        PlaceFiveCallsRandomly{},
    }
    policy := NewLecturePolicy(50, 0.1, operators)
    return SimulatedAnnealing(policy)
}
