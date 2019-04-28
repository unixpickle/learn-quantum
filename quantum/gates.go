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
		TInv(c, t.Bit)
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
