package simulator

import (
	agent "OCSS/FoosballGeneticLearning/pkg/Agent"
	system "OCSS/FoosballGeneticLearning/pkg/System"
)

func SimulateSystem(system system.System, agents []*agent.Agent) {
	state := system.InitializeState()
	agentActions := make([]agent.AgentAction, len(agents))

	for !state.TerminalState {
		for agentIndex := range agents {
			agentActions[agentIndex] = agents[agentIndex].GetAction(state.StateVector)
		}

		state = system.AdvanceState(state, agentActions)
	}
}
