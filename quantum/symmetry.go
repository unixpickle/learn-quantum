package quantum

import (
	"crypto/md5"
	"math"
	"math/cmplx"
	"math/rand"
)

type symHasher struct {
	startState *Simulation
}

// NewSymHasher creates a CircuitHasher that yields equal
// hashes for circuits that are equivalent up to a
// permutation of the qubits.
func NewSymHasher(numBits int) CircuitHasher {
	coefficients := make([]complex128, numBits+1)
	for i := range coefficients {
		coefficients[i] = complex(rand.Float64(), rand.Float64())
	}

	state := NewSimulation(numBits)
	for i := range state.Phases {
		state.Phases[i] = coefficients[countOnes(numBits, i)]
	}

	var mag float64
	for _, p := range state.Phases {
		mag += math.Pow(cmplx.Abs(p), 2)
	}
	normalizer := complex(1/math.Sqrt(mag), 0)
	for i := range state.Phases {
		state.Phases[i] *= normalizer
	}

	roundHashStart(state)
	return &symHasher{startState: state}
}

func (s *symHasher) NumBits() int {
	return s.startState.NumBits()
}

func (s *symHasher) Hash(g Gate) CircuitHash {
	sim := s.startState.Copy()
	g.Apply(sim)

	bitSums := make([]uint64, sim.NumBits()+1)
	for i, phase := range sim.Phases {
		r := uint64(uint32(int32(math.Round(valueScale * real(phase)))))
		im := uint64(uint32(int32(math.Round(valueScale * imag(phase)))))

		enc := r | (im << 32)
		bitSums[countOnes(sim.NumBits(), i)] += enc
	}

	data := make([]byte, len(bitSums)*8)
	for i, x := range bitSums {
		for j := 0; j < 8; j++ {
			data[i*8+j] = byte(x >> uint(j*8))
		}
	}

	return md5.Sum(data)
}

func (s *symHasher) Prefix(g Gate) CircuitHasher {
	sim := s.startState.Copy()
	g.Apply(sim)
	return &symHasher{startState: sim}
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

func countOnes(numBits, n int) int {
	var numOnes int
	for i := 0; i < numBits; i++ {
		if n&(1<<uint(i)) != 0 {
			numOnes += 1
		}
	}
	return numOnes
}
