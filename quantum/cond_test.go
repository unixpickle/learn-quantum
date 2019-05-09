package quantum

import (
	"math"
	"math/cmplx"
	"math/rand"
	"testing"
)

func TestCUnitary(t *testing.T) {
	testMat := func(t *testing.T, m *Matrix2) {
		s1 := RandomSimulation(3)
		s2 := s1.Copy()
		rawCUnitary(s1, 2, 1, m)
		CUnitary(s2, 2, 1, m)
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
	t.Run("NumericalErrors", func(t *testing.T) {
		num := complex(1/math.Sqrt2, 0)
		h := Matrix2{num, num, num, -num}
		T := Matrix2{1, 0, 0, cmplx.Exp(complex(0, math.Pi/4))}
		invT := Matrix2{1, 0, 0, cmplx.Exp(complex(0, -math.Pi/4))}

		// Create an _almost_ diagonal matrix with some
		// off-diagonal terms on the order of 1e-17.
		res := T
		res.Mul(&h)
		res.Mul(&T)
		res.Mul(&h)
		res.Mul(&T)
		res.Mul(&invT)
		res.Mul(&h)
		res.Mul(&invT)
		res.Mul(&h)

		testMat(t, &res)

		// Swap the rows and try again.
		res.M11, res.M12, res.M21, res.M22 = res.M21, res.M22, res.M11, res.M12
		testMat(t, &res)
	})
}

func rawCUnitary(s *Simulation, control, target int, m *Matrix2) {
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
