package operator

import "github.com/NixEngh/pickup_and_delivery/internal/solution"

type CombineOperator struct {
	Name    string
	removal Removal
	insert  Insert
}

func NewCombineOperator(removal Removal, insert Insert, name string) *CombineOperator {
	return &CombineOperator{name, removal, insert}
}

func (c *CombineOperator) Apply(s *solution.Solution) int {
	calls := c.removal.Choose(s)
	c.insert.Reinsert(s, calls)
	return len(calls)
}
