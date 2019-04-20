package main

import (
	"math/rand"

	"github.com/unixpickle/learn-quantum/quantum"
)

func RandomGate(numBits int) quantum.Gate {
	gid := rand.Intn(5)
	if gid == 0 {
		return &quantum.XGate{Bit: rand.Intn(numBits)}
	} else if gid == 1 {
		return &quantum.YGate{Bit: rand.Intn(numBits)}
	} else if gid == 2 {
		return &quantum.ZGate{Bit: rand.Intn(numBits)}
	} else if gid == 3 {
		source := rand.Intn(numBits)
		target := rand.Intn(numBits - 1)
		if target >= source {
			target += 1
		}
		return &quantum.CNotGate{Control: source, Target: target}
	} else if gid == 4 {
		return &quantum.HGate{Bit: rand.Intn(numBits)}
	}
	panic("unreachable")
}

func RandomCircuit(numBits, numGates int) quantum.Circuit {
	var c quantum.Circuit
	for i := 0; i < numGates; i++ {
		c = append(c, RandomGate(numBits))
	}
	return c
}
