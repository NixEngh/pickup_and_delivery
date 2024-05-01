package run

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/policy"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
)

func RunFinalExam(problems []*problem.Problem) {
	algorithms := map[string]algo.Algorithm{
		"Mix": Mix(),
	}

	Run(algorithms, problems)

}

func Optimal() algo.Algorithm {
	pol := policy.NewAdaptivePolicy(
		operator.NewCombineOperator(
			operator.NewRemoveCostly(5),
			operator.NewInsertGreedy(),
			"RemoveCostly + InsertGreedy",
		),
	)

	iterations := 1000

	acceptor := algo.NewIterationR2RAcceptor(iterations)

	stopper := algo.NewIterationBasedStopper(iterations)

	return algo.ALNS(pol, acceptor, stopper)

}

func Mix() algo.Algorithm {
	pol := policy.NewAdaptivePolicy(
		operator.NewCombineOperator(
			operator.NewRemoveRandom(5),
			operator.NewInsertGreedy(),
			"RemoveRandom + InsertGreedy",
		),

		operator.NewCombineOperator(
			operator.NewRemoveCostly(5),
			operator.NewInsertGreedy(),
			"RemoveCostly + InsertGreedy",
		),
	)
	iterations := 1000
	acceptor := algo.NewIterationR2RAcceptor(iterations)
	stopper := algo.NewIterationBasedStopper(iterations)

	return algo.ALNS(pol, acceptor, stopper)
}
