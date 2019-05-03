package main

import "github.com/unixpickle/learn-quantum/quantum"

// A SlidingGate is a 4-qubit gate that is run
// sequentially with stride 2.
type SlidingGate struct {
	Gate   quantum.Gate
	Invert bool
}

func (s *SlidingGate) String() string {
	return "SlidingGate(" + s.Gate.String() + ")"
}

func (s *SlidingGate) Apply(c quantum.Computer) {
	numPairs := c.NumBits() / 2
	if s.Invert {
		inv := s.Gate.Inverse()
		for i := numPairs - 2; i >= 0; i-- {
			start := i * 2
			mapped := &quantum.MappedComputer{
				C:       c,
				Mapping: []int{start, start + 1, start + 2, start + 3},
			}
			inv.Apply(mapped)
		}
	} else {
		for i := 0; i < numPairs-1; i++ {
			start := i * 2
			mapped := &quantum.MappedComputer{
				C:       c,
				Mapping: []int{start, start + 1, start + 2, start + 3},
			}
			s.Gate.Apply(mapped)
		}
	}
}

func (s *SlidingGate) Inverse() quantum.Gate {
	return &SlidingGate{Gate: s.Gate, Invert: !s.Invert}
}

// An EndGate is a 4-qubit gate that is run on the final
// qubits in a circuit.
type EndGate struct {
	Gate quantum.Gate
}

func (e *EndGate) String() string {
	return "EndGate(" + e.Gate.String() + ")"
}

func (e *EndGate) Apply(c quantum.Computer) {
	i := c.NumBits() - 4
	mapped := &quantum.MappedComputer{
		C:       c,
		Mapping: []int{i, i + 1, i + 2, i + 3},
	}
	e.Gate.Apply(mapped)
}

func (e *EndGate) Inverse() quantum.Gate {
	return &EndGate{Gate: e.Gate.Inverse()}
}
