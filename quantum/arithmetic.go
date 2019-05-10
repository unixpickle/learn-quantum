package quantum

// Add performs basic addition in base 2.
// The value stored in source is added to the value stored
// in target. If the carry argument is non-nil, it is used
// as the qubit to be flipped if the addition wraps.
// The source and target must be the same number of bits.
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
