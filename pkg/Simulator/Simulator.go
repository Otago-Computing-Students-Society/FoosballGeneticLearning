package simulator

import (
	agent "OCSS/FoosballGeneticLearning/pkg/Agent"
	system "OCSS/FoosballGeneticLearning/pkg/System"
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
