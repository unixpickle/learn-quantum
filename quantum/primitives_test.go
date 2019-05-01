package quantum

import (
	"math/cmplx"
	"testing"
)

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
