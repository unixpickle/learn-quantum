package quantum

import "math/rand"

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
