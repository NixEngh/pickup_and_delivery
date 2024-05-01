package run

import (
	"github.com/NixEngh/pickup_and_delivery/internal/algo"
	"github.com/NixEngh/pickup_and_delivery/internal/policy"
	"github.com/NixEngh/pickup_and_delivery/internal/problem"
)

func RunFinalExam(problems []*problem.Problem) {
	algorithms := map[string]algo.Algorithm{
		"AlwaysOptimal": AlwaysOptimal(),
	}

	Run(algorithms, problems)

}

func AlwaysOptimal() algo.Algorithm {
	pol := policy.NewAdaptivePolicy()
	iterations := 10000
	acceptor := algo.NewIterationR2RAcceptor(iterations)
	stopper := algo.NewIterationBasedStopper(iterations)

	return algo.ALNS(pol, acceptor, stopper)
}
