package quantum

import (
	"fmt"
	"strconv"
	"strings"
)

// A Gate is a generic object that modifies a quantum
// computer in some primitive way.
type Gate interface {
	fmt.Stringer
	Apply(c Computer)
	Inverse() Gate
}

type Circuit []Gate

func (c Circuit) String() string {
	var parts []string
	for _, g := range c {
		parts = append(parts, g.String())
	}
	return strings.Join(parts, " ")
}

func (c Circuit) Apply(comp Computer) {
	for _, g := range c {
		g.Apply(comp)
	}
}

func (c Circuit) Inverse() Gate {
	var res Circuit
	for i := len(c) - 1; i >= 0; i-- {
		res = append(res, c[i].Inverse())
	}
	return res
}

type HGate struct {
	Bit int
}

func (h *HGate) String() string {
	return "H(" + strconv.Itoa(h.Bit) + ")"
}

func (h *HGate) Apply(c Computer) {
	H(c, h.Bit)
}

func (h *HGate) Inverse() Gate {
	return h
}

type CHGate struct {
	Control int
	Target  int
}

func (c *CHGate) String() string {
	return fmt.Sprintf("CH(%d, %d)", c.Control, c.Target)
}

func (c *CHGate) Apply(comp Computer) {
	CH(comp, c.Control, c.Target)
}

func (c *CHGate) Inverse() Gate {
	return c
}

type TGate struct {
	Bit       int
	Conjugate bool
}

func (t *TGate) String() string {
	conjStr := ""
	if t.Conjugate {
		conjStr = "*"
	}
	return "T" + conjStr + "(" + strconv.Itoa(t.Bit) + ")"
}

func (t *TGate) Apply(c Computer) {
	if t.Conjugate {
		InvT(c, t.Bit)
	} else {
		T(c, t.Bit)
	}
}

func (t *TGate) Inverse() Gate {
	return &TGate{
		Bit:       t.Bit,
		Conjugate: !t.Conjugate,
	}
}

type XGate struct {
	Bit int
}

func (x *XGate) String() string {
	return "X(" + strconv.Itoa(x.Bit) + ")"
}

func (x *XGate) Apply(c Computer) {
	X(c, x.Bit)
}

func (x *XGate) Inverse() Gate {
	return x
}

type YGate struct {
	Bit int
}

func (y *YGate) String() string {
	return "Y(" + strconv.Itoa(y.Bit) + ")"
}

func (y *YGate) Apply(c Computer) {
	Y(c, y.Bit)
}

func (y *YGate) Inverse() Gate {
	return y
}

type ZGate struct {
	Bit int
}

func (z *ZGate) String() string {
	return "Z(" + strconv.Itoa(z.Bit) + ")"
}

func (z *ZGate) Apply(c Computer) {
	Z(c, z.Bit)
}

func (z *ZGate) Inverse() Gate {
	return z
}

type CNotGate struct {
	Control int
	Target  int
}

func (c *CNotGate) String() string {
	return fmt.Sprintf("CNot(%d, %d)", c.Control, c.Target)
}

func (c *CNotGate) Apply(comp Computer) {
	comp.CNot(c.Control, c.Target)
}

func (c *CNotGate) Inverse() Gate {
	return c
}

type CCNotGate struct {
	Control1 int
	Control2 int
	Target   int
}

func (c *CCNotGate) String() string {
	return fmt.Sprintf("CCNot(%d, %d, %d)", c.Control1, c.Control2, c.Target)
}

func (c *CCNotGate) Apply(comp Computer) {
	CCNot(comp, c.Control1, c.Control2, c.Target)
}

func (c *CCNotGate) Inverse() Gate {
	return c
}

type SqrtNotGate struct {
	Bit    int
	Invert bool
}

func (s *SqrtNotGate) String() string {
	if s.Invert {
		return "SqrtNot*(" + strconv.Itoa(s.Bit) + ")"
	} else {
		return "SqrtNot(" + strconv.Itoa(s.Bit) + ")"
	}
}

func (s *SqrtNotGate) Apply(c Computer) {
	if s.Invert {
		InvSqrtNot(c, s.Bit)
	} else {
		SqrtNot(c, s.Bit)
	}
}

func (s *SqrtNotGate) Inverse() Gate {
	return &SqrtNotGate{
		Bit:    s.Bit,
		Invert: !s.Invert,
	}
}

type CSqrtNotGate struct {
	Control int
	Target  int
	Invert  bool
}

func (c *CSqrtNotGate) String() string {
	if c.Invert {
		return fmt.Sprintf("CSqrtNot*(%d, %d)", c.Control, c.Target)
	} else {
		return fmt.Sprintf("CSqrtNot(%d, %d)", c.Control, c.Target)
	}
}

func (c *CSqrtNotGate) Apply(comp Computer) {
	if c.Invert {
		InvCSqrtNot(comp, c.Control, c.Target)
	} else {
		CSqrtNot(comp, c.Control, c.Target)
	}
}

func (c *CSqrtNotGate) Inverse() Gate {
	return &CSqrtNotGate{
		Control: c.Control,
		Target:  c.Target,
		Invert:  !c.Invert,
	}
}

type SqrtTGate struct {
	Bit       int
	Conjugate bool
}

func (s *SqrtTGate) String() string {
	conjStr := ""
	if s.Conjugate {
		conjStr = "*"
	}
	return "SqrtT" + conjStr + "(" + strconv.Itoa(s.Bit) + ")"
}

func (s *SqrtTGate) Apply(c Computer) {
	if s.Conjugate {
		InvSqrtT(c, s.Bit)
	} else {
		SqrtT(c, s.Bit)
	}
}

func (s *SqrtTGate) Inverse() Gate {
	return &SqrtTGate{Bit: s.Bit, Conjugate: !s.Conjugate}
}

type SwapGate struct {
	A int
	B int
}

func (s *SwapGate) String() string {
	return fmt.Sprintf("Swap(%d, %d)", s.A, s.B)
}

func (s *SwapGate) Apply(c Computer) {
	Swap(c, s.A, s.B)
}

func (s *SwapGate) Inverse() Gate {
	return s
}

type CSwapGate struct {
	Control int
	A       int
	B       int
}

func (c *CSwapGate) String() string {
	return fmt.Sprintf("CSwap(%d, %d, %d)", c.Control, c.A, c.B)
}

func (c *CSwapGate) Apply(comp Computer) {
	CSwap(comp, c.Control, c.A, c.B)
}

func (c *CSwapGate) Inverse() Gate {
	return c
}

// FnGate is a gate that calls contained functions.
type FnGate struct {
	Forward  func(c Computer)
	Backward func(c Computer)
	Str      string
}

func (f *FnGate) String() string {
	return f.Str
}

func (f *FnGate) Apply(c Computer) {
	f.Forward(c)
}

func (f *FnGate) Inverse() Gate {
	return &FnGate{Forward: f.Backward, Backward: f.Forward, Str: "Inv(" + f.Str + ")"}
}

// A ClassicalGate applies a bitwise function to classical
// bases states.
// It can only be applied to *Simulator computers.
type ClassicalGate struct {
	F        func(b []bool) []bool
	Inverted bool
	Str      string
}

func NewClassicalGate(F func(b []bool) []bool, str string) *ClassicalGate {
	if str == "" {
		str = "ClassicalGate"
	}
	return &ClassicalGate{
		F:   F,
		Str: str,
	}
}

func (c *ClassicalGate) Apply(qc Computer) {
	s := qc.(*Simulation)
	s1 := s.Copy()
	if c.Inverted {
		for i := range s1.Phases {
			input := make([]bool, s.NumBits())
			for j := range input {
				input[j] = (i&(1<<uint(j)) != 0)
			}
			output := c.F(input)
			outIdx := 0
			for j, b := range output {
				if b {
					outIdx |= 1 << uint(j)
				}
			}
			s.Phases[i] = s1.Phases[outIdx]
		}
	} else {
		for i, phase := range s1.Phases {
			input := make([]bool, s.NumBits())
			for j := range input {
				input[j] = (i&(1<<uint(j)) != 0)
			}
			output := c.F(input)
			outIdx := 0
			for j, b := range output {
				if b {
					outIdx |= 1 << uint(j)
				}
			}
			s.Phases[outIdx] = phase
		}
	}
}

func (c *ClassicalGate) Inverse() Gate {
	return &ClassicalGate{F: c.F, Inverted: !c.Inverted, Str: "Inv(" + c.Str + ")"}
}

func (c *ClassicalGate) String() string {
	return c.Str
}
