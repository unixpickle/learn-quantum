package quantum

import "math"

// Hadamard performs an H gate.
func Hadamard(c Computer, bitIdx int) {
	s := complex(1/math.Sqrt2, 0)
	c.Unitary(bitIdx, s, s, s, -s)
}

// Not performs a not gate.
func Not(c Computer, bitIdx int) {
	c.Unitary(bitIdx, 0, 1, 1, 0)
}
