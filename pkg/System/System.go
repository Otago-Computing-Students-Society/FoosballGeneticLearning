package system

import (
	agent "OCSS/FoosballGeneticLearning/pkg/Agent"
	systemstate "OCSS/FoosballGeneticLearning/pkg/SystemState"
)

type System interface {
	NumPercepts() int
	NumActions() int
	NumAgentsPerSimulation() int
	InitializeState() *systemstate.SystemState
	AdvanceState(systemstate.SystemState, []agent.AgentAction) *systemstate.SystemState
}
