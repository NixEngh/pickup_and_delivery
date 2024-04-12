package main

import (
	"fmt"
	"math"
	"math/rand"
)

type OperatorScore struct {
	operator    Operator
	probability float64
	score       int
	timesUsed   int
}

func ChooseWeightedOperator(operators []OperatorScore) *OperatorScore {
	var total float64
	for _, op := range operators {
		total += op.probability
	}

	r := rand.Float64() * total
	for _, op := range operators {
		if r -= op.probability; r < 0 {
			return &op
		}
	}

	panic("Should not reach here")
}

type OperatorPolicy interface {
	Apply(s *Solution)
	Name() string
}

type ChooseRandomOperator struct {
	operators []OperatorScore
	name      string
}

func (c *ChooseRandomOperator) Apply(s *Solution) {
	os := ChooseWeightedOperator(c.operators)
	operator := os.operator
	operator.Apply(s)
}

func (c *ChooseRandomOperator) UpdateProbabilities(s *Solution) {
	return
}

func (c *ChooseRandomOperator) Name() string {
	return c.name
}

type LecturePolicy struct {
	operators     []OperatorScore
	name          string
	iteration     int
	segmentLength int
	r             float64
	compareSet   CompareSet
	bestCost      int
}

func NewLecturePolicy(segmentLength int, r float64, operators []Operator) *LecturePolicy {

	operatorScores := make([]OperatorScore, 0)

	for _, operator := range operators {
		operatorScores = append(operatorScores, OperatorScore{
			operator:    operator,
			probability: 1 / float64(len(operators)),
		})
	}

	return &LecturePolicy{
		name:          "LecturePolicy",
		operators:     operatorScores,
		segmentLength: segmentLength,
		r:             r,
		compareSet:   *NewCompareSet(),
		bestCost:      math.MaxInt32,
	}
}

func (c *LecturePolicy) Apply(s *Solution) {
	os := ChooseWeightedOperator(c.operators)
	os.timesUsed++

	operator := os.operator
	costBefore := s.Cost()
	newCost := operator.Apply(s)
	c.UpdateScore(costBefore, newCost, s, os)
	c.iteration++
	if c.iteration == c.segmentLength {
		c.iteration = 0
		c.UpdateProbabilities()
	}
}

func (c *LecturePolicy) UpdateScore(costBefore, newCost int, s *Solution, score *OperatorScore) {
	scoreToAdd := 0
    if !c.compareSet.HasVisitedSolution(s.Solution) {
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
		operator.probability = operator.probability*(1-c.r) + c.r*(float64(operator.score)/float64(operator.timesUsed))
		operator.timesUsed = 0
		operator.score = 0
	}
}

func (c *LecturePolicy) Name() string {
	return c.name
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

func NewCompareSet() *CompareSet {
    return &CompareSet{
        visited: make(map[string]bool),
    }
}
