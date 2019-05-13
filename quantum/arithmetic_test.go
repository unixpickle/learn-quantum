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
			for numBits := 1; numBits < 7; numBits++ {
				for i := 0; i < 10; i++ {
					s1 := RandomSimulation(numBits*2 + 1)
					s2 := s1.Copy()
					bits := rand.Perm(s1.NumBits())
					source := Reg(bits[:numBits])
					target := Reg(bits[numBits : numBits*2])
					var carryField *int
					if carry {
						carryField = &bits[numBits*2]
					}
					Add(s1, source, target, carryField)
					simulatedAdd(s2, source, target, carryField)
					if !s1.ApproxEqual(s2, 1e-8) {
						t.Error("bad results", numBits)
					}
				}
			}
		})
	}
}

func TestSub(t *testing.T) {
	for _, carry := range []bool{false, true} {
		name := "NoCarry"
		if carry {
			name = "Carry"
		}
		t.Run(name, func(t *testing.T) {
			for numBits := 1; numBits < 7; numBits++ {
				for i := 0; i < 10; i++ {
					s1 := RandomSimulation(numBits*2 + 1)
					s2 := s1.Copy()
					bits := rand.Perm(s1.NumBits())
					source := Reg(bits[:numBits])
					target := Reg(bits[numBits : numBits*2])
					var carryField *int
					if carry {
						carryField = &bits[numBits*2]
					}
					Add(s1, source, target, carryField)
					Sub(s1, source, target, carryField)
					if !s1.ApproxEqual(s2, 1e-8) {
						t.Error("bad results", numBits)
					}
				}
			}
		})
	}
}

func TestLt(t *testing.T) {
	for numBits := 1; numBits < 7; numBits++ {
		for i := 0; i < 10; i++ {
			s1 := RandomSimulation(numBits*2 + 1)
			s2 := s1.Copy()
			bits := rand.Perm(s1.NumBits())
			a := Reg(bits[:numBits])
			b := Reg(bits[numBits : numBits*2])
			target := bits[numBits*2]
			Lt(s1, a, b, target)
			simulatedLt(s2, a, b, target)
			if !s1.ApproxEqual(s2, 1e-8) {
				t.Error("bad results", numBits)
			}
		}
	}
}

func TestModAdd(t *testing.T) {
	for numBits := 1; numBits < 5; numBits++ {
		for i := 0; i < 10; i++ {
			s1 := RandomSimulation(numBits*3 + 2)
			bits := rand.Perm(s1.NumBits())
			source := Reg(bits[:numBits])
			target := Reg(bits[numBits : numBits*2])
			modulus := Reg(bits[numBits*2 : numBits*3])
			working1, working2 := bits[numBits*3], bits[numBits*3+1]
			working := Reg{working1, working2}
			for i := range s1.Phases {
				if source.Extract(uint(i)) >= modulus.Extract(uint(i)) ||
					target.Extract(uint(i)) >= modulus.Extract(uint(i)) ||
					modulus.Extract(uint(i)) == 0 || working.Extract(uint(i)) != 0 {
					s1.Phases[i] = 0
				}
			}
			s2 := s1.Copy()
			ModAdd(s1, source, target, modulus, working1, working2)
			simulatedModAdd(s2, source, target, modulus)
			if !s1.ApproxEqual(s2, 1e-8) {
				t.Error("bad results", numBits)
			}
		}
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
		if carry != nil && carryBit != 0 {
			newState ^= 1 << uint(*carry)
		}
		newPhases[newState] = ph
	}
	sim.Phases = newPhases
}

func simulatedLt(sim *Simulation, a, b Reg, target int) {
	newPhases := make([]complex128, len(sim.Phases))
	for i, ph := range sim.Phases {
		n1 := a.Extract(uint(i))
		n2 := b.Extract(uint(i))
		newState := uint(i)
		if n1 < n2 {
			newState ^= 1 << uint(target)
		}
		newPhases[newState] = ph
	}
	sim.Phases = newPhases
}

func simulatedModAdd(sim *Simulation, source, target, modulus Reg) {
	newPhases := make([]complex128, len(sim.Phases))
	for i, ph := range sim.Phases {
		if ph == 0 {
			continue
		}
		n1 := source.Extract(uint(i))
		n2 := target.Extract(uint(i))
		m := modulus.Extract(uint(i))
		sum := (n1 + n2) % m
		newState := uint(i)
		newState = target.Inject(newState, sum)
		newPhases[newState] = ph
	}
	sim.Phases = newPhases
}
