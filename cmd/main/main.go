package main

import (
	sys "OCSS/FoosballGeneticLearning/cmd/main/multiAgentTest"
	manager "OCSS/FoosballGeneticLearning/pkg/Manager"
)

func main() {
	targetSystem := sys.TestSystem{}
	manager := manager.NewManager(&targetSystem, 100)
	for generationIndex := 0; generationIndex < 1000; generationIndex++ {
		manager.SimulateGeneration()
	}
}
