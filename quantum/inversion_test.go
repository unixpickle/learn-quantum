package quantum

import "testing"

func TestInvert(t *testing.T) {
	circuit := Circuit{
		&HGate{Bit: 0},
		&HGate{Bit: 3},
		&SqrtNotGate{Bit: 0},
		&CSqrtNotGate{Control: 2, Target: 1},
		&TGate{Bit: 0},
		&CCNotGate{Control1: 2, Control2: 0, Target: 3},
	}
	s1 := RandomSimulation(4)
	s2 := s1.Copy()
	Invert(s1, circuit.Apply)
	circuit.Inverse().Apply(s2)
	if !s1.ApproxEqual(s2, 1e-8) {
		t.Error("invalid result")
	}
}

func TestConj(t *testing.T) {
	a := Circuit{
		&HGate{Bit: 0},
		&HGate{Bit: 3},
		&SqrtNotGate{Bit: 0},
		&CSqrtNotGate{Control: 2, Target: 1},
		&TGate{Bit: 0},
		&CCNotGate{Control1: 2, Control2: 0, Target: 3},
	}
	b := Circuit{
		&HGate{Bit: 1},
		&HGate{Bit: 2},
		&SqrtNotGate{Bit: 3},
		&CSqrtNotGate{Control: 0, Target: 1},
		&TGate{Bit: 2},
		&CCNotGate{Control1: 3, Control2: 0, Target: 1},
	}
	s1 := RandomSimulation(4)
	s2 := s1.Copy()
	Conj(s1, a.Apply, b.Apply)
	a.Apply(s2)
	b.Apply(s2)
	a.Inverse().Apply(s2)
	if !s1.ApproxEqual(s2, 1e-8) {
		t.Error("invalid result")
	}
}
