package quantum

import (
	"fmt"
	"math"
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
