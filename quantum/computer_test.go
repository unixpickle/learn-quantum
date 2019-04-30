package quantum

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"testing"
)

func ExampleSimulation() {
	// Create s1 = |+> |->
	s1 := NewSimulation(2)
	X(s1, 1)
	H(s1, 0)
	H(s1, 1)

	// Create s2 = |-> |->
	s2 := NewSimulation(2)
	X(s2, 0)
	X(s2, 1)
	H(s2, 0)
	H(s2, 1)

	// Apply CNot on s1.
	s1.CNot(0, 1)

	fmt.Println(s1)
	fmt.Println(s2)

	// Output:
	// 0.5|00> + -0.5|10> + -0.5|01> + 0.5|11>
	// 0.5|00> + -0.5|10> + -0.5|01> + 0.5|11>
}

func TestSimulationSample(t *testing.T) {
	s := NewSimulation(2)
	H(s, 0)
	s.CNot(0, 1)
	counts := map[int]int{}
	for i := 0; i < 100000; i++ {
		b := s.Copy().Sample()
		n := 0
		if b[0] {
			n |= 1
		}
		if b[1] {
			n |= 2
		}
		counts[n]++
	}
	if counts[1] != 0 || counts[2] != 0 {
		t.Errorf("unexpected result")
	}
	// Stddev should be on the order of ~100.
	if math.Abs(float64(counts[0]-counts[3])) > 1000 {
		t.Errorf("incorrect sample counts, delta is %d", counts[0]-counts[3])
	}
}

func TestInverses(t *testing.T) {
	s := RandomSimulation(8)
	original := s.Copy()
	makeX := func(idx int) func() {
		return func() {
			X(s, idx)
		}
	}
	makeY := func(idx int) func() {
		return func() {
			Y(s, idx)
		}
	}
	makeZ := func(idx int) func() {
		return func() {
			Z(s, idx)
		}
	}
	makeH := func(idx int) func() {
		return func() {
			H(s, idx)
		}
	}
	makeCNot := func(control, target int) func() {
		return func() {
			s.CNot(control, target)
		}
	}
	ops := []func(){}
	for i := 0; i < 1000; i++ {
		x := rand.Intn(5)
		if x == 0 {
			ops = append(ops, makeX(rand.Intn(8)))
		} else if x == 1 {
			ops = append(ops, makeH(rand.Intn(8)))
		} else if x == 2 {
			ops = append(ops, makeY(rand.Intn(8)))
		} else if x == 3 {
			ops = append(ops, makeZ(rand.Intn(8)))
		} else {
			n1 := rand.Intn(8)
			n2 := n1
			for n2 == n1 {
				n2 = rand.Intn(8)
			}
			ops = append(ops, makeCNot(n1, n2))
		}
	}
	for _, op := range ops {
		op()
	}
	if s.ApproxEqual(original, epsilon) {
		t.Error("states should be nowhere close to equal")
	}
	for i := len(ops) - 1; i >= 0; i-- {
		ops[i]()
	}
	if !s.ApproxEqual(original, epsilon) {
		t.Error("states not equal")
	}
}

func TestCCNot(t *testing.T) {
	for i := 0; i < 8; i++ {
		s := NewSimulationBits(3, uint(i))
		output := i ^ (((i & 1) << 2) & ((i & 2) << 1))
		CCNot(s, 0, 1, 2)
		if cmplx.Abs(s.Phases[output]-1) > 1e-8 {
			t.Error("incorrect result for", i)
		}
	}
}

func TestSimulationInterference(t *testing.T) {
	for i := 0; i < 4; i++ {
		s := NewSimulationBits(2, uint(i))
		H(s, 0)
		H(s, 1)
		s.CNot(0, 1)
		H(s, 0)
		H(s, 1)
		output := i ^ ((i & 2) >> 1)
		if cmplx.Abs(s.Phases[output]-1) > 1e-8 {
			t.Error("unexpected output for", i)
		}
	}
}
