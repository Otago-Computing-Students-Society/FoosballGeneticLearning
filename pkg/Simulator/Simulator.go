package simulator

import (
	agent "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/Agent"
	datacollector "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/DataCollector"
	system "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/System"
)

const MAXIMUM_SIMULATION_ITERATIONS = 1000000

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
