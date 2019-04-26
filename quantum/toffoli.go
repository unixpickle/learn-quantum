package quantum

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

// CCCNot performs a four-bit Toffoli gate.
//
// This requires that there is some spare qubit on the
// computer not involved in this gate.
func CCCNot(c Computer, control1, control2, control3, target int) {
	working := allocWorking(c, control1, control2, control3, target)
	if len(working) == 0 {
		panic("an extra qubit is required")
	}
	for i := 0; i < 2; i++ {
		CCNot(c, control1, control2, working[0])
		CCNot(c, control3, working[0], target)
	}
}

// allocWorking finds the available working bits.
func allocWorking(c Computer, a int, b ...int) []int {
	var res []int
OuterLoop:
	for i := 0; i < c.NumBits(); i++ {
		if i == a {
			continue
		}
		for _, x := range b {
			if i == x {
				continue OuterLoop
			}
		}
		res = append(res, i)
	}
	return res
}
