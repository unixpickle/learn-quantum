package quantum

import (
	"math/cmplx"
	"testing"
)

func TestRandomMatrix2(t *testing.T) {
	for i := 0; i < 10; i++ {
		m := RandomMatrix2()
		mH := m
		mH.ConjTranspose()
		m.Mul(&mH)
		if cmplx.Abs(m.M11-1) > 1e-8 || cmplx.Abs(m.M22-1) > 1e-8 {
			t.Error("invalid diagonal", m.M11, m.M22)
		}
		if cmplx.Abs(m.M12) > 1e-8 || cmplx.Abs(m.M21) > 1e-8 {
			t.Error("invalid off-diagonal", m.M12, m.M21)
		}
	}
}

func TestMatrix2Sqrt(t *testing.T) {
	m := RandomMatrix2()
	s := m
	s.Sqrt()
	s.Mul(&s)

	if cmplx.Abs(m.M11-s.M11) > 1e-8 || cmplx.Abs(m.M12-s.M12) > 1e-8 ||
		cmplx.Abs(m.M21-s.M21) > 1e-8 || cmplx.Abs(m.M22-s.M22) > 1e-8 {
		t.Error("incorrect square")
	}
}
