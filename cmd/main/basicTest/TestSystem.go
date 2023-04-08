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

// Returns the initial state of the system
//
// In this case it is very boring - the state is always 1.0 flat
func (system *TestSystem) InitializeState() *systemstate.SystemState {
	return &systemstate.SystemState{
		StateVector:   mat.NewVecDense(system.NumPercepts(), []float64{1.0}),
		TerminalState: false,
	}
}

// Defines the behavior of the system
//
// This method should take the current system state and agents, find the agent actions,
// then apply those actions to the state to update it.
//
// This method is also responsible for updating the agent scores!
//
// The implementation here is very boring - we immediately call the simulation done (terminal = true)
// and set the score of the agent to the sum of the action vector.
func (system *TestSystem) AdvanceState(state *systemstate.SystemState, agents []*agent.Agent) {
	agentActions := agent.GetAllAgentActions(agents, state.StateVector)
	state.TerminalState = true
	agentScore := 0.0
	for _, elem := range agentActions[0].RawVector().Data {
		agentScore += elem
	}
	agents[0].Score += agentScore
}
