package system

import (
	agent "OCSS/FoosballGeneticLearning/pkg/Agent"
	systemstate "OCSS/FoosballGeneticLearning/pkg/SystemState"
)

type System interface {
	NumPercepts() int
	NumActions() int
	ScoreFunction(systemstate.StateHistory) float64
	InitializeState() systemstate.SystemState
	AdvanceState(systemstate.SystemState, []agent.AgentAction) systemstate.SystemState
}
