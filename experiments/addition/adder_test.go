package main

import (
	"math/cmplx"
	"math/rand"
	"testing"

	"github.com/unixpickle/learn-quantum/quantum"
)

func TestAdder(t *testing.T) {
	// A: 10001
	// B: 10101
	// Sum: 00110
	state := 803    // 1100100011
	endState := 297 // 0100101001

	qc := quantum.NewSimulationBits(10, uint(state))
	AddGate{}.Apply(qc)
	if cmplx.Abs(qc.Phases[endState]-1) > 1e-8 {
		t.Error("incorrect end state:", qc.Sample())
	}
}

func TestSubtractor(t *testing.T) {
	for i := 0; i < 100; i++ {
		state := rand.Intn(1 << 10)
		qc := quantum.NewSimulationBits(10, uint(state))
		AddGate{}.Apply(qc)
		SubGate{}.Apply(qc)
		if cmplx.Abs(qc.Phases[state]-1) > 1e-8 {
			t.Error("incorrect end state")
		}
	}
}
