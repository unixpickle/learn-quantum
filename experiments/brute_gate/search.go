package main

import (
	"crypto/md5"
	"math/rand"

	"github.com/unixpickle/learn-quantum/quantum"
)

func Search(numBits, maxGates int, gate func(b []bool) []bool, results chan<- quantum.Circuit) {
	inToOut := make([]int, 1<<uint(numBits))
	for i := 0; i < 1<<uint(numBits); i++ {
		input := make([]bool, numBits)
		for j := range input {
			input[j] = (i&(1<<uint(j)) != 0)
		}
		output := gate(input)
		outInt := 0
		for j, b := range output {
			if b {
				outInt |= 1 << uint(j)
			}
		}
		inToOut[i] = outInt
	}

	forwards := map[string]quantum.Circuit{}
	backwards := map[string]quantum.Circuit{}
	for {
		c := RandomCircuit(numBits, rand.Intn(maxGates/2)+1)
		var fwdHash string
		for i := 0; i < (1 << uint(numBits)); i++ {
			sim := simulationFromBits(numBits, i)
			c.Apply(sim)
			fwdHash += "  " + sim.String()
		}
		var bwdHash string
		for _, out := range inToOut {
			sim := simulationFromBits(numBits, out)
			c.Invert(sim)
			bwdHash += "  " + sim.String()
		}

		fwdHash = hashStr(fwdHash)
		bwdHash = hashStr(bwdHash)

		if c1, ok := forwards[fwdHash]; !ok || len(c1) > len(c) {
			forwards[fwdHash] = c
		}
		if c1, ok := backwards[bwdHash]; !ok || len(c1) > len(c) {
			backwards[bwdHash] = c
		}
		if c1, ok := backwards[fwdHash]; ok {
			results <- append(append(quantum.Circuit{}, c...), c1...)
		} else if c1, ok = forwards[bwdHash]; ok {
			results <- append(append(quantum.Circuit{}, c1...), c...)
		}
	}
}

func simulationFromBits(numBits, bits int) *quantum.Simulation {
	res := quantum.NewSimulation(numBits)
	for i := 0; i < numBits; i++ {
		if bits&(1<<uint(i)) != 0 {
			quantum.X(res, i)
		}
	}
	return res
}

func hashStr(s string) string {
	sum := md5.Sum([]byte(s))
	return string(sum[:])
}
