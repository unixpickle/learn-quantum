package quantum

import (
	"math/cmplx"
	"testing"
)

func TestSqrtCNot(t *testing.T) {
	for i := 0; i < 4; i++ {
		s := NewSimulationBits(2, uint(i))
		SqrtCNot(s, 0, 1)
		SqrtCNot(s, 0, 1)
		s.CNot(0, 1)
		if cmplx.Abs(s.Phases[i]-1) > 1e-8 {
			t.Error("invalid square")
		}

		SqrtCNot(s, 0, 1)
		InvSqrtCNot(s, 0, 1)
		if cmplx.Abs(s.Phases[i]-1) > 1e-8 {
			t.Error("invalid inverse")
		}
	}
}
