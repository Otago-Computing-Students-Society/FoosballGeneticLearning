package foosballsystem

import (
	agent "OCSS/FoosballGeneticLearning/pkg/Agent"
	systemstate "OCSS/FoosballGeneticLearning/pkg/SystemState"
)

const (
	NUM_PERCEPTS              = 8
	NUM_ACTIONS               = 8
	NUM_AGENTS_PER_SIMULATION = 2
)

type FoosballSystem struct {
}

// Gets the number of percepts for this system.
func (system *FoosballSystem) NumPercepts() int {
	return NUM_PERCEPTS
}

// Gets the number of actions for this system.
func (system *FoosballSystem) NumActions() int {
	return NUM_ACTIONS
}

// Gets the number of agents per simulation for this system.
func (system *FoosballSystem) NumAgentsPerSimulation() int {
	return NUM_AGENTS_PER_SIMULATION
}

// Defines the function to create an initial state of the system.
//
// For Foosball, this should include all rods in some neutral position and a ball
// in the center of the table with some random perturbation in position of velocity.
//
// TODO(hayden): Implement this function correctly for our system.
func (system *FoosballSystem) InitializeState(state *systemstate.SystemState, agentActions []agent.AgentAction) *systemstate.SystemState {
	return &systemstate.SystemState{}
}

// Defines the function to advance the system state forward a step,
// given the agent actions in this state.
//
// TODO(hayden): Implement this function correctly for our system.
// Note this function is essentially the entire simulation!
//
// # Arguments
//
// - `state`: The state the system is *currently* in. Prior state.
//
// - `agentActions`: The actions taken by all agents in `state`. Note this is an array
// so we can have flexibility in the number of agents involved in a system.
//
// # Returns
//
// A SystemState representing the updated state given the prior state and agent actions
func (system *FoosballSystem) AdvanceState(state *systemstate.SystemState, agentActions []agent.AgentAction) *systemstate.SystemState {
	return &systemstate.SystemState{}
}
