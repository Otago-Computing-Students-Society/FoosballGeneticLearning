package agent

import (
	systemstate "OCSS/FoosballGeneticLearning/pkg/SystemState"

	"gonum.org/v1/gonum/mat"
)

type AgentAction *mat.VecDense

type Agent struct {
	Chromosome   *mat.Dense
	AgentHistory systemstate.StateHistory
	Score        float64
}

func NewAgent(chromosome *mat.Dense) *Agent {
	return &Agent{
		Chromosome:   chromosome,
		AgentHistory: make(systemstate.StateHistory, 0),
		Score:        0.0,
	}
}

func (agent *Agent) GetAction(stateVector *mat.VecDense) AgentAction {
	numActions, _ := agent.Chromosome.Dims()
	actionVector := mat.NewVecDense(numActions, nil)
	actionVector.MulVec(agent.Chromosome, stateVector)
	return actionVector
}
