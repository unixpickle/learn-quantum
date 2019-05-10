package main

import "github.com/unixpickle/learn-quantum/quantum"

type constAdder struct {
	Value  uint
	Target quantum.Reg
	Carry  *int
	Invert bool
}

func (c *constAdder) String() string {
	return "ConstAdder"
}

func (c *constAdder) Apply(comp quantum.Computer) {
	sim := comp.(*quantum.Simulation)
	newPhases := make([]complex128, len(sim.Phases))
	for i, ph := range sim.Phases {
		input := c.Target.Extract(uint(i))
		var sum uint
		if c.Invert {
			sum = input - c.Value
		} else {
			sum = input + c.Value
		}
		carryBit := sum & (1 << uint(len(c.Target)))
		sum &= (1 << uint(len(c.Target))) - 1
		newState := uint(i)
		newState = c.Target.Inject(newState, sum)
		if c.Carry != nil && carryBit != 0 {
			newState ^= 1 << uint(*c.Carry)
		}
		newPhases[newState] = ph
	}
	sim.Phases = newPhases
}

func (c *constAdder) Inverse() quantum.Gate {
	c1 := *c
	c1.Invert = !c1.Invert
	return &c1
}
