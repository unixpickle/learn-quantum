package quantum

import (
	"math/rand"
	"testing"

	"github.com/unixpickle/essentials"
)

func TestToffoliN(t *testing.T) {
	for i := 0; i < 1000; i++ {
		numBits := rand.Intn(10) + 1
		target := rand.Intn(numBits)
		var numControl int
		if numBits <= 3 {
			numControl = rand.Intn(numBits)
		} else {
			numControl = rand.Intn(numBits - 2)
		}
		var control []int
		for j := 0; j < numControl; j++ {
			idx := rand.Intn(numBits)
			for idx == target || essentials.Contains(control, idx) {
				idx = rand.Intn(numBits)
			}
			control = append(control, idx)
		}
		s := RandomSimulation(numBits)
		expected := rawToffoliN(s, target, control)
		ToffoliN(s, target, control...)
		if expected.String() != s.String() {
			t.Errorf("error for %d bits with target %d and control %v", numBits, target, control)
		}
	}
}

func rawToffoliN(s *Simulation, target int, control []int) *Simulation {
	s1 := s.Copy()
	for i, phase := range s.Phases {
		matches := true
		for _, x := range control {
			if (i & (1 << uint(x))) == 0 {
				matches = false
				break
			}
		}
		if matches {
			s1.Phases[i^(1<<uint(target))] = phase
		}
	}
	return s1
}
