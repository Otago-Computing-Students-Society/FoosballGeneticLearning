package main

import (
	sys "OCSS/FoosballGeneticLearning/cmd/main/multiAgentSystem"
	manager "OCSS/FoosballGeneticLearning/pkg/Manager"
)

func main() {
	targetSystem := sys.MultiAgentSystem{}
	manager := manager.NewManager(&targetSystem, 100, true)
	for generationIndex := 0; generationIndex < 1000; generationIndex++ {
		manager.SimulateGeneration()
	}
}
