package main

import (
	"fmt"

	"github.com/unixpickle/learn-quantum/quantum"
)

const (
	BitSize  = 5
	NumBits  = BitSize * 2
	MaxCache = 5000000
)

func main() {
	hasher := quantum.NewCircuitHasher(NumBits)
	gen := quantum.NewCircuitGen(4, generateBasis(), MaxCache)
	backward := quantum.NewBackwardsMap(hasher, AddGate{})

	for i := 1; i < 5; i++ {
		ch, count := gen.Generate(i)
		fmt.Println("Doing backward search of depth", i, "with", count, "permutations...")
		for c := range ch {
			c1 := append(quantum.Circuit{}, c...)
			for j, g := range c {
				c1[j] = &EndGate{Gate: g}
			}
			backward.AddCircuit(c1)
		}
	}

	for i := 1; i < 100; i++ {
		ch, count := gen.Generate(i)
		fmt.Println("Doing forward search of depth", i, "with", count, "permutations...")
		for c := range ch {
			forward := &SlidingGate{Gate: c}
			if tail := backward.Lookup(forward); tail != nil {
				fmt.Println(append(c, tail))
				return
			}
		}
	}
}

func generateBasis() []quantum.Gate {
	const numBits = 4
	var result []quantum.Gate
	for i := 0; i < numBits; i++ {
		result = append(result, &quantum.HGate{Bit: i})
		result = append(result, &quantum.TGate{Bit: i})
		result = append(result, &quantum.TGate{Bit: i, Conjugate: true})
		result = append(result, &quantum.XGate{Bit: i})
		result = append(result, &quantum.YGate{Bit: i})
		result = append(result, &quantum.ZGate{Bit: i})
		for j := 0; j < numBits; j++ {
			if j != i {
				result = append(result, &quantum.CNotGate{Control: i, Target: j})
				result = append(result, &quantum.CSqrtNotGate{Control: i, Target: j})
				result = append(result, &quantum.CSqrtNotGate{Control: i, Target: j, Invert: true})
				result = append(result, &quantum.CHGate{Control: i, Target: j})
				if i < j {
					for k := 0; k < numBits; k++ {
						if k != i && k != j {
							result = append(result, &quantum.CCNotGate{
								Control1: i,
								Control2: j,
								Target:   k,
							})
						}
					}
				}
			}
		}
	}
	return result
}
