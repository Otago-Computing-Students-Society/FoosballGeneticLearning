package geneticbreeder

import (
	agent "OCSS/FoosballGeneticLearning/pkg/Agent"
	"OCSS/FoosballGeneticLearning/pkg/utils"
	"sort"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

type GeneticBreeder struct {
	randomGenerator             *rand.Rand
	numParentsDistribution      distuv.Rander
	kCrossoverDistribution      distuv.Rander
	mutationRate                float64
	mutationSegmentDistribution distuv.Rander
}

func NewGeneticBreeder(randomSource rand.Source, numberParents int) *GeneticBreeder {
	// Define the numParents distribution
	// // Current implementation has number of parents selected as
	// 0.5 chance of 2 parents, 0.5 change of 3 parents.
	//
	// See https://pkg.go.dev/gonum.org/v1/gonum@v0.12.0/stat/distuv#Categorical for explanation
	numParentsWeights := []float64{0.0, 0.0, 0.5, 0.5}
	numParentsDistribution := distuv.NewCategorical(numParentsWeights, randomSource)

	// Define the crossover distribution. This should be proportional to the size of
	// the chromosome, but we can also have a set value (since the proportionality shouldn't be huge)
	// In this implementation, we have a set probability of some small number for k.
	// See https://pkg.go.dev/gonum.org/v1/gonum@v0.12.0/stat/distuv#Categorical for explanation
	kWeights := []float64{0.0, 0.0, 0.0, 0.2, 0.2, 0.2, 0.2, 0.2}
	kCrossoverDistribution := distuv.NewCategorical(kWeights, randomSource)

	// Define the mutationSegmentDistribution - which determines how many
	// contiguous genes in the chromosome are updated. Currently implemented is
	// a distribution to update some finite number of genes, from 1 to 5
	mutationSegmentWeights := []float64{0.0, 0.2, 0.2, 0.2, 0.2, 0.2}
	mutationSegmentDistribution := distuv.NewCategorical(mutationSegmentWeights, randomSource)

	return &GeneticBreeder{
		randomGenerator:             rand.New(randomSource),
		numParentsDistribution:      numParentsDistribution,
		kCrossoverDistribution:      kCrossoverDistribution,
		mutationRate:                0.01,
		mutationSegmentDistribution: mutationSegmentDistribution,
	}
}

// Given the current generation of agents, as well as the agent scores,
// calculate the next generation of agents. This is done by, for each new agent
// 1. Finding the parents of the agent (based on fitness score)
// 2. Combining those parents in some way (see combineAgents function)
// 3. Applying any mutations.
func NextGeneration(currentGeneration []*agent.Agent, generationScores []float64, randomSource *rand.Source) []*agent.Agent {
	numAgents := len(currentGeneration)
	newGeneration := make([]*agent.Agent, numAgents)

	for agentIndex := range newGeneration {
		newAgentParents := selectParents(currentGeneration, generationScores, nil)
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
func selectParents(possibleParents []*agent.Agent, parentScores []float64, randomSource rand.Source) []*agent.Agent {
	if randomSource == nil {
		randomSource = rand.NewSource(uint64(time.Now().UnixNano()))
	}
	// Generate a probability distribution of parents based on the parentScores
	parentSelectionDistribution := distuv.NewCategorical(parentScores, randomSource)

	// Also create a quick distribution to select the number of parents
	triangleDistLow := 2.0
	triangleDistMode := 3.5
	triangleDistHigh := 5.0
	if len(possibleParents) < int(triangleDistHigh) {
		fmt.Fprintf(os.Stderr, "ERROR: Not enough parents to allow for %v as parent upper distribution\n", triangleDistHigh)
		os.Exit(1)
	}

	// Yes... the NewTriangle signature is indeed low, high, mode... unintuitive...
	numParentsDistribution := distuv.NewTriangle(triangleDistLow, triangleDistHigh, triangleDistMode, randomSource)
	numParents := int(numParentsDistribution.Rand())

	// Now we can select a number of selectedParents based on these probabilities
	selectedParents := []*agent.Agent{}
	selectedParentIndices := []int{}

	// For a number of times equal to the number of parents requested
	for i := 0; i < numParents; i++ {
		// Loop forever, trying a new random parent index (weighted by parent fitness)
		// Actually only loop for a very large number, to avoid infinite loops by mistake.
		for sanityValue := 0; sanityValue < 1000*len(possibleParents); sanityValue++ {
			selectedParentIndex := int(parentSelectionDistribution.Rand())
			// If we have seen this parent before, we try again...
			if !utils.IsElementInSlice(selectedParentIndices, selectedParentIndex) {
				selectedParentIndices = append(selectedParentIndices, selectedParentIndex)
				break
			}
		}
	}

	// Translate the selected parent indices into parents, and return
	for _, selectedParentIndex := range selectedParentIndices {
		selectedParents = append(selectedParents, possibleParents[selectedParentIndex])
	}
	return selectedParents
}

// Given a set of parents, combine each of their chromosomes in some sensible fashion
// to create a new chromosome (and hence a new agent).
//
// This function should accept an arbitrary number of agents as parents, rather than
// (for example) exactly two parents
//
// TODO(hayden): Come up with a good implementation for this. I suggest k-point crossover.
func combineParents(parents []*agent.Agent, parentScores []float64) *agent.Agent {
	return nil
}

// Apply a random mutation with some probability. This probability should be small, but non-zero.
//
// TODO(hayden): Implement some reasonable method for determining mutations. Experiment should
// reveal a decent value for mutation rate, but we should be slightly clever in how mutations actually occur.
// A single mutation position may not be enough to ensure mutations produce fitter agents every so often,
// and perhaps entire sections of a chromosome must be mutated...
func applyMutation(agent *agent.Agent) *agent.Agent {
	return nil
}
