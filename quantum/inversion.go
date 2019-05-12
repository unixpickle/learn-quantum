package quantum

// Invert applies the inverse of the function f to the
// computer c.
func Invert(c Computer, f func(c Computer)) {
	tape := &invertTape{c: c, inverse: func() {}}
	f(tape)
	tape.inverse()
}

// Conj applies a "conjugate". In other words, it applies
// the function a, then the function b, then the inverse
// of a. This can be used for clean computation.
func Conj(c Computer, a func(c Computer), b func(c Computer)) {
	tape := &invertTape{c: c, inverse: func() {}, forward: true}
	a(tape)
	b(c)
	tape.inverse()
}

type invertTape struct {
	c       Computer
	inverse func()
	forward bool
}

func (i *invertTape) NumBits() int {
	return i.c.NumBits()
}

func (i *invertTape) InUse(bit int) bool {
	return i.c.InUse(bit)
}

func (i *invertTape) Measure(bitIdx int) bool {
	panic("measurement is not invertible")
}

func (i *invertTape) Unitary(target int, m *Matrix2) {
	m1 := *m
	m1.ConjTranspose()
	oldInv := i.inverse
	i.inverse = func() {
		i.c.Unitary(target, &m1)
		oldInv()
	}
	if i.forward {
		i.c.Unitary(target, m)
	}
}

func (i *invertTape) CNot(control, target int) {
	oldInv := i.inverse
	i.inverse = func() {
		i.c.CNot(control, target)
		oldInv()
	}
	if i.forward {
		i.c.CNot(control, target)
	}
}
