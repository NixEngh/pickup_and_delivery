package policy

import (
	"fmt"

	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type OperatorStruct struct {
	Operator    operator.Operator
	Probability float64
	score       int
	timesUsed   int
}

func NewOperatorStruct(o operator.Operator, p float64) *OperatorStruct {
	return &OperatorStruct{
		Operator:    o,
		Probability: p,
	}
}

type OperatorPolicy interface {
	Apply(s *solution.Solution)
	Name() string
}

type CompareSet struct {
	visited map[string]bool
}

func (cs *CompareSet) HasVisitedSolution(solution []int) bool {
	hash := hashSolution(solution)
	if _, exists := cs.visited[hash]; exists {
		return true
	}
	cs.visited[hash] = true
	return false
}

func hashSolution(solution []int) string {
	hash := ""
	for _, v := range solution {
		hash += fmt.Sprintf(":%v", v) // Simple concatenation, consider a better hashing function
	}
	return hash
}

func NewCompareSet() CompareSet {
	return CompareSet{
		visited: make(map[string]bool),
	}
}
