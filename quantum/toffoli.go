package quantum

import "math"

// CCNot performs a Toffoli gate.
func CCNot(c Computer, control1, control2, target int) {
	// https://quantum.country/qcvc
	H(c, target)
	c.CNot(control2, target)
	InvT(c, target)
	c.CNot(control1, target)
	T(c, target)
	c.CNot(control2, target)
	InvT(c, target)
	c.CNot(control1, target)
	T(c, control2)
	T(c, target)
	H(c, target)
	c.CNot(control1, control2)
	T(c, control1)
	InvT(c, control2)
	c.CNot(control1, control2)
}

// ToffoliN performs an n-bit Toffoli gate.
//
// This typically requires that there is at least one
// spare qubit on the computer.
func ToffoliN(c Computer, target int, control ...int) {
	if len(control) == 0 {
		X(c, target)
	} else if len(control) == 1 {
		c.CNot(control[0], target)
	} else if len(control) == 2 {
		CCNot(c, control[0], control[1], target)
	} else {
		working := allocWorking(c, target, control...)
		if len(working) == 0 {
			panic("not enough working qubits")
		} else if len(working) >= len(control)-2 {
			// Lemma 7.2 from Barenco et al. 1995
			targets := append(append([]int{}, working[:len(control)-2]...), target)

			for i := len(control) - 1; i > 1; i-- {
				CCNot(c, control[i], targets[i-2], targets[i-1])
			}
			CCNot(c, control[0], control[1], targets[0])
			for i := 2; i < len(control); i++ {
				CCNot(c, control[i], targets[i-2], targets[i-1])
			}

			// Undo side-effects
			for i := len(control) - 2; i > 1; i-- {
				CCNot(c, control[i], targets[i-2], targets[i-1])
			}
			CCNot(c, control[0], control[1], targets[0])
			for i := 2; i < len(control)-1; i++ {
				CCNot(c, control[i], targets[i-2], targets[i-1])
			}
		} else {
			// Lemma 7.3 from Barenco et al. 1995
			size := int(math.Ceil(float64(len(control)+1) / 2))
			control1 := control[:size]
			control2 := append([]int{working[0]}, control[size:]...)
			ToffoliN(c, working[0], control1...)
			ToffoliN(c, target, control2...)
			ToffoliN(c, working[0], control1...)
			ToffoliN(c, target, control2...)
		}
	}
}

// allocWorking finds the available working bits.
func allocWorking(c Computer, a int, b ...int) []int {
	var res []int
OuterLoop:
	for i := 0; i < c.NumBits(); i++ {
		if i == a || c.InUse(i) {
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
