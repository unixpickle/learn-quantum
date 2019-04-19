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
	Sample() []bool
	Unitary(target int, m11, m12, m21, m22 complex128)
	CNot(control, target int)
}

// A Simulation is a classical simulation of a quantum
// computer.
type Simulation struct {
	numBits int
	phases  []complex128
}

// Create a new Simulation with all qubits set to 0.
func NewSimulation(numBits int) *Simulation {
	s := &Simulation{
		numBits: numBits,
		phases:  make([]complex128, 1<<uint(numBits)),
	}
	s.phases[0] = 1
	return s
}

// Create a new Simulation in a random state.
func RandomSimulation(numBits int) *Simulation {
	s := NewSimulation(numBits)
	mag := 0.0
	for i := range s.phases {
		s.phases[i] = complex(rand.NormFloat64(), rand.NormFloat64())
		mag += math.Pow(cmplx.Abs(s.phases[i]), 2)
	}
	scale := complex(1/math.Sqrt(mag), 0)
	for i := range s.phases {
		s.phases[i] *= scale
	}
	return s
}

func (s *Simulation) NumBits() int {
	return s.numBits
}

func (s *Simulation) Sample() []bool {
	x := s.sampleIndex()
	res := make([]bool, s.numBits)
	for i := range res {
		res[i] = x&(1<<uint(i)) != 0
	}
	return res
}

func (s *Simulation) sampleIndex() int {
	x := rand.Float64()
	for i, phase := range s.phases {
		v := real(phase)*real(phase) + imag(phase)*imag(phase)
		x -= v
		if x <= 0 {
			return i
		}
	}
	return len(s.phases) - 1
}

func (s *Simulation) Unitary(target int, m11, m12, m21, m22 complex128) {
	if target < 0 || target >= s.numBits {
		panic("bit index out of range")
	}
	res := make([]complex128, len(s.phases))
	for i := range s.phases {
		if i&(1<<uint(target)) != 0 {
			continue
		}
		other := i | (1 << uint(target))
		p0 := s.phases[i]
		p1 := s.phases[other]
		res[i] += m11*p0 + m12*p1
		res[other] += m21*p0 + m22*p1
	}
	s.phases = res
}

func (s *Simulation) CNot(control, target int) {
	if control < 0 || control >= s.numBits || target < 0 || target >= s.numBits {
		panic("bit index out of range")
	}
	res := make([]complex128, len(s.phases))
	for i, phase := range s.phases {
		b1 := (i & (1 << uint(control))) >> uint(control)
		res[i^(b1<<uint(target))] = phase
	}
	s.phases = res
}

func (s *Simulation) Copy() *Simulation {
	res := &Simulation{
		numBits: s.numBits,
		phases:  make([]complex128, len(s.phases)),
	}
	for i, phase := range s.phases {
		res.phases[i] = phase
	}
	return res
}

func (s *Simulation) ApproxEqual(s1 *Simulation, tol float64) bool {
	for i, phase := range s.phases {
		if cmplx.Abs(phase-s1.phases[i]) > tol {
			return false
		}
	}
	return true
}

func (s *Simulation) String() string {
	pieces := []string{}
	for i, phase := range s.phases {
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
