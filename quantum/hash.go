package quantum

import (
	"crypto/md5"
	"math"
	"math/cmplx"
	"math/rand"
)

const valueScale = (1 << 30)

type CircuitHash [md5.Size]byte

type CircuitHasher interface {
	NumBits() int
	Hash(g Applier) CircuitHash

	// Prefix creates a hasher that applies the prefix g
	// before every gate it hashes.
	Prefix(g Applier) CircuitHasher
}

type circuitHasher struct {
	startState *Simulation
}

// NewCircuitHasher creates a random hash function for
// quantum circuits. Each call may yield different hash
// functions, depending on the global random seed.
func NewCircuitHasher(numBits int) CircuitHasher {
	start := RandomSimulation(numBits)
	return &circuitHasher{startState: start}
}

func (c *circuitHasher) NumBits() int {
	return c.startState.NumBits()
}

func (c *circuitHasher) Hash(g Gate) CircuitHash {
	s := c.startState.Copy()
	g.Apply(s)
	data := make([]byte, 0, len(s.Phases)*8)
	for _, phase := range s.Phases {
		r := uint32(int32(math.Round(valueScale * real(phase))))
		im := uint32(int32(math.Round(valueScale * imag(phase))))
		data = append(data,
			byte(r>>24), byte(r>>16), byte(r>>8), byte(r),
			byte(im>>24), byte(im>>16), byte(im>>8), byte(im))
	}
	return md5.Sum(data)
}

func (c *circuitHasher) Prefix(g Gate) CircuitHasher {
	s := c.startState.Copy()
	g.Apply(s)
	return &circuitHasher{startState: s}
}

type symHasher struct {
	startState *Simulation
}

// NewSymHasher creates a CircuitHasher that yields equal
// hashes for circuits that are equivalent up to a
// permutation of the qubits.
//
// There may be false collisions, as the SymHasher is just
// a heuristic. For example, the conditional swap gate
// will yield a hash equivalent to the identity.
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
