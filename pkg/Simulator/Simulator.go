package simulator

import (
	agent "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/Agent"
	datacollector "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/DataCollector"
	system "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/System"
)

const MAXIMUM_SIMULATION_ITERATIONS = 5000

// Define a wrapper around the simulation functions to allow for easy concurrency
//
// This wrapper function is intended to be run as a goroutine, with the agentChannel
// getting the agents for each simulation.
//
// Note that because the simulation affects the agents and state directly, we do not have
// to return anything from this routine through another channel!
// Instead we send back only a simple signal through the simulationFinishedSignalChannel.
// This allows for the manager to count how many simulations have passed
func ConcurrentSimulationRoutine(system system.System, agentChannel <-chan []*agent.Agent, simulationFinishedSignalChannel chan<- struct{}) {
	// This loop will take items out of the agentChannel
	//
	// Once the agent channel is closed (done by the manager, once all agents are sent),
	// this loop will exit and the goroutine will terminate
	for agents := range agentChannel {
		SimulateSystem(system, agents)
		simulationFinishedSignalChannel <- struct{}{}
	}
}

// Simulate the given system until the state is found to be terminal
func SimulateSystem(system system.System, agents []*agent.Agent) {
	state := system.InitializeState()

	// Loop forever (until very large value)
	// or until the state is found to be terminal
	for sanityCheck := 0; sanityCheck < MAXIMUM_SIMULATION_ITERATIONS; sanityCheck++ {
		if state.TerminalState {
			break
		}
		system.AdvanceState(state, agents)
	}
}

// Simulate the given system until state is terminal
// Save each state to a file for easy inspection
func SimulateSystemWithSave(system system.System, agents []*agent.Agent, simulationDataCollector *datacollector.SimulationDataCollector) {

	state := system.InitializeState()
	simulationDataCollector.CollectSimulationData(state)

	// Loop forever (until very large value)
	// or until the state is found to be terminal
	for sanityCheck := 0; sanityCheck < MAXIMUM_SIMULATION_ITERATIONS; sanityCheck++ {
		if state.TerminalState {
			break
		}
		system.AdvanceState(state, agents)
		simulationDataCollector.CollectSimulationData(state.DeepCopyState())
	}
}
