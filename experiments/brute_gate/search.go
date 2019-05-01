package main

import (
	"fmt"

	"github.com/unixpickle/learn-quantum/quantum"
)

const (
	maxCircuitCache = 5000000
)

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
	ctx := newSearchContext(numBits, gates)
	goal := ctx.Hasher.Hash(target)
	backHasher := ctx.Hasher.Prefix(target)
	backwards := map[quantum.CircuitHash]quantum.Circuit{}

	for i := 1; i < len(ctx.Cache); i++ {
		count, ch := ctx.Enumerate(i)
		fmt.Println("Doing backward search of depth", i, "with", count, "permutations...")
		for c := range ch {
			count -= 1
			if ctx.Hasher.Hash(c) == goal {
				return c
			}
			backwards[backHasher.Hash(c.Inverse())] = c
		}
	}

	for i := 1; i < 100; i++ {
		count, ch := ctx.Enumerate(i)
		fmt.Println("Doing forward search of depth", i, "with", count, "permutations...")
		for c := range ch {
			if c1, ok := backwards[ctx.Hasher.Hash(c)]; ok {
				return append(c, c1...)
			}
		}
	}

	return nil
}

func SearchSqrt(numBits int, gates []quantum.Gate, target quantum.Gate) quantum.Circuit {
	ctx := newSearchContext(numBits, gates)
	goal := ctx.Hasher.Hash(target)

	for i := 1; i < 100; i++ {
		count, ch := ctx.Enumerate(i)
		fmt.Println("Doing forward search of depth", i, "with", count, "permutations...")
		for c := range ch {
			c1 := append(append(quantum.Circuit{}, c...), c...)
			if ctx.Hasher.Hash(c1) == goal {
				return c
			}
		}
	}

	return nil
}

func SearchCtrl(numBits int, gates []quantum.Gate, gate quantum.Gate) quantum.Circuit {
	return Search(numBits, gates, &ctrlGate{gate})
}

type searchContext struct {
	NumBits int
	Gates   []quantum.Gate
	Cache   [][]quantum.Circuit
	Hasher  quantum.CircuitHasher
}

func newSearchContext(numBits int, gates []quantum.Gate) *searchContext {
	var oneStep []quantum.Circuit
	for _, g := range gates {
		oneStep = append(oneStep, quantum.Circuit{g})
	}
	res := &searchContext{
		NumBits: numBits,
		Gates:   gates,
		Cache: [][]quantum.Circuit{
			[]quantum.Circuit{quantum.Circuit{}},
			oneStep,
		},
		Hasher: quantum.NewCircuitHasher(numBits),
	}

	var numCircuits int

CacheLoop:
	for i := 2; i <= 15; i++ {
		fmt.Println("Generating circuit cache at depth", i, "...")
		next := map[quantum.CircuitHash]quantum.Circuit{}
		_, ch := res.Enumerate(i)
		for c := range ch {
			hash := res.Hasher.Hash(c)
			if _, ok := next[hash]; !ok {
				next[hash] = c
				numCircuits++
			}
			if numCircuits > maxCircuitCache {
				break CacheLoop
			}
		}
		var circuits []quantum.Circuit
		for _, c := range next {
			circuits = append(circuits, c)
		}
		res.Cache = append(res.Cache, circuits)
	}
	return res
}

func (s *searchContext) Enumerate(numGates int) (int, <-chan quantum.Circuit) {
	if numGates < len(s.Cache) {
		cached := s.Cache[numGates]
		ch := make(chan quantum.Circuit, len(cached))
		for _, c := range cached {
			ch <- c
		}
		close(ch)
		return len(cached), ch
	}

	cached := s.Cache[len(s.Cache)-1]
	subCount, subChan := s.Enumerate(numGates - (len(s.Cache) - 1))

	ch := make(chan quantum.Circuit, 1)
	go func() {
		defer close(ch)
		for c1 := range subChan {
			for _, c2 := range cached {
				ch <- append(append(quantum.Circuit{}, c1...), c2...)
			}
		}
	}()

	return len(cached) * subCount, ch
}

type ctrlGate struct {
	G quantum.Gate
}

func (c *ctrlGate) Apply(qc quantum.Computer) {
	s1 := qc.(*quantum.Simulation)
	s2 := s1.Copy()
	for i := range s1.Phases {
		if i&1 == 0 {
			s2.Phases[i] = 0
		} else {
			s1.Phases[i] = 0
		}
	}
	c.G.Apply(s2)
	for i, phase := range s2.Phases {
		s1.Phases[i] += phase
	}
}

func (c *ctrlGate) Inverse() quantum.Gate {
	return &ctrlGate{G: c.G.Inverse()}
}

func (c *ctrlGate) String() string {
	return "Ctrl(" + c.G.String() + ")"
}
