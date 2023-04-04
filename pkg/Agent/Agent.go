package agent

import (
	systemstate "OCSS/FoosballGeneticLearning/pkg/SystemState"

	"gonum.org/v1/gonum/mat"
)

type Agent struct {
	Chromosome   *mat.Dense
	AgentHistory systemstate.StateHistory
}

type AgentAction struct {
	Action *mat.VecDense
}
