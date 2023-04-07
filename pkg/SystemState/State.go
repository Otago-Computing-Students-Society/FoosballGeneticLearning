package systemstate

import "gonum.org/v1/gonum/mat"

type SystemState struct {
	stateVector   *mat.Dense
	terminalState bool
}

type StateHistory []*SystemState
