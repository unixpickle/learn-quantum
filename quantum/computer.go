package quantum

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"strconv"
	"strings"
)

const epsilon = 1e-8

// A Computer is a generic quantum computer.
type Computer interface {
	NumBits() int
	Measure(bitIdx int) bool
	Unitary(target int, m11, m12, m21, m22 complex128)
	CNot(control, target int)
}

// A Simulation is a classical simulation of a quantum
// computer.
type Simulation struct {
	numBits int
	Phases  []complex128
}

// Create a new Simulation with all qubits set to 0.
func NewSimulation(numBits int) *Simulation {
	return NewSimulationBits(numBits, 0)
}

// Create a new Simulation with a given bit-string.
func NewSimulationBits(numBits int, value uint) *Simulation {
	s := &Simulation{
		numBits: numBits,
		Phases:  make([]complex128, 1<<uint(numBits)),
	}
	s.Phases[int(value)] = 1
	return s
}

// Create a new Simulation in a random state.
func RandomSimulation(numBits int) *Simulation {
	s := NewSimulation(numBits)
	mag := 0.0
	for i := range s.Phases {
		s.Phases[i] = complex(rand.NormFloat64(), rand.NormFloat64())
		mag += math.Pow(cmplx.Abs(s.Phases[i]), 2)
	}
	scale := complex(1/math.Sqrt(mag), 0)
	for i := range s.Phases {
		s.Phases[i] *= scale
	}
	return s
}

func (s *Simulation) NumBits() int {
	return s.numBits
}

func (s *Simulation) Measure(bitIdx int) bool {
	var zeroProb float64
	var oneProb float64
	for i, ph := range s.Phases {
		prob := math.Pow(cmplx.Abs(ph), 2)
		if i&(1<<uint(bitIdx)) != 0 {
			oneProb += prob
		} else {
			zeroProb += prob
		}
	}
	isOne := rand.Float64() > zeroProb
	var scale float64
	if isOne {
		scale = 1 / math.Sqrt(oneProb)
	} else {
		scale = 1 / math.Sqrt(zeroProb)
	}
	for i := range s.Phases {
		if (i&(1<<uint(bitIdx)) != 0) != isOne {
			s.Phases[i] = 0
		} else {
			s.Phases[i] *= complex(scale, 0)
		}
	}
	return isOne
}

func (s *Simulation) Unitary(target int, m11, m12, m21, m22 complex128) {
	if target < 0 || target >= s.numBits {
		panic("bit index out of range")
	}
	res := make([]complex128, len(s.Phases))
	for i := range s.Phases {
		if i&(1<<uint(target)) != 0 {
			continue
		}
		other := i | (1 << uint(target))
		p0 := s.Phases[i]
		p1 := s.Phases[other]
		res[i] += m11*p0 + m12*p1
		res[other] += m21*p0 + m22*p1
	}
	s.Phases = res
}

func (s *Simulation) CNot(control, target int) {
	if control < 0 || control >= s.numBits || target < 0 || target >= s.numBits {
		panic("bit index out of range")
	}
	res := make([]complex128, len(s.Phases))
	for i, phase := range s.Phases {
		b1 := (i & (1 << uint(control))) >> uint(control)
		res[i^(b1<<uint(target))] = phase
	}
	s.Phases = res
}

func (s *Simulation) Phase(value []bool) complex128 {
	var idx int
	for i, x := range value {
		if x {
			idx |= 1 << uint(i)
		}
	}
	return s.Phases[idx]
}

func (s *Simulation) Copy() *Simulation {
	res := &Simulation{
		numBits: s.numBits,
		Phases:  make([]complex128, len(s.Phases)),
	}
	for i, phase := range s.Phases {
		res.Phases[i] = phase
	}
	return res
}

func (s *Simulation) ApproxEqual(s1 *Simulation, tol float64) bool {
	for i, phase := range s.Phases {
		if cmplx.Abs(phase-s1.Phases[i]) > tol {
			return false
		}
	}
	return true
}

func (s *Simulation) String() string {
	pieces := []string{}
	for i, phase := range s.Phases {
		if cmplx.Abs(phase) < epsilon {
			continue
		}
		var coeff string
		if math.Abs(imag(phase)) < epsilon {
			coeff = formatFloat(real(phase))
		} else if math.Abs(real(phase)) < epsilon {
			coeff = formatFloat(imag(phase)) + "i"
		} else {
			if imag(phase) > 0 {
				coeff = fmt.Sprintf("(%s+%si)", formatFloat(real(phase)), formatFloat(imag(phase)))
			} else {
				coeff = fmt.Sprintf("(%s-%si)", formatFloat(real(phase)), formatFloat(-imag(phase)))
			}
		}
		pieces = append(pieces, coeff+s.classicalString(i))
	}
	return strings.Join(pieces, " + ")
}

func (s *Simulation) Sample() []bool {
	var res []bool
	for i := 0; i < s.numBits; i++ {
		res = append(res, s.Measure(i))
	}
	return res
}

// ControlGate runs the gate g on states where control is
// not set.
// The gate g must not modify the control bit.
func (s *Simulation) ControlGate(control int, g Gate) {
	s1 := s.Copy()
	for i := range s.Phases {
		if i&1 == 0 {
			s1.Phases[i] = 0
		} else {
			s.Phases[i] = 0
		}
	}
	g.Apply(s1)
	for i, phase := range s1.Phases {
		if phase != 0 && s.Phases[i] != 0 {
			panic("gate must not modify control bit")
		}
		s.Phases[i] += phase
	}
}

func (s *Simulation) classicalString(i int) string {
	res := ""
	for j := 0; j < s.numBits; j++ {
		res += strconv.Itoa((i & (1 << uint(j))) >> uint(j))
	}
	return "|" + res + ">"
}

func formatFloat(f float64) string {
	res := fmt.Sprintf("%f", f)
	for strings.Contains(res, ".") && res[len(res)-1] == '0' {
		res = res[:len(res)-1]
	}
	if res[len(res)-1] == '.' {
		return res[:len(res)-1]
	}
	return res
}
