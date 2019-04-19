package quantum

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func ExampleState() {
	// Create s1 = |+> |->
	s1 := NewState(2)
	s1 = s1.Not(1)
	s1 = s1.Hadamard(0)
	s1 = s1.Hadamard(1)

	// Create s2 = |-> |->
	s2 := NewState(2)
	s2 = s2.Not(0)
	s2 = s2.Not(1)
	s2 = s2.Hadamard(0)
	s2 = s2.Hadamard(1)

	// Apply CNot on s1.
	s1 = s1.CNot(0, 1)

	fmt.Println(s1)
	fmt.Println(s2)

	// Output:
	// 0.5|00> + -0.5|10> + -0.5|01> + 0.5|11>
	// 0.5|00> + -0.5|10> + -0.5|01> + 0.5|11>
}

func TestStateSample(t *testing.T) {
	s := NewState(2)
	s = s.Hadamard(0)
	s = s.CNot(0, 1)
	counts := map[uint64]int{}
	for i := 0; i < 100000; i++ {
		counts[s.Sample()]++
	}
	if counts[1] != 0 || counts[2] != 0 {
		fmt.Println("unexpected result")
	}
	// Stddev should be on the order of ~100.
	if math.Abs(float64(counts[0]-counts[3])) > 1000 {
		fmt.Println("incorrect sample counts, delta is", counts[0]-counts[3])
	}
}

func TestInverses(t *testing.T) {
	s := RandomState(8)
	original := s.Copy()
	makeNot := func(idx int) func() {
		return func() {
			s = s.Not(idx)
		}
	}
	makeHadamard := func(idx int) func() {
		return func() {
			s = s.Hadamard(idx)
		}
	}
	makeCNot := func(control, target int) func() {
		return func() {
			s = s.CNot(control, target)
		}
	}
	ops := []func(){}
	for i := 0; i < 1000; i++ {
		x := rand.Intn(3)
		if x == 0 {
			ops = append(ops, makeNot(rand.Intn(8)))
		} else if x == 1 {
			ops = append(ops, makeHadamard(rand.Intn(8)))
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
	if s.ApproxEqual(original, 0) {
		t.Error("states should be nowhere close to equal")
	}
	for i := len(ops) - 1; i >= 0; i-- {
		ops[i]()
	}
	if !s.ApproxEqual(original, 0) {
		t.Error("states not equal")
	}
}
