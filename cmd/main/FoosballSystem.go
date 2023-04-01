package main

import (
	agent "OCSS/FoosballGeneticLearning/pkg/Agent"
	systemstate "OCSS/FoosballGeneticLearning/pkg/SystemState"
)

const (
	NUM_PERCEPTS = 8
	NUM_ACTIONS  = 8
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

// Defines the scoring of an agents state history.
//
// TODO(hayden): Implement a score function that makes sense for the system implementation
//
// # Arguments
//
// - `history`: the StateHistory of an agent to score
//
// # Returns
//
// A float representing the "score" of the agent represented by `history`.
// This float should be larger for a more "fit" agent.
func (system *FoosballSystem) ScoreFunction(history systemstate.StateHistory) float64 {
	return 0.0
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
func (system *FoosballSystem) AdvanceState(state systemstate.SystemState, agentActions []agent.AgentAction) systemstate.SystemState {
	return nil
}
