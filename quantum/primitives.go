package quantum

import "math"

// Hadamard performs an H gate.
func Hadamard(c Computer, bitIdx int) {
	s := complex(1/math.Sqrt2, 0)
	c.Unitary(bitIdx, s, s, s, -s)
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
