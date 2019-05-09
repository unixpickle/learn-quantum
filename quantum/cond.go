package quantum

import (
	"math"
	"math/cmplx"
)

// Cond runs a function that only has an effect if a given
// control bit is set. The function should not attempt to
// modify the control bit.
func Cond(c Computer, control int, f func(c Computer)) {
	f(&CondComputer{Computer: c, Control: control})
}

// A CondComputer is a Computer that applies gates
// conditioned on some qubit being set. Under the hood it
// changes CNot gates to Toffoli gates, and Unitary gates
// to controlled unitary gates.
type CondComputer struct {
	Computer Computer
	Control  int
}

func (c *CondComputer) NumBits() int {
	return c.Computer.NumBits()
}

func (c *CondComputer) InUse(bit int) bool {
	return bit == c.Control || c.Computer.InUse(bit)

}

func (c *CondComputer) Measure(bitIdx int) bool {
	return c.Computer.Measure(bitIdx)
}

func (c *CondComputer) Unitary(target int, m *Matrix2) {
	if target == c.Control {
		panic("cannot change control bit")
	}
	CUnitary(c.Computer, c.Control, target, m)
}

func (c *CondComputer) CNot(control, target int) {
	if target == c.Control {
		panic("cannot change control bit")
	}
	if control == c.Control {
		c.Computer.CNot(control, target)
	} else {
		CCNot(c.Computer, c.Control, control, target)
	}
}

// CUnitary applies the unitary matrix m to the target
// qubit if the control qubit is set.
//
// The underlying implementation uses four unitaries and
// two CNot gates.
func CUnitary(c Computer, control, target int, m *Matrix2) {
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
			CUnitary(c, control, target, &Matrix2{phase, 0, 0, phase})
			CUnitary(c, control, target, &Matrix2{1, 0, 0, m.M22 / phase})
		}
		return
	} else if m.M22 == 0 && m.M11 == 0 {
		CUnitary(c, control, target, &Matrix2{m.M21, 0, 0, m.M12})
		c.CNot(control, target)
		return
	}

	// https://arxiv.org/abs/quant-ph/9503016

	theta := 2 * math.Atan2(cmplx.Abs(m.M12), cmplx.Abs(m.M11))

	a := imag(cmplx.Log(m.M11 / -m.M21))
	var b, delta float64
	if cmplx.Abs(m.M11) > cmplx.Abs(m.M21) {
		b = imag(cmplx.Log((m.M11 / m.M22) * cmplx.Exp(complex(0, -a))))
		delta = imag(cmplx.Log(m.M11 / cmplx.Exp(complex(0, a/2+b/2))))
	} else {
		b = imag(cmplx.Log((-m.M21 / m.M12) * cmplx.Exp(complex(0, a))))
		delta = imag(cmplx.Log(m.M12 / cmplx.Exp(complex(0, a/2-b/2))))
	}

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
	CUnitary(c, control, target, &mat4)
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
