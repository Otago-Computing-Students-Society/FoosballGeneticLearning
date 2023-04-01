package systemstate

import "gonum.org/v1/gonum/mat"

type SystemState *mat.Dense

type StateHistory []SystemState
