package policy

import (
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type AdaptivePolicy struct {
	operators []OperatorScore
}

func NewAdaptivePolicy() *AdaptivePolicy {
	policy := AdaptivePolicy{
		operators: []OperatorScore{
			OperatorScore{
				Operator: operator.NewCombineOperator(
					5,
					operator.RemoveRandom{},
					operator.InsertGreedy{},
					"RemoveRandom + InsertGreedy",
				),
				Probability: 1,
			},
		},
	}

}

func (p AdaptivePolicy) Apply(s *solution.Solution) {

}

func (p AdaptivePolicy) Name() string {
	return ""
}
