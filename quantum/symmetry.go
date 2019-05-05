package quantum

import (
	"crypto/md5"
	"math"
	"math/cmplx"

	"github.com/unixpickle/essentials"
)

type symHasher struct {
	startState *Simulation
}

// NewSymHasher creates a CircuitHasher that yields equal
// hashes for circuits that are equivalent up to a
// permutation of the qubits.
func NewSymHasher(numBits int) CircuitHasher {
	s1 := RandomSimulation(numBits)

	// Create a symmetric start state by adding all of the
	// permutations of s1 together.
	sum := NewSimulation(numBits)
	sum.Phases[0] = 0
	for _, perm := range permutations(numBits) {
		for i, p := range s1.Phases {
			sum.Phases[permuteBits(perm, i)] += p
		}
	}

	var mag float64
	for _, p := range sum.Phases {
		mag += math.Pow(cmplx.Abs(p), 2)
	}
	normalizer := complex(1/math.Sqrt(mag), 0)
	for i, p := range sum.Phases {
		sum.Phases[i] = p * normalizer
	}

	roundHashStart(sum)
	return &symHasher{startState: sum}
}

func (s *symHasher) NumBits() int {
	return s.startState.NumBits()
}

func (s *symHasher) Hash(g Gate) CircuitHash {
	sim := s.startState.Copy()
	g.Apply(sim)

	phaseEnc := make([]uint64, len(sim.Phases))
	bitSums := make([]uint64, sim.NumBits())

	for i, phase := range sim.Phases {
		r := uint64(uint32(int32(math.Round(valueScale * real(phase)))))
		im := uint64(uint32(int32(math.Round(valueScale * imag(phase)))))

		enc := r | (im << 32)
		phaseEnc[i] = enc

		for b := 0; b < sim.NumBits(); b++ {
			if i&(1<<uint(b)) != 0 {
				bitSums[b] += enc
			}
		}
	}

	perm := make([]int, sim.NumBits())
	for i := range perm {
		perm[i] = i
	}
	essentials.VoodooSort(bitSums, func(i, j int) bool {
		return bitSums[i] < bitSums[j]
	}, perm)

	data := make([]byte, len(sim.Phases)*8)
	for i := 0; i < len(sim.Phases); i++ {
		n := phaseEnc[invPermuteBits(perm, i)]
		for j := 0; j < 8; j++ {
			data[i*8+j] = byte(n >> uint(j*8))
		}
	}

	return md5.Sum(data)
}

func (s *symHasher) Prefix(g Gate) CircuitHasher {
	sim := s.startState.Copy()
	g.Apply(sim)
	return &symHasher{startState: sim}
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

func permuteBits(perm []int, num int) int {
	var res int
	for i, target := range perm {
		if num&(1<<uint(target)) != 0 {
			res |= 1 << uint(i)
		}
	}
	return res
}

func invPermuteBits(perm []int, num int) int {
	var res int
	for i, target := range perm {
		if num&(1<<uint(i)) != 0 {
			res |= 1 << uint(target)
		}
	}
	return res
}
