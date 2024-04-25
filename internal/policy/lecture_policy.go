package policy

import (
	"math"

	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/solution"
)

type LecturePolicy struct {
    operators         []OperatorScore
    name              string
	iteration         int
	segmentLength     int
	r                 float64
	compareSet        CompareSet
	noNewFoundCounter int
	bestCost          int
}

func NewLecturePolicy(segmentLength int, r float64, operators []operator.Operator) *LecturePolicy {

	operatorScores := make([]OperatorScore, 0)

	for _, operator := range operators {
		operatorScores = append(operatorScores, OperatorScore{
			Operator:    operator,
			Probability: 1 / float64(len(operators)),
		})
	}

	return &LecturePolicy{
		name:          "LecturePolicy",
		operators:     operatorScores,
		segmentLength: segmentLength,
		r:             r,
		compareSet:    *NewCompareSet(),
		bestCost:      math.MaxInt32,
	}
}

func (c *LecturePolicy) Apply(s *solution.Solution) {
	hasVisited := c.compareSet.HasVisitedSolution(s.Solution)
	if hasVisited {
		c.noNewFoundCounter++
	} else {
		c.noNewFoundCounter = 0
	}

	if c.noNewFoundCounter > 100 {
        op := operator.PlaceFiveCallsRandomly{}
        op.Apply(s)
	} else {
		os := ChooseWeightedOperator(c.operators)
		os.timesUsed++
		operator := os.Operator
		costBefore := s.Cost()
		newCost := operator.Apply(s)
		c.UpdateScore(costBefore, newCost, hasVisited, s, os)
	}

	c.iteration++
	if c.iteration == c.segmentLength {
		c.iteration = 0
		c.UpdateProbabilities()
	}
}

func (c *LecturePolicy) UpdateScore(costBefore, newCost int, hasVisited bool, s *solution.Solution, score *OperatorScore) {
	scoreToAdd := 0
	if !hasVisited {
		scoreToAdd += 1
	}
	if newCost < c.bestCost {
		c.bestCost = newCost
		scoreToAdd += 4
	}

	if newCost < costBefore {
		scoreToAdd += 1
	}

	score.score += scoreToAdd
}

func (c *LecturePolicy) UpdateProbabilities() {
	for _, operator := range c.operators {
		operator.Probability = operator.Probability*(1-c.r) + c.r*(float64(operator.score)/float64(operator.timesUsed))
		operator.timesUsed = 0
		operator.score = 0
	}
}

func (c *LecturePolicy) Name() string {
	return c.name
}
