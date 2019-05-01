package quantum

import (
	"crypto/md5"
	"math"
)

const valueScale = 30000

type CircuitHash [md5.Size]byte

type CircuitHasher interface {
	NumBits() int
	Hash(g Gate) CircuitHash

	// Prefix creates a hasher that applies the prefix g
	// before every gate it hashes.
	Prefix(g Gate) CircuitHasher
}

type circuitHasher struct {
	startState *Simulation
}

// NewCircuitHasher creates a random hash function for
// quantum circuits. Each call may yield different hash
// functions, depending on the global random seed.
func NewCircuitHasher(numBits int) CircuitHasher {
	start := RandomSimulation(numBits)

	// Move coefficients to the boundaries where we
	// discretize, that way rounding errors are unlikely
	// to affect the final hash.
	// This slightly unnormalizes the simulation, but as
	// long as we never sample, that should not be an
	// issue.
	for i, phase := range start.Phases {
		r := real(phase)
		im := imag(phase)
		start.Phases[i] = complex(
			float64(int16(r*valueScale))/valueScale,
			float64(int16(im*valueScale))/valueScale,
		)
	}

	return &circuitHasher{startState: start}
}

func (c *circuitHasher) NumBits() int {
	return c.startState.NumBits()
}

func (c *circuitHasher) Hash(g Gate) CircuitHash {
	s := c.startState.Copy()
	g.Apply(s)
	data := make([]byte, 0, len(s.Phases)*4)
	for _, phase := range s.Phases {
		r := uint16(int16(math.Round(valueScale * real(phase))))
		im := uint16(int16(math.Round(valueScale * imag(phase))))
		data = append(data, byte(r>>8), byte(r), byte(im>>8), byte(im))
	}
	return md5.Sum(data)
}

func (c *circuitHasher) Prefix(g Gate) CircuitHasher {
	s := c.startState.Copy()
	g.Apply(s)
	return &circuitHasher{startState: s}
}
