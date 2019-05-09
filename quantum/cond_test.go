package quantum

import (
	"math"
	"math/cmplx"
	"math/rand"
	"testing"
)

func TestCondUnitary(t *testing.T) {
	testMat := func(t *testing.T, m *Matrix2) {
		s1 := RandomSimulation(3)
		s2 := s1.Copy()
		rawCondUnitary(s1, 2, 1, m)
		CondUnitary(s2, 2, 1, m)
		if !s1.ApproxEqual(s2, 1e-8) {
			t.Fatal("incorrect result")
		}
	}
	t.Run("Diagonal", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			mat := Matrix2{cmplx.Exp(complex(0, rand.Float64()*2*math.Pi)), 0, 0,
				cmplx.Exp(complex(0, rand.Float64()*2*math.Pi))}
			testMat(t, &mat)
		}
	})
	t.Run("OffDiagonal", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			mat := Matrix2{0, cmplx.Exp(complex(0, rand.Float64()*2*math.Pi)),
				cmplx.Exp(complex(0, rand.Float64()*2*math.Pi)), 0}
			testMat(t, &mat)
		}
	})
	t.Run("Random", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			mat := RandomMatrix2()
			testMat(t, &mat)
		}
	})
}

func rawCondUnitary(s *Simulation, control, target int, m *Matrix2) {
	for i := range s.Phases {
		if i&(1<<uint(target)) != 0 || i&(1<<uint(control)) == 0 {
			continue
		}
		other := i | (1 << uint(target))
		p0 := s.Phases[i]
		p1 := s.Phases[other]
		s.Phases[i] = m.M11*p0 + m.M12*p1
		s.Phases[other] = m.M21*p0 + m.M22*p1
	}
}
