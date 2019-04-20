package main

import (
	"math/cmplx"
	"math/rand"

	"github.com/unixpickle/learn-quantum/quantum"
)

func Search(numBits, maxGates int, gate func(b []bool) []bool, results chan<- quantum.Circuit) {
OuterLoop:
	for {
		c := RandomCircuit(numBits, rand.Intn(maxGates)+1)
		for inIdx := 0; inIdx < 1<<uint(numBits); inIdx++ {
			input := make([]bool, numBits)
			for i := range input {
				input[i] = (inIdx&(1<<uint(i)) != 0)
			}
			sim := quantum.NewSimulation(numBits)
			for i, b := range input {
				if b {
					quantum.X(sim, i)
				}
			}
			c.Apply(sim)
			output := gate(input)
			if cmplx.Abs(sim.Phase(output)-1) > 1e-8 {
				continue OuterLoop
			}
		}
		results <- c
	}
}
