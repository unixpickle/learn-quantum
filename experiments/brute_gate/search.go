package main

import (
	"fmt"

	"github.com/unixpickle/learn-quantum/quantum"
)

const maxCircuitCache = 5000000

func AllGates(numBits int, includeCCNot bool) []quantum.Gate {
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
				if includeCCNot && i < j {
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

func Search(numBits int, gates []quantum.Gate, target quantum.Gate) quantum.Circuit {
	gen := quantum.NewCircuitGen(numBits, gates, maxCircuitCache)
	hasher := quantum.NewCircuitHasher(numBits)
	goal := hasher.Hash(target)
	backwards := quantum.NewBackwardsMap(hasher, target)

	for i := 1; true; i++ {
		circuits := gen.GenerateSlice(i)
		if circuits == nil {
			ch, count := gen.Generate(i)
			fmt.Println("Doing backward search at depth", i, "with", count, "permutations...")
			for c := range ch {
				if hasher.Hash(c) == goal {
					return c
				}
				backwards.AddCircuit(c)
			}
			break
		} else {
			fmt.Println("Doing backward search of depth", i, "with", len(circuits), "permutations...")
			for _, c := range circuits {
				if hasher.Hash(c) == goal {
					return c
				}
				backwards.AddCircuit(c)
			}
		}
	}

	for i := 1; i < 100; i++ {
		ch, count := gen.Generate(i)
		fmt.Println("Doing forward search of depth", i, "with", count, "permutations...")
		for c := range ch {
			if c1 := backwards.Lookup(c); c1 != nil {
				return append(c, c1...)
			}
		}
	}

	return nil
}

func SearchSqrt(numBits int, gates []quantum.Gate, target quantum.Gate) quantum.Circuit {
	gen := quantum.NewCircuitGen(numBits, gates, maxCircuitCache)
	hasher := quantum.NewCircuitHasher(numBits)
	goal := hasher.Hash(target)

	for i := 1; i < 100; i++ {
		ch, count := gen.Generate(i)
		fmt.Println("Doing forward search of depth", i, "with", count, "permutations...")
		for c := range ch {
			c1 := append(append(quantum.Circuit{}, c...), c...)
			if hasher.Hash(c1) == goal {
				return c
			}
		}
	}

	return nil
}

func SearchCtrl(numBits int, gates []quantum.Gate, gate quantum.Gate) quantum.Circuit {
	return Search(numBits, gates, &ctrlGate{gate})
}

type ctrlGate struct {
	G quantum.Gate
}

func (c *ctrlGate) Apply(qc quantum.Computer) {
	qc.(*quantum.Simulation).ControlGate(0, c.G)
}

func (c *ctrlGate) Inverse() quantum.Gate {
	return &ctrlGate{G: c.G.Inverse()}
}

func (c *ctrlGate) String() string {
	return "Ctrl(" + c.G.String() + ")"
}
