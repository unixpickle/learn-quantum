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

// InvT performs an inverse rotation by pi/4.
func InvT(c Computer, bitIdx int) {
	c.Unitary(bitIdx, 1, 0, 0, cmplx.Exp(complex(0, -math.Pi/4)))
}

// SqrtT performs the positive square root of the T gate.
func SqrtT(c Computer, bitIdx int) {
	c.Unitary(bitIdx, 1, 0, 0, cmplx.Exp(complex(0, math.Pi/8)))
}

// InvSqrtT performs the negative square root of the T gate.
func InvSqrtT(c Computer, bitIdx int) {
	c.Unitary(bitIdx, 1, 0, 0, cmplx.Exp(complex(0, -math.Pi/8)))
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

// SqrtNot performs the square root of the Not gate.
func SqrtNot(c Computer, bitIdx int) {
	H(c, bitIdx)
	InvT(c, bitIdx)
	InvT(c, bitIdx)
	H(c, bitIdx)
}

// InvSqrtNot performs the inverse of SqrtNot.
func InvSqrtNot(c Computer, bitIdx int) {
	H(c, bitIdx)
	T(c, bitIdx)
	T(c, bitIdx)
	H(c, bitIdx)
}

// CSqrtNot performs a controlled SqrtNot gate.
func CSqrtNot(c Computer, control, target int) {
	// Found via search.
	H(c, target)
	InvT(c, control)
	c.CNot(control, target)
	T(c, target)
	c.CNot(control, target)
	InvT(c, target)
	H(c, target)
}

// InvCSqrtNot is the inverse of CSqrtNot.
func InvCSqrtNot(c Computer, control, target int) {
	H(c, target)
	T(c, target)
	c.CNot(control, target)
	InvT(c, target)
	c.CNot(control, target)
	T(c, control)
	H(c, target)
}

// Swap swaps two qubits.
func Swap(c Computer, a, b int) {
	c.CNot(a, b)
	c.CNot(b, a)
	c.CNot(a, b)
}

// SqrtSwap perfroms the square root of the Swap gate.
func SqrtSwap(c Computer, a, b int) {
	// Found via search.
	InvT(c, a)
	InvT(c, a)
	c.CNot(a, b)
	H(c, a)
	InvT(c, b)
	c.CNot(a, b)
	T(c, b)
	T(c, b)
}

// InvSqrtSwap perfroms the inverse of SqrtSwap.
func InvSqrtSwap(c Computer, a, b int) {
	InvT(c, b)
	InvT(c, b)
	c.CNot(a, b)
	T(c, b)
	H(c, a)
	c.CNot(a, b)
	T(c, a)
	T(c, a)
}

// CH applies a controlled Hadamard gate.
func CH(c Computer, control, target int) {
	// Found via search.
	T(c, target)
	T(c, target)
	H(c, target)
	T(c, target)
	c.CNot(control, target)
	InvT(c, target)
	H(c, target)
	InvT(c, target)
	InvT(c, target)
}
