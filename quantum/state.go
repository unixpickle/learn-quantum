package quantum

import (
	"math"
	"math/rand"
)

type State struct {
	phases  []complex128
	numBits int
}

func NewState(numBits int) *State {
	s := &State{
		phases:  make([]complex128, 1<<uint(numBits)),
		numBits: numBits,
	}
	s.phases[0] = 1
	return s
}

func (s *State) Sample() uint64 {
	x := rand.Float64()
	for i, phase := range s.phases {
		v := real(phase)*real(phase) + imag(phase)*imag(phase)
		x -= v
		if x <= 0 {
			return uint64(i)
		}
	}
	return uint64(len(s.phases) - 1)
}

func (s *State) Not(bitIdx int) *State {
	if bitIdx < 0 || bitIdx >= s.numBits {
		panic("bit index out of range")
	}
	return s.Map(func(x uint64) uint64 {
		return x ^ (1 << uint(bitIdx))
	})
}

func (s *State) CNot(source, target int) *State {
	if source < 0 || source >= s.numBits || target < 0 || target >= s.numBits {
		panic("bit index out of range")
	}
	return s.Map(func(x uint64) uint64 {
		b1 := (x & (1 << uint(source))) >> uint(source)
		return x ^ (b1 << uint(target))
	})
}

func (s *State) Hadamard(bitIdx int) *State {
	if bitIdx < 0 || bitIdx >= s.numBits {
		panic("bit index out of range")
	}
	s1 := NewState(s.numBits)
	s1.phases[0] = 0
	coeff := complex(1/math.Sqrt2, 0)
	for i, phase := range s.phases {
		flip := i ^ (1 << uint(bitIdx))
		if i&(1<<uint(bitIdx)) == 0 {
			s1.phases[i] += phase * coeff
			s1.phases[flip] += phase * coeff
		} else {
			s1.phases[i] -= phase * coeff
			s1.phases[flip] += phase * coeff
		}
	}
	return s1
}

func (s *State) Map(f func(uint64) uint64) *State {
	s1 := NewState(s.numBits)
	s1.phases[0] = 0
	for i, phase := range s.phases {
		s1.phases[int(f(uint64(i)))] += phase
	}
	return s1
}

func getBit(value uint64, idx int) bool {
	return value<<uint(idx) != 0
}

func setBit(value uint64, idx int, b bool) uint64 {
	if getBit(value, idx) == b {
		return value
	}
	return value ^ (1 << uint(idx))
}
