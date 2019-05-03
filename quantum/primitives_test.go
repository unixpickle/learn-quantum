package quantum

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"
)

func BenchmarkUnitary(b *testing.B) {
	for _, size := range []int{1, 5, 10} {
		b.Run(fmt.Sprintf("Bits%d", size), func(b *testing.B) {
			s := RandomSimulation(size)
			coeff := complex(1.0/math.Sqrt2, 0)
			for i := 0; i < b.N; i++ {
				s.Unitary(i%size, coeff, coeff, coeff, -coeff)
			}
		})
	}
}

func BenchmarkCNot(b *testing.B) {
	for _, size := range []int{2, 5, 10} {
		b.Run(fmt.Sprintf("Bits%d", size), func(b *testing.B) {
			s := RandomSimulation(size)
			for i := 0; i < b.N; i++ {
				s.CNot(i%size, (i+1)%size)
			}
		})
	}
}

func TestCSqrtNot(t *testing.T) {
	testControl(t, SqrtNot, InvSqrtNot, CSqrtNot, InvCSqrtNot)
}

func TestCH(t *testing.T) {
	testControl(t, H, H, CH, CH)
}

func TestSqrtSwap(t *testing.T) {
	for i := 0; i < 4; i++ {
		s := NewSimulationBits(2, uint(i))
		SqrtSwap(s, 0, 1)
		SqrtSwap(s, 0, 1)
		Swap(s, 0, 1)
		if cmplx.Abs(s.Phases[i]-1) > 1e-8 {
			t.Error("invalid square")
		}

		SqrtSwap(s, 0, 1)
		InvSqrtSwap(s, 0, 1)
		if cmplx.Abs(s.Phases[i]-1) > 1e-8 {
			t.Error("invalid inverse")
		}
	}
}

func TestCSwap(t *testing.T) {
	g1 := &CSwapGate{1, 2, 3}
	g2 := &ClassicalGate{
		F: func(b []bool) []bool {
			res := append([]bool{}, b...)
			if res[1] {
				res[2], res[3] = res[3], res[2]
			}
			return res
		},
	}
	hasher := NewCircuitHasher(4)
	if hasher.Hash(g1) != hasher.Hash(g2) {
		t.Error("invalid circuit")
	}
}

func testControl(t *testing.T, fwd func(c Computer, bit int), inv func(c Computer, bit int),
	controlled func(c Computer, ctrl, targ int),
	controlledInv func(c Computer, ctrl, targ int)) {
	rawGate := &FnGate{
		Forward: func(c Computer) {
			fwd(c, 2)
		},
		Backward: func(c Computer) {
			inv(c, 2)
		},
	}
	idealGate := &FnGate{
		Forward: func(c Computer) {
			c.(*Simulation).ControlGate(0, rawGate)
		},
		Backward: func(c Computer) {
			c.(*Simulation).ControlGate(0, rawGate.Inverse())
		},
	}
	actualGate := &FnGate{
		Forward: func(c Computer) {
			controlled(c, 0, 2)
		},
		Backward: func(c Computer) {
			controlledInv(c, 0, 2)
		},
	}

	hasher := NewCircuitHasher(3)
	hash1 := hasher.Hash(idealGate)
	hash2 := hasher.Hash(actualGate)
	if hash1 != hash2 {
		t.Error("invalid forward hash")
	}

	hash1 = hasher.Hash(idealGate.Inverse())
	hash2 = hasher.Hash(actualGate.Inverse())
	if hash1 != hash2 {
		t.Error("invalid backward hash")
	}
}
