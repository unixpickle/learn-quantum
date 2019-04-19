package quantum

import (
	"fmt"
	"strconv"
)

// A Gate is a generic object that modifies a quantum
// computer in some primitive way.
type Gate interface {
	fmt.Stringer
	Apply(c Computer)
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

type YGate struct {
	Bit int
}

func (y *YGate) String() string {
	return "Y(" + strconv.Itoa(y.Bit) + ")"
}

func (y *YGate) Apply(c Computer) {
	Y(c, y.Bit)
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

type CNotGate struct {
	Control int
	Target  int
}

func (c *CNotGate) String() string {
	return fmt.Sprintf("CNot(%d, %d)", c.Control, c.Target)
}

func (c *CNotGate) Apply(computer Computer) {
	computer.CNot(c.Control, c.Target)
}
