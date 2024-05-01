package policy

import (
	"math"

	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type AdaptivePolicy struct {
	operators     []*OperatorStruct
	iteration     int
	segmentLength int
	r             float64
	bestCost      int
	compareSet    *CompareSet
}

func NewAdaptivePolicy(operators ...operator.Operator) *AdaptivePolicy {
	operatorStructs := make([]*OperatorStruct, len(operators))

	for i, op := range operators {
		operatorStructs[i] = &OperatorStruct{
			Operator:    op,
			Probability: 1 / float64(len(operators)),
		}
	}

	policy := &AdaptivePolicy{
		operators:     operatorStructs,
		r:             0.1,
		segmentLength: 100,
		bestCost:      math.MaxInt32,
		compareSet:    NewCompareSet(),
	}

	return policy

}

func (p *AdaptivePolicy) Apply(s *solution.Solution) {
	operatorStruct := ChooseWeightedOperator(p.operators)
	operatorStruct.timesUsed++

	operator := operatorStruct.Operator
	costBefore := s.Cost()
	newCost := operator.Apply(s)

	p.UpdateScore(costBefore, newCost, s, operatorStruct)

	if p.iteration == p.segmentLength {
		p.iteration = 0
		p.UpdateProbabilities()
	}
}

func (p *AdaptivePolicy) UpdateScore(costBefore, newCost int, s *solution.Solution, score *OperatorStruct) {
	scoreToAdd := 0
	if !p.compareSet.HasVisitedSolution(s.Solution) {
		scoreToAdd += 1
	}
	if newCost < p.bestCost {
		p.bestCost = newCost
		scoreToAdd += 4
	}

	if newCost < costBefore {
		scoreToAdd += 2
	}

	score.score += scoreToAdd
}

func (p *AdaptivePolicy) UpdateProbabilities() {
	for _, operator := range p.operators {
		operator.Probability = operator.Probability*(1-p.r) + p.r*(float64(operator.score)/float64(operator.timesUsed))
		operator.timesUsed = 0
		operator.score = 0
	}
}

func (p *AdaptivePolicy) Name() string {
	return "AdaptivePolicy"
}
