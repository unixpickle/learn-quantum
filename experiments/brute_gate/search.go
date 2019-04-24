package main

import (
	"crypto/md5"
	"fmt"

	"github.com/unixpickle/learn-quantum/quantum"
)

const (
	maxCircuitCache = 5000000
)

type SimHash [md5.Size]byte

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

func Search(numBits int, gates []quantum.Gate, gate func(b []bool) []bool) quantum.Circuit {
	ctx := newSearchContext(numBits, gates, gate)
	goal := HashClassicalGate(numBits, ctx.InToOut)
	backwards := map[SimHash]quantum.Circuit{}

	for i := 1; i < len(ctx.Cache); i++ {
		count, ch := ctx.Enumerate(i)
		fmt.Println("Doing backward search of depth", i, "with", count, "permutations...")
		for c := range ch {
			count -= 1
			if HashCircuit(numBits, c) == goal {
				return c
			}
			backwards[HashCircuitBackwards(numBits, c, ctx.InToOut)] = c
		}
	}

	for i := 1; i < 100; i++ {
		count, ch := ctx.Enumerate(i)
		fmt.Println("Doing forward search of depth", i, "with", count, "permutations...")
		for c := range ch {
			if c1, ok := backwards[HashCircuit(numBits, c)]; ok {
				return append(c, c1...)
			}
		}
	}

	return nil
}

type searchContext struct {
	NumBits int
	InToOut []int
	Gates   []quantum.Gate
	Cache   [][]quantum.Circuit
}

func newSearchContext(numBits int, gates []quantum.Gate, gate func([]bool) []bool) *searchContext {
	var oneStep []quantum.Circuit
	for _, g := range gates {
		oneStep = append(oneStep, quantum.Circuit{g})
	}
	res := &searchContext{
		NumBits: numBits,
		InToOut: computeInToOut(numBits, gate),
		Gates:   gates,
		Cache: [][]quantum.Circuit{
			[]quantum.Circuit{quantum.Circuit{}},
			oneStep,
		},
	}
	var numCircuits int

CacheLoop:
	for i := 2; i <= 100; i++ {
		fmt.Println("Generating circuit cache at depth", i, "...")
		next := map[SimHash]quantum.Circuit{}
		_, ch := res.Enumerate(i)
		for c := range ch {
			hash := HashCircuit(res.NumBits, c)
			if _, ok := next[hash]; !ok {
				next[HashCircuit(res.NumBits, c)] = c
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

func computeInToOut(numBits int, gate func(b []bool) []bool) []int {
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
	return inToOut
}
