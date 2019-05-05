package quantum

import (
	"math/rand"
	"testing"
)

const symTestNumBits = 10

func TestSymHasher(t *testing.T) {
	hasher := NewSymHasher(symTestNumBits)

	// Check that the hash doesn't identify different
	// circuits.
	// if hasher.Hash(&XGate{Bit: 0}) == hasher.Hash(Circuit{}) {
	// 	t.Error("hash collision")
	// }

	// if hasher.Hash(&HGate{Bit: 0}) == hasher.Hash(Circuit{}) {
	// 	t.Error("hash collision")
	// }

	// Check that the hash identifies permutations of the
	// same circuits.
	for seed := 0; seed < 10; seed++ {
		hash := hasher.Hash(randomizedCircuit(rand.Perm(symTestNumBits), seed))
		for i := 0; i < 10; i++ {
			if hasher.Hash(randomizedCircuit(rand.Perm(symTestNumBits), seed)) != hash {
				t.Fatal("mismatching hash")
			}
		}
	}
}

func randomizedCircuit(perm []int, seed int) Circuit {
	gen := rand.New(rand.NewSource(int64(seed)))
	var c Circuit
	for i := 0; i < 20; i++ {
		n := gen.Intn(3)
		if n == 0 {
			a := gen.Intn(symTestNumBits)
			b := gen.Intn(symTestNumBits)
			for b == a {
				b = gen.Intn(symTestNumBits)
			}
			c = append(c, &CNotGate{Control: perm[a], Target: perm[b]})
		} else if n == 1 {
			c = append(c, &HGate{Bit: perm[gen.Intn(symTestNumBits)]})
		} else if n == 2 {
			c = append(c, &TGate{Bit: perm[gen.Intn(symTestNumBits)]})
		}
	}
	return c
}
