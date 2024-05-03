package run

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/policy"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
)

func RunFinalExam(problems []*problem.Problem) {

	RunUltimate(problems, "Finalrun", Attempt4)
}

var timer int64 = 60 * 15

func Attempt4() algo.Algorithm {

	pol := policy.NewAdaptivePolicy(

		operator.NewCombineOperator(
			operator.NewRemoveRandom(5),
			operator.NewInsertGreedy(),
			"Operator 1",
		),
		operator.NewCombineOperator(
			operator.NewRemoveRandom(10),
			operator.NewInsertGreedy(),
			"Operator 7",
		),
		operator.NewCombineOperator(
			operator.NewRemoveRandom(15),
			operator.NewInsertGreedy(),
			"Operator 4",
		),
		operator.NewCombineOperator(
			operator.NewRemoveCostly(5),
			operator.NewInsertGreedy(),
			"Operator 5",
		),
		operator.NewCombineOperator(
			operator.NewRemoveCostly(10),
			operator.NewInsertGreedy(),
			"Operator 6",
		),
		operator.NewCombineOperator(
			operator.NewRemoveCostly(15),
			operator.NewInsertGreedy(),
			"Operator 7",
		),
	)

	acceptor := algo.NewTimeR2RAcceptor(timer)

	stopper := algo.NewTimeBasedStopper(timer)

	return algo.ALNS(pol, acceptor, stopper)
}
