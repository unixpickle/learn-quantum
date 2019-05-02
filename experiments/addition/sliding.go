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
