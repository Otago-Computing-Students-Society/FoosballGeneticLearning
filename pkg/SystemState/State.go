package systemstate

import "gonum.org/v1/gonum/mat"

type SystemState struct {
	StateVector   *mat.VecDense
	StateIndex    int
	TerminalState bool
}
