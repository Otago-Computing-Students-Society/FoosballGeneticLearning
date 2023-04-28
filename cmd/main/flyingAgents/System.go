package flyingagents

// Agents are a body with a thruster on each side. Agents must navigate to target locations to score points.
//
// Agents will receive information on the next target position, agent position, orientation, velocity.
//
// Agents will act to produce a thrust on each side of the body.

import (
	"math"
	"time"

	agent "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/Agent"
	systemstate "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/SystemState"
	"github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/utils"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
)

const (
	// A state vector consists of:
	// 0 - AgentX
	// 1 - AgentY
	// 2 - AgentVelX
	// 3 - AgentVelY
	// 4 - TargetLocationX
	// 5 - TargetLocationY
	// 6 - NumTargetLocationsVisited
	STATE_VECTOR_LEN = 7

	// An agent will receive:
	// 0 - AgentX
	// 1 - AgentY
	// 2 - AgentVelX
	// 3 - AgentVelY
	// 4 - TargetLocationX
	// 5 - TargetLocationY
	NUM_PERCEPTS = 6

	// The agent must act on:
	// 0 - UpThruster
	// 1 - DownThruster
	// 2 - LeftThruster
	// 3 - RightThruster
	NUM_ACTIONS = 4

	// There is only a single agent per simulation
	NUM_AGENTS_PER_SIMULATION = 1

	// --------------------------------------------------------------------------------------------

	// Defines the bounds of the simulation
	//
	// If an agent reaches this bound (in ANY direction, positive and negative X and Y) then the agent "loses"
	SIMULATION_BOUND = 100.0

	// Defines how large each step is
	TIME_DELTA = 0.5

	// Defines how strong gravity is - each frame, the agents Y velocity is decremented by some amount to simulate this.
	GRAVITY = 0.9

	// Defines how powerful the agents thrust is allowed to be.
	MAX_THRUST = 1.0

	// Defines the radius in which the agent can pick up the reward.
	AGENT_RADIUS = 1.0

	// Defines how many locations the agent is allowed to visit in a single trial
	MAX_LOCATIONS = 10

	// --------------------------------------------------------------------------------------------

	// Defines the penalty for "losing" (running into the simulation bound)
	LOSING_PENALTY = -1000

	// Defines how much reaching a target location is worth.
	TARGET_LOCATION_REWARD = 0

	// Defines how much reward is attributed to moving towards the target location
	MOVEMENT_TOWARDS_LOCATION_REWARD = 1
)

// ------------------------------------------------------------------------------------------------

type FlyingAgentSystem struct {
	randomGenerator *rand.Rand
}

func NewFlyingAgentSystem() *FlyingAgentSystem {
	randomGenerator := rand.New(rand.NewSource(uint64(time.Now().Nanosecond())))

	return &FlyingAgentSystem{
		randomGenerator: randomGenerator,
	}
}

func (system *FlyingAgentSystem) NumPercepts() int {
	return NUM_PERCEPTS
}
func (system *FlyingAgentSystem) NumActions() int {
	return NUM_ACTIONS
}
func (system *FlyingAgentSystem) NumAgentsPerSimulation() int {
	return NUM_AGENTS_PER_SIMULATION
}

// ------------------------------------------------------------------------------------------------

// Determine the next location of the target location.
//
// In future, this function may take the previous target location, to choose a new location sufficiently far away.
func (system *FlyingAgentSystem) chooseNextTargetLocation() (float64, float64) {
	targetLocationX := 2*SIMULATION_BOUND*system.randomGenerator.Float64() - SIMULATION_BOUND
	targetLocationY := 2*SIMULATION_BOUND*system.randomGenerator.Float64() - SIMULATION_BOUND
	return targetLocationX, targetLocationY
}

// Give the initial state of a system
//
// We must position the agent, as well as determine the first target location
func (system *FlyingAgentSystem) InitializeState() *systemstate.SystemState {
	// agent always starts in the middle of the simulation with 0 orientation and 0 velocity
	agentX := 0.0
	agentY := 0.0
	agentVelX := 0.0
	agentVelY := 0.0

	// First target location is given by a random number chosen in the valid bounds
	targetLocationX, targetLocationY := system.chooseNextTargetLocation()
	numTargetLocationsVisited := 0.0

	return &systemstate.SystemState{
		StateVector: mat.NewVecDense(STATE_VECTOR_LEN, []float64{
			agentX,
			agentY,
			agentVelX,
			agentVelY,
			targetLocationX,
			targetLocationY,
			numTargetLocationsVisited,
		}),
		TerminalState: false,
	}
}

func (system *FlyingAgentSystem) AdvanceState(state *systemstate.SystemState, agents []*agent.Agent) {
	agentX := state.StateVector.AtVec(0)
	agentY := state.StateVector.AtVec(1)
	agentVelX := state.StateVector.AtVec(2)
	agentVelY := state.StateVector.AtVec(3)
	targetLocationX := state.StateVector.AtVec(4)
	targetLocationY := state.StateVector.AtVec(5)
	numTargetLocationsVisited := state.StateVector.AtVec(6)

	// Craft the agent input vector
	agentPercepts := mat.NewVecDense(NUM_PERCEPTS, []float64{
		agentX,
		agentY,
		agentVelX,
		agentVelY,
		targetLocationX,
		targetLocationY,
	})
	agent := agents[0]
	agentAction := agent.GetAction(agentPercepts)

	// Decode the agent action
	upThruster := agentAction.AtVec(0)
	downThruster := agentAction.AtVec(1)
	leftThruster := agentAction.AtVec(2)
	rightThruster := agentAction.AtVec(3)
	// Ensure thrusters are in correct bound
	upThruster = utils.ClipToBounds(upThruster, 0, MAX_THRUST)
	downThruster = utils.ClipToBounds(downThruster, 0, MAX_THRUST)
	leftThruster = utils.ClipToBounds(leftThruster, 0, MAX_THRUST)
	rightThruster = utils.ClipToBounds(rightThruster, 0, MAX_THRUST)

	// Update the agents position and velocities --------------------------------------------------

	newAgentX := agentX + TIME_DELTA*agentVelX
	newAgentY := agentY + TIME_DELTA*agentVelY

	newAgentVelX := agentVelX + TIME_DELTA*(rightThruster-leftThruster)
	newAgentVelY := agentVelY + TIME_DELTA*(upThruster-downThruster-GRAVITY)

	// Check over agent rewards -------------------------------------------------------------------

	// Check if agent have left simulation
	if newAgentX < -SIMULATION_BOUND || newAgentX > SIMULATION_BOUND || newAgentY < -SIMULATION_BOUND || newAgentY > SIMULATION_BOUND {
		agent.Score += LOSING_PENALTY
		state.TerminalState = true
	}

	// Check if agent has moved towards the reward (in both directions)
	if math.Abs(targetLocationX-newAgentX) < math.Abs(targetLocationX-agentX) && math.Abs(targetLocationY-newAgentY) < math.Abs(targetLocationY-agentY) {
		agent.Score += MOVEMENT_TOWARDS_LOCATION_REWARD
	}

	// Check if agent has collected the reward
	if math.Pow(targetLocationX-newAgentX, 2)+math.Pow(targetLocationY-newAgentY, 2) < math.Pow(AGENT_RADIUS, 2) {
		agent.Score += TARGET_LOCATION_REWARD

		targetLocationX, targetLocationY = system.chooseNextTargetLocation()
		numTargetLocationsVisited += 1
		if numTargetLocationsVisited >= MAX_LOCATIONS {
			state.TerminalState = true
		}
	}

	// Update state and finish up
	state.StateIndex += 1
	state.StateVector.SetVec(0, newAgentX)
	state.StateVector.SetVec(1, newAgentY)
	state.StateVector.SetVec(2, newAgentVelX)
	state.StateVector.SetVec(3, newAgentVelY)
	state.StateVector.SetVec(4, targetLocationX)
	state.StateVector.SetVec(5, targetLocationY)
	state.StateVector.SetVec(6, numTargetLocationsVisited)
}
