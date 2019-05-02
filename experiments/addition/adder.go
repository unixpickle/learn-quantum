package main

import (
	"github.com/unixpickle/learn-quantum/quantum"
)

type AdderGate struct{}

func (_ AdderGate) String() string {
	return "Adder"
}

func (_ AdderGate) Apply(c quantum.Computer) {
	s := c.(*quantum.Simulation)
	s1 := s.Copy()
	for i, phase := range s1.Phases {
		a, b := adderPairs(s.NumBits(), i)
		dest := replaceAddResult(s.NumBits(), i, a+b)
		s.Phases[dest] = phase
	}
}

func (_ AdderGate) Inverse() quantum.Gate {
	return SubtractorGate{}
}

type SubtractorGate struct{}

func (_ SubtractorGate) String() string {
	return "Subtractor"
}

func (_ SubtractorGate) Apply(c quantum.Computer) {
	s := c.(*quantum.Simulation)
	s1 := s.Copy()
	for i, phase := range s1.Phases {
		a, b := adderPairs(s.NumBits(), i)
		dest := replaceAddResult(s.NumBits(), i, b-a)
		s.Phases[dest] = phase
	}
}

func (_ SubtractorGate) Inverse() quantum.Gate {
	return AdderGate{}
}

func adderPairs(numBits, state int) (uint32, uint32) {
	var state1 uint32
	var state2 uint32
	for i := 0; i < numBits; i++ {
		if i&1 == 0 {
			state1 |= uint32(state&(1<<uint(i))) >> uint(i/2)
		} else {
			state2 |= uint32(state&(1<<uint(i))) >> uint(i/2+1)
		}
	}
	return state1, state2
}

func replaceAddResult(numBits, state int, state2 uint32) int {
	for i := 1; i < numBits; i += 2 {
		mask := 1 << uint(i)
		if state2&(1<<uint(i/2)) == 0 {
			state &= ^mask
		} else {
			state |= mask
		}
	}
	return state
}
