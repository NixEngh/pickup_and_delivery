package operator

import "github.com/NixEngh/pickup_and_delivery/internal/solution"

type CombineOperator struct {
	Name    string
	nCalls  int
	removal Removal
	insert  Insert
}

func NewCombineOperator(nCalls int, removal Removal, insert Insert, name string) *CombineOperator {
	return &CombineOperator{name, nCalls, removal, insert}
}

func (c *CombineOperator) Apply(s *solution.Solution) int {
	calls := c.removal.Choose(s)
	c.insert.Reinsert(s, calls)
	return len(calls)
}
