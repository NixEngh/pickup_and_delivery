package policy

import (
	"math/rand"

	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type ChooseRandomOperator struct {
	Operators []*OperatorStruct
	name      string
}

func NewChooseRandomOperator(operators []*OperatorStruct, name string) *ChooseRandomOperator {

	return &ChooseRandomOperator{
		Operators: operators,
		name:      name,
	}
}

func (c *ChooseRandomOperator) Apply(s *solution.Solution) {
	os := ChooseWeightedOperator(c.Operators)
	operator := os.Operator
	operator.Apply(s)
}

func (c *ChooseRandomOperator) UpdateProbabilities(s *solution.Solution) {
	return
}

func (c *ChooseRandomOperator) Name() string {
	return c.name
}

func ChooseWeightedOperator(operators []*OperatorStruct) *OperatorStruct {
	var total float64
	for _, op := range operators {
		total += op.Probability
	}

	r := rand.Float64() * total
	for _, op := range operators {
		if r -= op.Probability; r < 0 {
			return op
		}
	}

	panic("Should not reach here")
}
