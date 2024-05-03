package run

import (
	"fmt"
	"sync"

	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/operator"
	"github.com/NixEngh/pickup_and_delivery/internal/policy"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
)

func RunFinalExam(problems []*problem.Problem) {
	algos := []func() algo.Algorithm{
		Attempt1,
		Attempt2,
		Attempt3,
		Attempt4,
	}

	wg := sync.WaitGroup{}
	for i, al := range algos {
		wg.Add(1)
		go func(i int, algo func() algo.Algorithm) {
			defer wg.Done()
			RunUltimate(problems, fmt.Sprintf("attempt_%d", i), algo)
		}(i, al)
	}

	wg.Wait()
}

var timer int64 = 60 * 10

func Attempt1() algo.Algorithm {

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

	acceptor := algo.NewTimeR2RAcceptor(timer)

	stopper := algo.NewTimeBasedStopper(timer)

	return algo.ALNS(pol, acceptor, stopper)
}

func Attempt2() algo.Algorithm {

	pol := policy.NewAdaptivePolicy(

		operator.NewCombineOperator(
			operator.NewRemoveRandom(5),
			operator.NewInsertGreedy(),
			"Operator 1",
		),
		operator.NewCombineOperator(
			operator.NewRemoveRandom(20),
			operator.NewInsertGreedy(),
			"Operator 2",
		),
		operator.NewCombineOperator(
			operator.NewRemoveRandom(40),
			operator.NewInsertGreedy(),
			"Operator 4",
		),
		operator.NewCombineOperator(
			operator.NewRemoveCostly(5),
			operator.NewInsertGreedy(),
			"Operator 5",
		),
		operator.NewCombineOperator(
			operator.NewRemoveCostly(20),
			operator.NewInsertGreedy(),
			"Operator 6",
		),
	)

	acceptor := algo.NewTimeR2RAcceptor(timer)

	stopper := algo.NewTimeBasedStopper(timer)

	//acceptor := algo.NewIterationR2RAcceptor(10000)
	//stopper := algo.NewIterationBasedStopper(10000)

	return algo.ALNS(pol, acceptor, stopper)
}

func Attempt3() algo.Algorithm {

	pol := policy.NewAdaptivePolicy(

		operator.NewCombineOperator(
			operator.NewRemoveRandom(5),
			operator.NewInsertGreedy(),
			"Operator 1",
		),
		operator.NewCombineOperator(
			operator.NewRemoveRandom(20),
			operator.NewInsertGreedy(),
			"Operator 7",
		),
		operator.NewCombineOperator(
			operator.NewRemoveCostly(40),
			operator.NewInsertGreedy(),
			"Operator 7",
		),
	)

	acceptor := algo.NewTimeR2RAcceptor(timer)

	stopper := algo.NewTimeBasedStopper(timer)

	return algo.ALNS(pol, acceptor, stopper)
}

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
