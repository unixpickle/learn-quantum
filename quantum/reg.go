package quantum

// A Reg is a "quantum register". It is simply a list of
// distinct qubit indices.
type Reg []int

// Valid ensures that no qubits are repeated in the list
// and that no qubit indices are less than 0.
func (r Reg) Valid() bool {
	set := map[int]bool{}
	for _, x := range r {
		if set[x] || x < 0 {
			return false
		}
		set[x] = true
	}
	return true
}

// Overlaps checks if r shares any bits in common with r1.
func (r Reg) Overlaps(r1 Reg) bool {
	for _, x := range r {
		for _, y := range r1 {
			if x == y {
				return true
			}
		}
	}
	return false
}

// Extract extracts the register from a computer state.
// The bits in the computer state are stored lowest to
// highest, as are the bits in the register.
func (r Reg) Extract(state uint) uint {
	var res uint
	for i, x := range r {
		if state&(1<<uint(x)) != 0 {
			res |= 1 << uint(i)
		}
	}
	return res
}

// Inject performs the inverse of Extract. It sets the
// value of the register in a state, and returns the new
// state.
func (r Reg) Inject(state, value uint) uint {
	for i, x := range r {
		if (state&(1<<uint(x)) != 0) != (value&(1<<uint(i)) != 0) {
			state ^= 1 << uint(x)
		}
	}
	return state
}
