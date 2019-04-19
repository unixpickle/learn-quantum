package quantum

import "fmt"

func ExampleCNotSource() {
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
