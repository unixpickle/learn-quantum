package quantum

import (
	"math"
	"math/cmplx"
	"math/rand"
)

type Matrix2 struct {
	// Row 1, column 1
	M11 complex128

	// Row 1, column 2
	M12 complex128

	// Row 2, column 1
	M21 complex128

	// Row 2, column 2.
	M22 complex128
}

// NewMatrix2 creates the identity.
func NewMatrix2() Matrix2 {
	return Matrix2{1, 0, 0, 1}
}

// RandomMatrix2 creates a random unitary 2x2 matrix.
func RandomMatrix2() Matrix2 {
	var m Matrix2

	nums := []*complex128{&m.M11, &m.M12, &m.M21, &m.M22}
	for _, num := range nums {
		*num = complex(rand.NormFloat64(), rand.NormFloat64())
	}

	// Normalize the first column
	norm := complex(math.Sqrt(math.Pow(cmplx.Abs(m.M11), 2)+math.Pow(cmplx.Abs(m.M21), 2)), 0)
	m.M11 /= norm
	m.M21 /= norm

	// Project the first column out of the second.
	dot := m.M12*cmplx.Conj(m.M11) + m.M22*cmplx.Conj(m.M21)
	m.M12 -= m.M11 * dot
	m.M22 -= m.M21 * dot

	// Normalize the second column
	norm = complex(math.Sqrt(math.Pow(cmplx.Abs(m.M12), 2)+math.Pow(cmplx.Abs(m.M22), 2)), 0)
	m.M12 /= norm
	m.M22 /= norm

	return m
}

func (m *Matrix2) ConjTranspose() {
	m.M12, m.M21 = m.M21, m.M12
	m.M11 = cmplx.Conj(m.M11)
	m.M12 = cmplx.Conj(m.M12)
	m.M21 = cmplx.Conj(m.M21)
	m.M22 = cmplx.Conj(m.M22)
}

func (m *Matrix2) Mul(other *Matrix2) {
	m.M11, m.M12, m.M21, m.M22 = m.M11*other.M11+m.M12*other.M21,
		m.M11*other.M12+m.M12*other.M22,
		m.M21*other.M11+m.M22*other.M21,
		m.M21*other.M12+m.M22*other.M22
}

func (m *Matrix2) Sub(other *Matrix2) {
	m.M11, m.M12, m.M21, m.M22 = m.M11-other.M11, m.M12-other.M12, m.M21-other.M21,
		m.M22-other.M22
}

func (m *Matrix2) Sqrt() {
	det := m.M11*m.M22 - m.M12*m.M21
	trace := m.M11 + m.M22

	// Using formula from https://en.wikipedia.org/wiki/Square_root_of_a_2_by_2_matrix

	s := cmplx.Sqrt(det)
	t := cmplx.Sqrt(trace + 2*s)

	m.M11 = (m.M11 + s) / t
	m.M12 = m.M12 / t
	m.M21 = m.M21 / t
	m.M22 = (m.M22 + s) / t
}
