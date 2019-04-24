package quantum

import (
	"math"
	"math/cmplx"
)

// H performs a Hadamard gate.
func H(c Computer, bitIdx int) {
	s := complex(1/math.Sqrt2, 0)
	c.Unitary(bitIdx, s, s, s, -s)
}

// T performs a rotation by pi/4.
func T(c Computer, bitIdx int) {
	c.Unitary(bitIdx, 1, 0, 0, cmplx.Exp(complex(0, math.Pi/4)))
}

// TInv performs an inverse rotation by pi/4.
func TInv(c Computer, bitIdx int) {
	c.Unitary(bitIdx, 1, 0, 0, cmplx.Exp(complex(0, -math.Pi/4)))
}

// X performs a not gate.
func X(c Computer, bitIdx int) {
	c.Unitary(bitIdx, 0, 1, 1, 0)
}

// Y performs a Pauli Y-gate.
func Y(c Computer, bitIdx int) {
	c.Unitary(bitIdx, 0, complex(0, -1), complex(0, 1), 0)
}

// Z performs a Pauli Z-gate.
func Z(c Computer, bitIdx int) {
	c.Unitary(bitIdx, 1, 0, 0, -1)
}

// SqrtCNot performs the square root of the CNot gate.
func SqrtCNot(c Computer, control, target int) {
	// Found via search.
	H(c, target)
	TInv(c, target)
	c.CNot(target, control)
	T(c, control)
	X(c, control)
	TInv(c, control)
	H(c, target)
}

// InvSqrtCNot performs the inverse of SqrtCNot.
func InvSqrtCNot(c Computer, control, target int) {
	H(c, target)
	T(c, control)
	X(c, control)
	TInv(c, control)
	c.CNot(target, control)
	T(c, target)
	H(c, target)
}

// CCNot performs a Toffoli gate.
func CCNot(c Computer, control1, control2, target int) {
	// https://quantum.country/qcvc
	H(c, target)
	c.CNot(control2, target)
	TInv(c, target)
	c.CNot(control1, target)
	T(c, target)
	c.CNot(control2, target)
	TInv(c, target)
	c.CNot(control1, target)
	T(c, control2)
	T(c, target)
	H(c, target)
	c.CNot(control1, control2)
	T(c, control1)
	TInv(c, control2)
	c.CNot(control1, control2)
}
