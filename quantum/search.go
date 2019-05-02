package quantum

// A CircuitGen generates exhaustive lists of circuits.
//
// It is not safe to call methods on a CircuitGen from
// multiple Goroutines concurrently.
type CircuitGen struct {
	numBits        int
	basis          []Gate
	hasher         CircuitHasher
	cache          [][]Circuit
	cacheRemaining int
}

// NewCircuitGen creates a new circuit generator.
//
// The maxCache argument specifies how many circuits the
// generator may store in memory to increase efficiency
// and reduce duplicates.
func NewCircuitGen(numBits int, basis []Gate, maxCache int) *CircuitGen {
	return &CircuitGen{
		numBits:        numBits,
		basis:          basis,
		hasher:         NewCircuitHasher(numBits),
		cache:          [][]Circuit{[]Circuit{Circuit{}}},
		cacheRemaining: maxCache,
	}
}

// GenerateSlice uses the circuit cache to provide an
// in-memory list of all the circuits of a given size.
//
// If all the circuits do not fit into memory, this
// returns nil.
func (c *CircuitGen) GenerateSlice(numGates int) []Circuit {
	for len(c.cache) <= numGates && c.cacheRemaining > 0 {
		c.extendCache()
	}
	if numGates >= len(c.cache) {
		return nil
	}
	return c.cache[numGates]
}

// Generate generates a (possibly redundant) sequence of
// circuits of a given size.
func (c *CircuitGen) Generate(numGates int) (<-chan Circuit, int) {
	for len(c.cache) <= numGates && c.cacheRemaining > 0 {
		c.extendCache()
	}

	ch := make(chan Circuit, 10)

	if numGates < len(c.cache) {
		go func() {
			defer close(ch)
			for _, circ := range c.cache[numGates] {
				ch <- circ
			}
		}()
		return ch, len(c.cache[numGates])
	}

	subCh, subCount := c.Generate(numGates - len(c.cache) + 1)
	go func() {
		defer close(ch)
		for subCirc := range subCh {
			for _, circ := range c.cache[len(c.cache)-1] {
				ch <- append(append(Circuit{}, subCirc...), circ...)
			}
		}
	}()

	return ch, subCount * len(c.cache[len(c.cache)-1])
}

func (c *CircuitGen) extendCache() bool {
	var next []Circuit
	found := map[CircuitHash]bool{}
	for _, prevCirc := range c.cache[len(c.cache)-1] {
		for _, gate := range c.basis {
			circ := append(Circuit{gate}, prevCirc...)
			hash := c.hasher.Hash(circ)
			if !found[hash] {
				found[hash] = true
				next = append(next, circ)
				c.cacheRemaining -= 1
				if c.cacheRemaining == 0 {
					return false
				}
			}
		}
	}
	c.cache = append(c.cache, next)
	return true
}
