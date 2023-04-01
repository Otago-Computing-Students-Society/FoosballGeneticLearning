package agent

import (
	system "OCSS/FoosballGeneticLearning/pkg/SystemState"

	"gonum.org/v1/gonum/mat"
)

type Agent struct {
	Chromosome   *mat.Dense
	AgentHistory system.StateHistory
}

type AgentAction struct {
	Action *mat.VecDense
}
