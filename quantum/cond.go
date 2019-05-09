package quantum

import (
	"math"
	"math/cmplx"
)

func CondUnitary(c Computer, control, target int, m *Matrix2) {
	if m.M12 == 0 && m.M21 == 0 {
		if m.M11 == m.M22 {
			c.Unitary(control, &Matrix2{1, 0, 0, m.M11})
		} else if m.M11 == 1 {
			sqrt := *m
			sqrt.M22 = cmplx.Sqrt(sqrt.M22)
			sqrtInv := sqrt
			sqrtInv.M22 = cmplx.Conj(sqrtInv.M22)
			c.Unitary(control, &sqrt)
			c.CNot(control, target)
			c.Unitary(target, &sqrtInv)
			c.CNot(control, target)
			c.Unitary(target, &sqrt)
		} else {
			phase := m.M11
			CondUnitary(c, control, target, &Matrix2{phase, 0, 0, phase})
			CondUnitary(c, control, target, &Matrix2{1, 0, 0, m.M22 / phase})
		}
		return
	} else if m.M22 == 0 && m.M11 == 0 {
		CondUnitary(c, control, target, &Matrix2{m.M21, 0, 0, m.M12})
		c.CNot(control, target)
		return
	}

	// https://arxiv.org/abs/quant-ph/9503016

	theta := 2 * math.Atan2(cmplx.Abs(m.M12), cmplx.Abs(m.M11))

	// TODO: more numerically-stable formulae here.
	b := imag(cmplx.Log(m.M11 / m.M12))
	a := imag(cmplx.Log(m.M11 / -m.M21))
	delta := imag(cmplx.Log(m.M11 / cmplx.Exp(complex(0, a/2+b/2))))

	mat1 := rotateZ(a)
	mat1A := rotateY(theta / 2)
	mat1.Mul(&mat1A)

	mat2 := rotateY(-theta / 2)
	mat2A := rotateZ(-(a + b) / 2)
	mat2.Mul(&mat2A)

	mat3 := rotateZ((b - a) / 2)
	mat4 := phaseShift(delta)

	c.Unitary(target, &mat3)
	c.CNot(control, target)
	c.Unitary(target, &mat2)
	c.CNot(control, target)
	c.Unitary(target, &mat1)
	CondUnitary(c, control, target, &mat4)
}

func rotateY(theta float64) Matrix2 {
	cos := complex(math.Cos(theta/2), 0)
	sin := complex(math.Sin(theta/2), 0)
	return Matrix2{cos, sin, -sin, cos}
}

func rotateZ(alpha float64) Matrix2 {
	num := cmplx.Exp(complex(0, alpha/2))
	return Matrix2{num, 0, 0, cmplx.Conj(num)}
}

func phaseShift(delta float64) Matrix2 {
	num := cmplx.Exp(complex(0, delta))
	return Matrix2{num, 0, 0, num}
}
