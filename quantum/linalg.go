package quantum

import (
	"math"
	"math/cmplx"
	"math/rand"
)

func randomUnitary2() (m11, m12, m21, m22 complex128) {
	nums := []*complex128{&m11, &m12, &m21, &m22}
	for _, num := range nums {
		*num = complex(rand.NormFloat64(), rand.NormFloat64())
	}

	// Normalize the first column
	norm := complex(math.Sqrt(math.Pow(cmplx.Abs(m11), 2)+math.Pow(cmplx.Abs(m21), 2)), 0)
	m11 /= norm
	m21 /= norm

	// Project the first column out of the second.
	dot := m12*cmplx.Conj(m11) + m22*cmplx.Conj(m21)
	m12 -= m11 * dot
	m22 -= m21 * dot

	// Normalize the second column
	norm = complex(math.Sqrt(math.Pow(cmplx.Abs(m12), 2)+math.Pow(cmplx.Abs(m22), 2)), 0)
	m12 /= norm
	m22 /= norm

	return
}
