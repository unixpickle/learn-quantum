package quantum

// Add performs basic addition in base 2.
// The value stored in source is added to the value stored
// in target. If the carry argument is non-nil, it is used
// as the qubit to be flipped if the addition wraps.
// The source and target must be the same number of bits.
// The source and target registers are stored lowest-bit
// first.
func Add(c Computer, source, target Reg, carry *int) {
	var carryReg Reg
	if carry != nil {
		carryReg = Reg{*carry}
	}
	if source.Overlaps(target) || source.Overlaps(carryReg) || target.Overlaps(carryReg) ||
		len(source) != len(target) || !source.Valid() || !target.Valid() {
		panic("invalid arguments")
	}

	if len(source) == 1 {
		if carry != nil {
			CCNot(c, source[0], target[0], *carry)
		}
		c.CNot(source[0], target[0])
		return
	}

	// Implementation of https://arxiv.org/abs/0910.2530

	// Step 1
	for i := 1; i < len(source); i++ {
		c.CNot(source[i], target[i])
	}

	// Step 2
	if carry != nil {
		c.CNot(source[len(source)-1], *carry)
	}
	for i := len(source) - 2; i > 0; i-- {
		c.CNot(source[i], source[i+1])
	}

	// Step 3
	for i := 0; i < len(source)-1; i++ {
		CCNot(c, source[i], target[i], source[i+1])
	}
	if carry != nil {
		CCNot(c, source[len(source)-1], target[len(target)-1], *carry)
	}

	// Step 4
	for i := len(source) - 1; i > 0; i-- {
		c.CNot(source[i], target[i])
		CCNot(c, source[i-1], target[i-1], source[i])
	}

	// Step 5
	for i := 1; i < len(source)-1; i++ {
		c.CNot(source[i], source[i+1])
	}

	// Step 6
	for i := 0; i < len(source); i++ {
		c.CNot(source[i], target[i])
	}
}

// Sub performs the inverse of Add.
func Sub(c Computer, source, target Reg, carry *int) {
	var carryReg Reg
	if carry != nil {
		carryReg = Reg{*carry}
	}
	if source.Overlaps(target) || source.Overlaps(carryReg) || target.Overlaps(carryReg) ||
		len(source) != len(target) || !source.Valid() || !target.Valid() {
		panic("invalid arguments")
	}

	if len(source) == 1 {
		c.CNot(source[0], target[0])
		if carry != nil {
			CCNot(c, source[0], target[0], *carry)
		}
		return
	}

	// Step 6
	for i := len(source) - 1; i >= 0; i-- {
		c.CNot(source[i], target[i])
	}

	// Step 5
	for i := len(source) - 2; i > 0; i-- {
		c.CNot(source[i], source[i+1])
	}

	// Step 4
	for i := 1; i < len(source); i++ {
		CCNot(c, source[i-1], target[i-1], source[i])
		c.CNot(source[i], target[i])
	}

	// Step 3
	if carry != nil {
		CCNot(c, source[len(source)-1], target[len(target)-1], *carry)
	}
	for i := len(source) - 2; i >= 0; i-- {
		CCNot(c, source[i], target[i], source[i+1])
	}

	// Step 2
	for i := 1; i < len(source)-1; i++ {
		c.CNot(source[i], source[i+1])
	}
	if carry != nil {
		c.CNot(source[len(source)-1], *carry)
	}

	// Step 1
	for i := len(source) - 1; i > 0; i-- {
		c.CNot(source[i], target[i])
	}
}

// Lt flips a bit if a < b in unsigned arithmetic.
func Lt(c Computer, a, b Reg, target int) {
	Sub(c, b, a, &target)
	Add(c, b, a, nil)
}

// ModAdd performs modular addition.
//
// Behavior is undefined if the inputs are greater than or
// equal to the modulus.
//
// The working qubits must start as zeros and will end as
// zeros.
func ModAdd(c Computer, source, target, modulus Reg, working int) {
	workingReg := Reg{working}
	if len(source) != len(target) || len(source) != len(modulus) || !source.Valid() ||
		!target.Valid() || !modulus.Valid() || source.Overlaps(target) ||
		target.Overlaps(modulus) || source.Overlaps(workingReg) || target.Overlaps(workingReg) ||
		modulus.Overlaps(workingReg) {
		panic("invalid inputs")
	}

	// Extend the target with one working bit, then add
	// the source to it.
	Add(c, source, target, &working)
	Sub(c, modulus, target, &working)

	// If the carry bit is set, it means we wrapped around
	// when subtracting the modulus.
	Cond(c, working, func(c Computer) {
		Add(c, modulus, target, nil)
	})

	// Reverse working1 if it was set, using the
	// assumption that source and target started out both
	// less than the modulus.
	X(c, working)
	Lt(c, target, source, working)
}

// ModSub is the inverse of ModAdd.
func ModSub(c Computer, source, target, modulus Reg, working int) {
	workingReg := Reg{working}
	if len(source) != len(target) || len(source) != len(modulus) || !source.Valid() ||
		!target.Valid() || !modulus.Valid() || source.Overlaps(target) ||
		target.Overlaps(modulus) || source.Overlaps(workingReg) || target.Overlaps(workingReg) ||
		modulus.Overlaps(workingReg) {
		panic("invalid inputs")
	}

	Lt(c, target, source, working)
	X(c, working)
	Cond(c, working, func(c Computer) {
		Sub(c, modulus, target, nil)
	})
	Add(c, modulus, target, &working)
	Sub(c, source, target, &working)
}
