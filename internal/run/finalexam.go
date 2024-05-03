package run

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/policy"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
)

func RunFinalExam(problems []*problem.Problem) {
	RunUltimate(problems, AttemptGenerator)
}

func AttemptGenerator() algo.Algorithm {
	return Attempt()
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
			"Operator 7",
		),
		operator.NewCombineOperator(
			operator.NewRemoveRandom(30),
			operator.NewInsertGreedy(),
			"Operator 4",
		),
		operator.NewCombineOperator(
			operator.NewRemoveCostly(10),
			operator.NewInsertGreedy(),
			"Operator 5",
		),
		operator.NewCombineOperator(
			operator.NewRemoveCostly(20),
			operator.NewInsertGreedy(),
			"Operator 6",
		),
		operator.NewCombineOperator(
			operator.NewRemoveCostly(30),
			operator.NewInsertGreedy(),
			"Operator 7",
		),
	)

	var time int64 = 60 * 15

	acceptor := algo.NewTimeR2RAcceptor(time)

	stopper := algo.NewTimeBasedStopper(time)

	//acceptor := algo.NewIterationR2RAcceptor(10000)
	//stopper := algo.NewIterationBasedStopper(10000)

	return algo.ALNS(pol, acceptor, stopper)
}
