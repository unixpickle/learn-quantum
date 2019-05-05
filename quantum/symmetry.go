package quantum

import (
	"crypto/md5"
	"math"
	"math/cmplx"
	"math/rand"

	"github.com/unixpickle/essentials"
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
		var numOnes int
		for j := 0; j < numBits; j++ {
			if i&(1<<uint(j)) != 0 {
				numOnes += 1
			}
		}
		state.Phases[i] = coefficients[numOnes]
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
	for i := range phaseEnc {
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

func invPermuteBits(perm []int, num int) int {
	var res int
	for i, target := range perm {
		if num&(1<<uint(i)) != 0 {
			res |= 1 << uint(target)
		}
	}
	return res
}
