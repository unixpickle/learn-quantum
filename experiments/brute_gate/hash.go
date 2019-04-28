package main

import (
	"crypto/md5"

	"github.com/unixpickle/learn-quantum/quantum"
)

func HashCircuit(numBits int, c quantum.Circuit) SimHash {
	data := make([]byte, 0, numBits*4*(1<<uint(numBits)))
	for i := 0; i < (1 << uint(numBits)); i++ {
		sim := quantum.NewSimulationBits(numBits, uint(i))
		c.Apply(sim)
		data = append(data, encodeQuantumState(sim)...)
	}
	return md5.Sum(data)
}

func HashCircuitBackwards(numBits int, c quantum.Circuit, inToOut []int) SimHash {
	inv := c.Inverse()
	data := make([]byte, 0, numBits*4*len(inToOut))
	for _, i := range inToOut {
		sim := quantum.NewSimulationBits(numBits, uint(i))
		inv.Apply(sim)
		data = append(data, encodeQuantumState(sim)...)
	}
	return md5.Sum(data)
}

func HashClassicalGate(numBits int, inToOut []int) SimHash {
	data := make([]byte, 0, numBits*4*len(inToOut))
	for _, i := range inToOut {
		sim := quantum.NewSimulationBits(numBits, uint(i))
		data = append(data, encodeQuantumState(sim)...)
	}
	return md5.Sum(data)
}

func encodeQuantumState(s *quantum.Simulation) []byte {
	data := make([]byte, 0, len(s.Phases)*4)
	for _, phase := range s.Phases {
		r := uint16(int16(30000 * real(phase)))
		i := uint16(int16(30000 * imag(phase)))
		data = append(data, byte(r>>8), byte(r), byte(i>>8), byte(i))
	}
	return data
}
