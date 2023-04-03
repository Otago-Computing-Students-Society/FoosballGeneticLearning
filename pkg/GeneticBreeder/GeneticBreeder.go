package geneticbreeder

import agent "OCSS/FoosballGeneticLearning/pkg/Agent"

// Given the current generation of agents, as well as the agent scores,
// calculate the next generation of agents. This is done by, for each new agent
// 1. Finding the parents of the agent (based on fitness score)
// 2. Combining those parents in some way (see combineAgents function)
// 3. Applying any mutations.
func NextGeneration(currentGeneration []*agent.Agent, generationScores []float64) []*agent.Agent {
	numAgents := len(currentGeneration)
	newGeneration := make([]*agent.Agent, numAgents)

	for agentIndex := range newGeneration {
		newAgentParents := selectParents(currentGeneration, generationScores)
		newAgent := combineParents(newAgentParents, generationScores)
		newAgent = applyMutation(newAgent)
		newGeneration[agentIndex] = newAgent
	}

	return newGeneration
}

// Given the possible parents (a set of agents) and those parents fitness functions,
// select a set of parents for a new agent.
//
// Note this function will return some number of agents as parents. It would be best
// for combineAgents to accept any number of agents to allow for this method
// to return an arbitrary number of agents
//
// TODO(hayden): Come up with a good solution for parent selection. I suggest using
// parent scores as a weighting into a probability distribution. This would allow for
// bad agents to still have a chance of passing on chromosomes, increasing genetic diversity.
func selectParents(possibleParents []*agent.Agent, parentScores []float64) []*agent.Agent {

}

// Given a set of parents, combine each of their chromosomes in some sensible fashion
// to create a new chromosome (and hence a new agent).
//
// This function should accept an arbitrary number of agents as parents, rather than
// (for example) exactly two parents
//
// TODO(hayden): Come up with a good implementation for this. I suggest k-point crossover.
func combineParents(parents []*agent.Agent, parentScores []float64) *agent.Agent {

}

// Apply a random mutation with some probability. This probability should be small, but non-zero.
//
// TODO(hayden): Implement some reasonable method for determining mutations. Experiment should
// reveal a decent value for mutation rate, but we should be slightly clever in how mutations actually occur.
// A single mutation position may not be enough to ensure mutations produce fitter agents every so often,
// and perhaps entire sections of a chromosome must be mutated...
func applyMutation(agent *agent.Agent) *agent.Agent {

}
