package systemstate

import "gonum.org/v1/gonum/mat"

type SystemState struct {
	StateVector   *mat.VecDense
	StateIndex    int
	TerminalState bool
}

func (state *SystemState) DeepCopyState() *SystemState {
	return &SystemState{
		StateIndex:  state.StateIndex,
		StateVector: mat.VecDenseCopyOf(state.StateVector),
	}
}
