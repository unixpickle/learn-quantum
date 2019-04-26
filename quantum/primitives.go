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

// SqrtNot performs the square root of the Not gate.
func SqrtNot(c Computer, bitIdx int) {
	H(c, bitIdx)
	TInv(c, bitIdx)
	TInv(c, bitIdx)
	H(c, bitIdx)
}

// InvSqrtNot performs the inverse of SqrtNot.
func InvSqrtNot(c Computer, bitIdx int) {
	H(c, bitIdx)
	T(c, bitIdx)
	T(c, bitIdx)
	H(c, bitIdx)
}

// SqrtCNot performs the square root of the CNot gate.
func SqrtCNot(c Computer, control, target int) {
	// Found via search.
	TInv(c, 0)
	H(c, 1)
	c.CNot(0, 1)
	T(c, 1)
	c.CNot(0, 1)
	TInv(c, 1)
	H(c, 1)
}

// InvSqrtCNot performs the inverse of SqrtCNot.
func InvSqrtCNot(c Computer, control, target int) {
	H(c, 1)
	T(c, 1)
	c.CNot(0, 1)
	TInv(c, 1)
	c.CNot(0, 1)
	H(c, 1)
	T(c, 0)
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
	TInv(c, a)
	TInv(c, a)
	c.CNot(a, b)
	H(c, a)
	TInv(c, b)
	c.CNot(a, b)
	T(c, b)
	T(c, b)
}

// InvSqrtSwap perfroms the inverse of SqrtSwap.
func InvSqrtSwap(c Computer, a, b int) {
	TInv(c, b)
	TInv(c, b)
	c.CNot(a, b)
	T(c, b)
	H(c, a)
	c.CNot(a, b)
	T(c, a)
	T(c, a)
}
