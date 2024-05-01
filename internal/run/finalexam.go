package run

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/policy"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
)

func RunFinalExam(problems []*problem.Problem) {
	algorithms := map[string]algo.Algorithm{
		"Attempt": Attempt(),
	}

	Run(algorithms, problems)

}

func Attempt() algo.Algorithm {
	pol := policy.NewAdaptivePolicy(

		operator.NewCombineOperator(
			operator.NewRemoveRandom(10),
			operator.NewInsertGreedy(),
			"Operator 1",
		),
		operator.NewCombineOperator(
			operator.NewRemoveRandom(20),
			operator.NewInsertGreedy(),
			"Operator 3",
		),
		operator.NewCombineOperator(
			operator.NewRemoveRandom(40),
			operator.NewInsertGreedy(),
			"Operator 4",
		),
		operator.NewCombineOperator(
			operator.NewRemoveCostly(10),
			operator.NewInsertGreedy(),
			"Operator 4",
		),
		operator.NewCombineOperator(
			operator.NewRemoveCostly(20),
			operator.NewInsertGreedy(),
			"Operator 5",
		),
	)

	iterations := 10000

	acceptor := algo.NewIterationR2RAcceptor(iterations)

	stopper := algo.NewIterationBasedStopper(iterations)

	return algo.ALNS(pol, acceptor, stopper)
}
