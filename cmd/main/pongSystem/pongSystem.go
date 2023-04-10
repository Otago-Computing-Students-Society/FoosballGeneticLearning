package pongsystem

import (
	"time"

	agent "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/Agent"
	systemstate "github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/SystemState"
	"github.com/Otago-Computer-Science-Society/FoosballGeneticLearning/pkg/utils"
	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/mat"
)

// This system implements a simple game of Pong
//
// The system will advance until one agent scores a point,
// at which time it will be made terminal
//
// The bounds of the system will be [-1,1] is both the X direction
// (between the agents) and [-0.5,0.5] Y direction (across the agents)
//
// Y
// ^
// |  ---------------------------------------------------
// |  |													|
// |  |												|	|
// |  |												|	|
// |  |							o					|	|
// |  |	|												|
// |  |	|												|
// |  |	|												|
// |  |													|
// |  ---------------------------------------------------
// |
// --------------------------------------------------> X
//
//
// The ball moves with a constant velocity, and bounces specularly  off walls.
// For simplicity we may start with the ball bouncing specularly off the paddles, but we could change this later

const (
	// State vector consists of:
	// - ballX
	// - ballY
	// - ballXVelocity
	// - ballYVelocity
	// - paddle0Position
	// - paddle1Position
	STATE_VECTOR_LEN = 6

	// Percepts are given in the following order:
	// - ballX
	// - ballY
	// - ballXVelocity
	// - ballYVelocity
	// - paddle1Position
	// - paddle2Position
	NUM_PERCEPTS = 6

	// Agent can take exactly 1 action:
	// - paddleVelocity
	NUM_ACTIONS = 1

	// Pong is a two player game
	NUM_AGENTS_PER_SIMULATION = 2

	// The X dimension of the game.
	// This extends both positive and negative about 0.0
	GAME_X_DIMENSION = 1.0

	// The Y dimension of the game.
	// This extends both positive and negative about 0.0
	GAME_Y_DIMENSION = 0.5

	// Sets how large the agent paddles are
	// Note the paddle extends half this amount above and below
	// the paddle position. Remember the total Y distance is of size 1
	PADDLE_SIZE = 0.2

	// The velocity cap of the paddle
	MAX_PADDLE_VELOCITY = 1.0

	// The time delta of the system
	// Defines how far to step the system physics each advancement
	TIME_DELTA = 0.01

	// How much score is given for scoring
	SCORING_SCORE = 100.0

	// How much score is given for bouncing
	BOUNCE_SCORE = 0.0

	// Score given for being "in front of" ball
	READY_SCORE = 1.0
)

// Given a velocity vector that is bouncing off a surface (with given normal)
// update the velocity vector to the velocity after bouncing
func bounceSpecularly(velocityVector *mat.VecDense, reflectionNormal *mat.VecDense) {
	// Actually ensure the reflectionNormal is unit
	reflectionNormal.ScaleVec(1/reflectionNormal.Norm(2), reflectionNormal)

	// Get the unit direction vector of the velocity
	velocityDirection := mat.VecDenseCopyOf(velocityVector)
	velocityDirection.ScaleVec(1/velocityVector.Norm(2), velocityVector)

	// Find the new direction of the velocity after bouncing
	// See https://en.wikipedia.org/wiki/Specular_reflection#Vector_formulation
	identityMatrix := mat.NewDiagDense(2, []float64{1, 1})
	householderTransformationMatrix := mat.NewDense(2, 2, nil)
	householderTransformationMatrix.RankOne(identityMatrix, -2, reflectionNormal, reflectionNormal)
	velocityDirection.MulVec(householderTransformationMatrix, velocityDirection)

	// Rescale to correct velocity after bounce
	// Here is where we could also apply non-elastic collisions...
	velocityVector.ScaleVec(velocityVector.Norm(2), velocityDirection)
}

type PongSystem struct {
	randomGenerator *rand.Rand
}

func NewPongSystem() *PongSystem {
	randomGenerator := rand.New(rand.NewSource(uint64(time.Now().Nanosecond())))

	return &PongSystem{
		randomGenerator: randomGenerator,
	}
}

func (system *PongSystem) NumPercepts() int {
	return NUM_PERCEPTS
}
func (system *PongSystem) NumActions() int {
	return NUM_ACTIONS
}
func (system *PongSystem) NumAgentsPerSimulation() int {
	return NUM_AGENTS_PER_SIMULATION
}

// Returns the initial state of the system
func (system *PongSystem) InitializeState() *systemstate.SystemState {
	// Ball always starts exactly halfway between agents
	ballX := 0.0
	// Ball starts at random Y position
	ballY := (2*system.randomGenerator.Float64() - 1) * GAME_Y_DIMENSION
	// Ball has some random initial velocity in both Y direction [-1.5,-0.5] U [0.5,1.5]
	ballYVelocity := (0.5 * (system.randomGenerator.Float64() + 1)) * GAME_Y_DIMENSION
	// Ball has a larger velocity in the X direction, and is randomly set to
	// either positive or negative X direction [-1.5,-0.5] U [0.5,1.5]
	ballXVelocity := (0.5 * (system.randomGenerator.Float64() + 1)) * GAME_X_DIMENSION
	if system.randomGenerator.NormFloat64() < 0 {
		ballXVelocity *= -1
	}
	// Paddles start in neutral position
	paddle0Position := 0.0
	paddle1Position := 0.0
	return &systemstate.SystemState{
		StateVector: mat.NewVecDense(STATE_VECTOR_LEN, []float64{
			ballX,
			ballY,
			ballXVelocity,
			ballYVelocity,
			paddle0Position,
			paddle1Position,
		}),
		TerminalState: false,
	}
}

// Defines the behavior of the system
//
// This method should take the current system state and agents, find the agent actions,
// then apply those actions to the state to update it.
//
// This method is also responsible for updating the agent scores!
//
// The implementation here is very boring - we immediately call the simulation done (terminal = true)
// and set the score of the agent to the sum of the action vector.
func (system *PongSystem) AdvanceState(state *systemstate.SystemState, agents []*agent.Agent) {
	// Get the data out of the state
	ballX := state.StateVector.AtVec(0)
	ballY := state.StateVector.AtVec(1)
	ballXVelocity := state.StateVector.AtVec(2)
	ballYVelocity := state.StateVector.AtVec(3)
	paddle0Position := state.StateVector.AtVec(4)
	paddle1Position := state.StateVector.AtVec(5)

	// Create the percept vectors
	perceptVector := mat.NewVecDense(NUM_PERCEPTS, []float64{
		ballX,
		ballY,
		ballXVelocity,
		ballYVelocity,
		paddle0Position,
		paddle1Position,
	})

	// Because the system is entirely symmetric on the X axis we can pass
	// The inverse state to agent1 and ensure that all agents
	// are always playing on the "left"
	inversePerceptVector := mat.NewVecDense(NUM_PERCEPTS, []float64{
		-1.0 * ballX,
		ballY,
		-1.0 * ballXVelocity,
		ballYVelocity,
		paddle0Position,
		paddle1Position,
	})

	// Get agent actions
	agent0Action := agents[0].GetAction(perceptVector)
	agent1Action := agents[1].GetAction(inversePerceptVector)

	paddle0Velocity := agent0Action.AtVec(0)
	paddle1Velocity := agent1Action.AtVec(0)

	// Check if ball is "scored"
	if ballX >= 1.0 {
		agents[0].Score += SCORING_SCORE
		state.TerminalState = true
		return
	}

	if ballX <= -1.0 {
		agents[1].Score += SCORING_SCORE
		state.TerminalState = true
		return
	}

	// Update the objects in the system
	ballPosition := mat.NewVecDense(2, []float64{ballX, ballY})
	ballVelocity := mat.NewVecDense(2, []float64{ballXVelocity, ballYVelocity})
	ballPosition.AddScaledVec(ballPosition, TIME_DELTA, ballVelocity)
	ballX = ballPosition.AtVec(0)
	ballY = ballPosition.AtVec(1)

	paddle0Velocity = utils.ClipToBounds(paddle0Velocity, -MAX_PADDLE_VELOCITY, MAX_PADDLE_VELOCITY)
	paddle1Velocity = utils.ClipToBounds(paddle1Velocity, -MAX_PADDLE_VELOCITY, MAX_PADDLE_VELOCITY)
	paddle0Position = utils.ClipToBounds(paddle0Position+TIME_DELTA*paddle0Velocity, -GAME_Y_DIMENSION, GAME_Y_DIMENSION)
	paddle1Position = utils.ClipToBounds(paddle1Position+TIME_DELTA*paddle1Velocity, -GAME_Y_DIMENSION, GAME_Y_DIMENSION)

	// Reflect ball off bottom wall
	if ballY <= -0.5 {
		bounceSpecularly(ballVelocity, mat.NewVecDense(2, []float64{0.0, 1.0}))
	}
	// Reflect ball off top wall
	if ballY >= 0.5 {
		bounceSpecularly(ballVelocity, mat.NewVecDense(2, []float64{0.0, -1.0}))
	}
	// Reflect ball off left paddle if and only if paddle0 is in the way
	if ballX <= -1.0 && (paddle0Position-PADDLE_SIZE < ballY && ballY < paddle0Position+PADDLE_SIZE) {
		bounceSpecularly(ballVelocity, mat.NewVecDense(2, []float64{1.0, 0.0}))
		agents[0].Score += BOUNCE_SCORE
		ballX = -0.9
	}
	// Reflect ball off right paddle if and only if paddle1 is in the way
	if ballX >= 1.0 && (paddle1Position-PADDLE_SIZE < ballY && ballY < paddle1Position+PADDLE_SIZE) {
		bounceSpecularly(ballVelocity, mat.NewVecDense(2, []float64{-1.0, 0.0}))
		agents[1].Score += BOUNCE_SCORE
		ballX = 0.9
	}

	ballXVelocity = ballVelocity.AtVec(0)
	ballYVelocity = ballVelocity.AtVec(1)

	// Check if agent has paddle in front of ball
	if paddle0Position-PADDLE_SIZE < ballY && ballY < paddle0Position+PADDLE_SIZE {
		agents[0].Score += READY_SCORE
	}
	if paddle1Position-PADDLE_SIZE < ballY && ballY < paddle1Position+PADDLE_SIZE {
		agents[1].Score += READY_SCORE
	}

	ballX = utils.ClipToBounds(ballX, -GAME_X_DIMENSION, GAME_X_DIMENSION)
	ballY = utils.ClipToBounds(ballY, -GAME_Y_DIMENSION, GAME_Y_DIMENSION)

	// Update the state with new data
	state.StateIndex += 1
	state.StateVector.SetVec(0, ballX)
	state.StateVector.SetVec(1, ballY)
	state.StateVector.SetVec(2, ballXVelocity)
	state.StateVector.SetVec(3, ballYVelocity)
	state.StateVector.SetVec(4, paddle0Position)
	state.StateVector.SetVec(5, paddle1Position)

	// fmt.Printf("%v %v\n", state.TerminalState, mat.Formatted(state.StateVector.T(), mat.Squeeze()))
}
