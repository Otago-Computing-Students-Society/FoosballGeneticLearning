package agent

import (
	systemstate "github.com/Otago-Computer-Science-Society/Foosball-Genetic-Learning/pkg/SystemState"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

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

// Return a new agent with a unit Gaussian random chromosome
//
// TODO(hayden): Actually implement this method!
func NewRandomGaussianAgent(numActions int, numPercepts int) *Agent {
	unitNormal := distuv.UnitNormal
	chromosomeData := make([]float64, numActions*numPercepts)
	for index := range chromosomeData {
		chromosomeData[index] = unitNormal.Rand()
	}
	chromosome := mat.NewDense(numActions, numPercepts, chromosomeData)
	return NewAgent(chromosome)
}

func (agent *Agent) GetAction(stateVector *mat.VecDense) *mat.VecDense {
	numActions, _ := agent.Chromosome.Dims()
	actionVector := mat.NewVecDense(numActions, nil)
	actionVector.MulVec(agent.Chromosome, stateVector)
	return actionVector
}

func GetAllAgentActions(agents []*Agent, stateVector *mat.VecDense) []*mat.VecDense {
	agentActions := make([]*mat.VecDense, len(agents))
	for agentIndex := range agents {
		agentActions[agentIndex] = agents[agentIndex].GetAction(stateVector)
	}
	return agentActions
}
