package quantum

import (
	"math/rand"
	"testing"
)

func TestAdd(t *testing.T) {
	for _, carry := range []bool{false, true} {
		name := "NoCarry"
		if carry {
			name = "Carry"
		}
		t.Run(name, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				s1 := NewSimulation(11)
				s2 := s1.Copy()
				bits := rand.Perm(s1.NumBits())
				source := Reg(bits[:5])
				target := Reg(bits[5:10])
				var carryField *int
				if carry {
					carryField = &bits[10]
				}
				Add(s1, source, target, carryField)
				simulatedAdd(s2, source, target, carryField)
				if !s1.ApproxEqual(s2, 1e-8) {
					t.Error("bad results")
				}
			}
		})
	}
}

func simulatedAdd(sim *Simulation, source, target Reg, carry *int) {
	newPhases := make([]complex128, len(sim.Phases))
	for i, ph := range sim.Phases {
		n1 := source.Extract(uint(i))
		n2 := target.Extract(uint(i))
		sum := n1 + n2
		carryBit := sum & (1 << uint(len(source)))
		sum &= (1 << uint(len(source))) - 1
		newState := uint(i)
		newState = target.Inject(newState, sum)
		if carry != nil && carryBit == 1 {
			newState ^= 1 << uint(*carry)
		}
		newPhases[newState] = ph
	}
	sim.Phases = newPhases
}
