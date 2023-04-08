package geneticbreeder

import (
	agent "OCSS/FoosballGeneticLearning/pkg/Agent"
	"OCSS/FoosballGeneticLearning/pkg/utils"
	"math"
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

func NewGeneticBreeder(randomSource rand.Source) *GeneticBreeder {
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
	kCrossoverWeights := []float64{0.0, 0.0, 0.0, 0.2, 0.2, 0.2, 0.2, 0.2}
	kCrossoverDistribution := distuv.NewCategorical(kCrossoverWeights, randomSource)

	// Define the mutationSegmentDistribution - which determines how many
	// contiguous genes in the chromosome are updated. Currently implemented is
	// a distribution to update some finite number of genes, from 1 to 5
	mutationSegmentWeights := []float64{0.0, 0.2, 0.2, 0.2, 0.2, 0.2}
	mutationSegmentDistribution := distuv.NewCategorical(mutationSegmentWeights, randomSource)

	return &GeneticBreeder{
		randomGenerator:             rand.New(randomSource),
		numParentsDistribution:      numParentsDistribution,
		kCrossoverDistribution:      kCrossoverDistribution,
		mutationRate:                math.Pow10(-6),
		mutationSegmentDistribution: mutationSegmentDistribution,
	}
}

// Given the current generation of agents, as well as the agent scores,
// calculate the next generation of agents. This is done by, for each new agent
// 1. Finding the parents of the agent (based on fitness score)
// 2. Combining those parents in some way (see combineAgents function)
// 3. Applying any mutations.
func (gb *GeneticBreeder) NextGeneration(currentGeneration []*agent.Agent) []*agent.Agent {
	numAgents := len(currentGeneration)
	newGeneration := make([]*agent.Agent, numAgents)

	generationScores := make([]float64, len(currentGeneration))
	for agentIndex := range currentGeneration {
		generationScores[agentIndex] = currentGeneration[agentIndex].Score
	}

	minimumScore := utils.MinElementInSlice(generationScores)
	if minimumScore < 0 {
		for index := range generationScores {
			generationScores[index] -= minimumScore
		}
	}

	for agentIndex := range newGeneration {
		newGeneration[agentIndex] = gb.breedNewAgent(currentGeneration, generationScores)
	}

	return newGeneration
}

// Creates a new agent given the previous generation
func (gb *GeneticBreeder) breedNewAgent(currentGeneration []*agent.Agent, generationScores []float64) *agent.Agent {
	newAgentParents := gb.selectParents(currentGeneration, generationScores)
	newAgent := gb.combineParents(newAgentParents)
	newAgent = gb.applyMutation(newAgent)
	return newAgent
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
func (gb *GeneticBreeder) selectParents(possibleParents []*agent.Agent, parentScores []float64) []*agent.Agent {
	// Determine how many parents we will select
	numParents := int(gb.numParentsDistribution.Rand())

	// Create a distribution (with random seed) to select parents, based on
	// parent fitness
	parentSelectionSource := gb.randomGenerator.Uint64()
	parentSelectionDistribution := distuv.NewCategorical(parentScores, rand.NewSource(parentSelectionSource))
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
func (gb *GeneticBreeder) combineParents(parents []*agent.Agent) *agent.Agent {
	// Process the agent chromosomes into a useful format, and collect the size information
	parentChromosomes := make([][]float64, len(parents))
	for index, parent := range parents {
		parentChromosomes[index] = parent.Chromosome.RawMatrix().Data
	}
	chromosomeRows, chromosomeCols := parents[0].Chromosome.Dims()
	chromosomeSize := chromosomeRows * chromosomeCols

	// First, we should determine the number of crossovers to take.
	// i.e. we select a value of k.
	numCrossovers := int(gb.kCrossoverDistribution.Rand())

	// Then we select the crossover points by creating an array of all possible crossover indices
	// then shuffling that, taking the first k elements, and sorting the result
	indexArray := make([]int, chromosomeSize)
	for i := 0; i < chromosomeSize; i++ {
		indexArray[i] = i
	}
	utils.ShuffleSlice(gb.randomGenerator, indexArray)
	crossoverPoints := indexArray[:numCrossovers]
	sort.Ints(crossoverPoints)

	// Next we can actually start putting the parent chromosomes together!  Yay!
	// We should start with random parent, to avoid biasing the fittest parent to the start of the chromosome
	childChromosomeData := make([]float64, chromosomeSize)
	childChromosomeIndex := 0
	parentIndex := gb.randomGenerator.Int() % len(parents)
	for _, crossoverPoint := range crossoverPoints {
		parentChromosome := parentChromosomes[parentIndex]
		copy(childChromosomeData[childChromosomeIndex:crossoverPoint], parentChromosome[childChromosomeIndex:crossoverPoint])
		childChromosomeIndex = crossoverPoint
		parentIndex = (parentIndex + 1) % len(parents)
	}

	childChromosome := mat.NewDense(chromosomeRows, chromosomeCols, childChromosomeData)
	return agent.NewAgent(childChromosome)
}

// Apply a random mutation with some probability. This probability should be small, but non-zero.
//
// TODO(hayden): Implement some reasonable method for determining mutations. Experiment should
// reveal a decent value for mutation rate, but we should be slightly clever in how mutations actually occur.
// A single mutation position may not be enough to ensure mutations produce fitter agents every so often,
// and perhaps entire sections of a chromosome must be mutated...
func (gb *GeneticBreeder) applyMutation(agent *agent.Agent) *agent.Agent {
	// If we do not roll a mutation - don't do anything!
	if distuv.UnitUniform.Rand() > gb.mutationRate {
		return agent
	}

	chromosomeData := agent.Chromosome.RawMatrix().Data
	chromosomeRows, chromosomeCols := agent.Chromosome.Dims()
	chromosomeSize := chromosomeRows * chromosomeCols
	// We have a mutation - let's figure out where we are applying this and how much mutation we apply!
	mutationSegmentSize := int(gb.mutationSegmentDistribution.Rand())
	mutationStartIndexDistribution := distuv.Uniform{
		Min: 0,
		Max: float64(chromosomeSize - mutationSegmentSize),
		Src: rand.NewSource(gb.randomGenerator.Uint64()),
	}
	mutationSegmentStartIndex := int(mutationStartIndexDistribution.Rand())

	// Current mutation distribution is simply a gaussian with same mean and scale as the chromosome
	chromosomeMean, chromosomeStd := utils.SummaryStatistics(chromosomeData)
	mutationDistribution := distuv.Normal{
		Mu:    chromosomeMean,
		Sigma: chromosomeStd,
		Src:   rand.NewSource(gb.randomGenerator.Uint64()),
	}

	for i := 0; i < mutationSegmentSize; i++ {
		chromosomeData[mutationSegmentStartIndex+i] = mutationDistribution.Rand()
	}

	// For peace of mind, let's recreate the chromosome matrix
	agent.Chromosome = mat.NewDense(chromosomeRows, chromosomeCols, chromosomeData)
	return agent
}
