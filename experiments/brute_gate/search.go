package main

import (
	"crypto/md5"
	"fmt"

	"github.com/unixpickle/learn-quantum/quantum"
)

const (
	circuitChunkSize = 4
	maxBackward      = 7
)

type SimHash [md5.Size]byte

func Search(numBits int, gate func(b []bool) []bool) quantum.Circuit {
	ctx := newSearchContext(numBits, gate)
	goal := HashClassicalGate(numBits, ctx.InToOut)
	backwards := map[SimHash]quantum.Circuit{}

	for i := 1; i <= maxBackward; i++ {
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
	NumBits     int
	InToOut     []int
	Gates       []quantum.Gate
	RawCircuits [][]quantum.Circuit
}

func newSearchContext(numBits int, gate func([]bool) []bool) *searchContext {
	res := &searchContext{
		NumBits: numBits,
		InToOut: computeInToOut(numBits, gate),
		Gates:   allGates(numBits),
	}
	for i := 0; i <= circuitChunkSize; i++ {
		res.RawCircuits = append(res.RawCircuits, rawEnumerateCircuits(res.Gates, res.NumBits, i))
	}
	return res
}

func (s *searchContext) Enumerate(numGates int) (int, <-chan quantum.Circuit) {
	if numGates < len(s.RawCircuits) {
		raw := s.RawCircuits[numGates]
		ch := make(chan quantum.Circuit, len(raw))
		for _, c := range raw {
			ch <- c
		}
		close(ch)
		return len(raw), ch
	}

	raw := s.RawCircuits[circuitChunkSize]
	subCount, subChan := s.Enumerate(numGates - circuitChunkSize)

	ch := make(chan quantum.Circuit, 1)
	go func() {
		defer close(ch)
		for c1 := range subChan {
			for _, c2 := range raw {
				ch <- append(append(quantum.Circuit{}, c1...), c2...)
			}
		}
	}()

	return len(raw) * subCount, ch
}

func allGates(numBits int) []quantum.Gate {
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
			}
		}
	}
	return result
}

func rawEnumerateCircuits(gates []quantum.Gate, numBits, numGates int) []quantum.Circuit {
	if numGates == 0 {
		return []quantum.Circuit{quantum.Circuit{}}
	}
	x := map[SimHash]quantum.Circuit{}
	subCircuits := rawEnumerateCircuits(gates, numBits, numGates-1)
	for _, firstGate := range gates {
		for _, tail := range subCircuits {
			c := append(quantum.Circuit{firstGate}, tail...)
			x[HashCircuit(numBits, c)] = c
		}
	}
	var res []quantum.Circuit
	for _, c := range x {
		res = append(res, c)
	}
	return res
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
