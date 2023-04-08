package basictest

import (
	agent "OCSS/FoosballGeneticLearning/pkg/Agent"
	systemstate "OCSS/FoosballGeneticLearning/pkg/SystemState"

	"gonum.org/v1/gonum/mat"
)

const (
	NUM_PERCEPTS              = 1
	NUM_ACTIONS               = 10
	NUM_AGENTS_PER_SIMULATION = 1
)

type TestSystem struct {
}

func (system *TestSystem) NumPercepts() int {
	return NUM_PERCEPTS
}
func (system *TestSystem) NumActions() int {
	return NUM_ACTIONS
}
func (system *TestSystem) NumAgentsPerSimulation() int {
	return NUM_AGENTS_PER_SIMULATION
}

func (system *TestSystem) InitializeState() *systemstate.SystemState {
	return &systemstate.SystemState{
		StateVector: mat.NewVecDense(system.NumPercepts(), []float64{1.0}),
	}
}

func (system *TestSystem) AdvanceState(state *systemstate.SystemState, agents []*agent.Agent) {
	agentActions := agent.GetAllAgentActions(agents, state.StateVector)
	state.TerminalState = true
	agentScore := 0.0
	for _, elem := range agentActions[0].RawVector().Data {
		agentScore += elem
	}
	agents[0].Score += agentScore
}
