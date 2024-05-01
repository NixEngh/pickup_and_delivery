package policy

import (
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type AdaptivePolicy struct {
	operators []*OperatorStruct
}

func NewAdaptivePolicy(operators ...operator.Operator) *AdaptivePolicy {

	operatorStructs := make([]*OperatorStruct, len(operators))

	for i, op := range operators {
		operatorStructs[i] = &OperatorStruct{
			Operator: op,
		}
	}

	policy := AdaptivePolicy{
		operators: operatorStructs,
	}

	return &policy

}

func (p *AdaptivePolicy) Apply(s *solution.Solution) {
	os := ChooseWeightedOperator(p.operators)
	operator := os.Operator
	operator.Apply(s)
}

func (p *AdaptivePolicy) Name() string {
	return "AdaptivePolicy"
}
