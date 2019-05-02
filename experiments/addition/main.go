package main

import (
	"github.com/unixpickle/learn-quantum/quantum"
)

const (
	BitSize  = 5
	NumBits  = BitSize * 2
	MaxCache = 5000000
)

func main() {
	hasher := quantum.NewCircuitHasher(NumBits)
	// TODO: generate basis.
	gen := quantum.NewCircuitGen(4, basis, MaxCache)
	backward := quantum.NewBackwardsMap(hasher, AdderGate{})

	for i := 0; i < 5; i++ {
		// TODO: generate backward circuits.
	}

	// TODO: forward search over sliding circuits.
}
