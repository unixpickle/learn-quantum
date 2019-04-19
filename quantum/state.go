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

type State struct {
	phases  []complex128
	numBits int
}

func NewState(numBits int) *State {
	s := &State{
		phases:  make([]complex128, 1<<uint(numBits)),
		numBits: numBits,
	}
	s.phases[0] = 1
	return s
}

func RandomState(numBits int) *State {
	s := NewState(numBits)
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

func (s *State) String() string {
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

func (s *State) Copy() *State {
	// TODO: do this in a less lazy way.
	return s.Not(0).Not(0)
}

func (s *State) ApproxEqual(s1 *State, tol float64) bool {
	if tol == 0 {
		tol = epsilon
	}
	for i, phase := range s.phases {
		if cmplx.Abs(phase-s1.phases[i]) > tol {
			return false
		}
	}
	return true
}

func (s *State) Sample() uint64 {
	x := rand.Float64()
	for i, phase := range s.phases {
		v := real(phase)*real(phase) + imag(phase)*imag(phase)
		x -= v
		if x <= 0 {
			return uint64(i)
		}
	}
	return uint64(len(s.phases) - 1)
}

func (s *State) Not(bitIdx int) *State {
	if bitIdx < 0 || bitIdx >= s.numBits {
		panic("bit index out of range")
	}
	return s.Map(func(x uint64) uint64 {
		return x ^ (1 << uint(bitIdx))
	})
}

func (s *State) CNot(control, target int) *State {
	if control < 0 || control >= s.numBits || target < 0 || target >= s.numBits {
		panic("bit index out of range")
	}
	return s.Map(func(x uint64) uint64 {
		b1 := (x & (1 << uint(control))) >> uint(control)
		return x ^ (b1 << uint(target))
	})
}

func (s *State) Hadamard(bitIdx int) *State {
	if bitIdx < 0 || bitIdx >= s.numBits {
		panic("bit index out of range")
	}
	s1 := NewState(s.numBits)
	s1.phases[0] = 0
	coeff := complex(1/math.Sqrt2, 0)
	for i, phase := range s.phases {
		flip := i ^ (1 << uint(bitIdx))
		if i&(1<<uint(bitIdx)) == 0 {
			s1.phases[i] += phase * coeff
			s1.phases[flip] += phase * coeff
		} else {
			s1.phases[i] -= phase * coeff
			s1.phases[flip] += phase * coeff
		}
	}
	return s1
}

func (s *State) Map(f func(uint64) uint64) *State {
	s1 := NewState(s.numBits)
	s1.phases[0] = 0
	for i, phase := range s.phases {
		s1.phases[int(f(uint64(i)))] += phase
	}
	return s1
}

func (s *State) classicalString(i int) string {
	res := ""
	for j := 0; j < s.numBits; j++ {
		res += strconv.Itoa((i & (1 << uint(j))) >> uint(j))
	}
	return "|" + res + ">"
}

func getBit(value uint64, idx int) bool {
	return value<<uint(idx) != 0
}

func setBit(value uint64, idx int, b bool) uint64 {
	if getBit(value, idx) == b {
		return value
	}
	return value ^ (1 << uint(idx))
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
