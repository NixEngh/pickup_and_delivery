package operator

import (
	"math/rand"

	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)
type PlaceFiveCallsRandomly struct{}

func (o PlaceFiveCallsRandomly) Apply(s *solution.Solution) int {
	callsToMove := rand.Perm(s.Problem.NumberOfCalls)
	count := 0
	for _, callToMove := range callsToMove {
		if ok := s.PlaceCallRandomly(callToMove + 1); ok {
			count += 1
		}

		if count == 5 {
			break
		}
	}
	return s.Cost()
}
