package quantum

import (
	"math/rand"
	"testing"
)

const symTestNumBits = 10

func TestSymHasher(t *testing.T) {
	hasher := NewSymHasher(symTestNumBits)

	t.Run("SimpleCollisions", func(t *testing.T) {
		if hasher.Hash(&XGate{Bit: 0}) == hasher.Hash(Circuit{}) {
			t.Error("hash collision")
		}

		if hasher.Hash(&HGate{Bit: 0}) == hasher.Hash(Circuit{}) {
			t.Error("hash collision")
		}
	})

	t.Run("CNot", func(t *testing.T) {
		hash := hasher.Hash(&CNotGate{Control: 0, Target: 1})
		for i := 0; i < symTestNumBits; i++ {
			for j := 0; j < symTestNumBits; j++ {
				if i == j {
					continue
				}
				if hasher.Hash(&CNotGate{Control: i, Target: j}) != hash {
					t.Errorf("CNot(0, 1) != CNot(%d, %d)", i, j)
				}
			}
		}
	})

	// Check that the hash identifies permutations of the
	// same circuits.
	t.Run("Random", func(t *testing.T) {
		for seed := 0; seed < 100; seed++ {
			real := randomizedCircuit(rand.Perm(symTestNumBits), seed, 20)
			hash := hasher.Hash(real)
			for i := 0; i < 10; i++ {
				if hasher.Hash(randomizedCircuit(rand.Perm(symTestNumBits), seed, 20)) != hash {
					t.Fatal("mismatching hash for", real)
				}
			}
		}
	})
}

func randomizedCircuit(perm []int, seed, size int) Circuit {
	gen := rand.New(rand.NewSource(int64(seed)))
	var c Circuit
	for i := 0; i < size; i++ {
		n := gen.Intn(4)
		if n == 0 {
			a := gen.Intn(len(perm))
			b := gen.Intn(len(perm))
			for b == a {
				b = gen.Intn(len(perm))
			}
			c = append(c, &CNotGate{Control: perm[a], Target: perm[b]})
		} else if n == 1 {
			c = append(c, &HGate{Bit: perm[gen.Intn(len(perm))]})
		} else if n == 2 {
			c = append(c, &TGate{Bit: perm[gen.Intn(len(perm))]})
		} else if n == 3 {
			var a, b, d int
			for a == b || b == d || a == d {
				a = gen.Intn(len(perm))
				b = gen.Intn(len(perm))
				d = gen.Intn(len(perm))
			}
			c = append(c, &CSwapGate{Control: perm[a], A: perm[b], B: perm[d]})
		}
	}
	return c
}

func permutations(length int) [][]int {
	if length == 0 {
		return [][]int{[]int{}}
	}
	var results [][]int
	for _, perm := range permutations(length - 1) {
		for i := 0; i <= len(perm); i++ {
			newPerm := make([]int, length)
			for j := 0; j <= len(perm); j++ {
				if j == i {
					newPerm[j] = length - 1
				} else if j < i {
					newPerm[j] = perm[j]
				} else {
					newPerm[j] = perm[j-1]
				}
			}
			results = append(results, newPerm)
		}
	}
	return results
}

type permGate struct {
	Perm []int
	Gate Gate
}

func (p *permGate) String() string {
	return ""
}

func (p *permGate) Apply(c Computer) {
	p.Gate.Apply(&MappedComputer{C: c, Mapping: p.Perm})
}

func (p *permGate) Inverse() Gate {
	return &permGate{Perm: p.Perm, Gate: p.Gate.Inverse()}
}
