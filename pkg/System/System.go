package system

import (
	agent "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/Agent"
	systemstate "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/SystemState"
)

type System interface {

	// Gets the number of percepts for this system.
	NumPercepts() int

	// Gets the number of actions for this system.
	NumActions() int

	// Gets the number of agents per simulation for this system.
	NumAgentsPerSimulation() int

	// Defines the function to create an initial state of the system.
	InitializeState() *systemstate.SystemState

	// Defines the function to advance the system state forward a step,
	// given the agents in this state.
	//
	// Note you must get the actions of the agents manually (see `agent/GetAllAgentActions`)
	//
	// Note you must manually update the score of each agent in this method. This is vital!!
	//
	// # Arguments
	//
	// state: The state the system is *currently* in. Prior state.
	//
	// agents : The agents that are operating in this state. Note this is an array
	// so we can have flexibility in the number of agents involved in a system.
	AdvanceState(*systemstate.SystemState, []*agent.Agent)
}
