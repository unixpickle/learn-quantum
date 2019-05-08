package quantum

import (
	"math/cmplx"
	"testing"
)

func TestRandomUnitary2(t *testing.T) {
	for i := 0; i < 10; i++ {
		m11, m12, m21, m22 := randomUnitary2()

		s11 := m11*cmplx.Conj(m11) + m12*cmplx.Conj(m12)
		s12 := m11*cmplx.Conj(m21) + m12*cmplx.Conj(m22)
		s21 := m21*cmplx.Conj(m11) + m22*cmplx.Conj(m12)
		s22 := m21*cmplx.Conj(m21) + m22*cmplx.Conj(m22)
		if cmplx.Abs(s11-1) > 1e-8 || cmplx.Abs(s22-1) > 1e-8 {
			t.Error("invalid diagonal", s11, s22)
		}
		if cmplx.Abs(s12) > 1e-8 || cmplx.Abs(s21) > 1e-8 {
			t.Error("invalid off-diagonal", s12, s21)
		}
	}
}
