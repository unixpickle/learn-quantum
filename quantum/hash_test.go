package quantum

import (
	"testing"
)

func TestEquivalence(t *testing.T) {
	c1 := Circuit{&TGate{Bit: 2}, &HGate{Bit: 4}, &TGate{Bit: 1}}
	c2 := Circuit{&TGate{Bit: 2}, &TGate{Bit: 1}, &HGate{Bit: 4}}
	hasher := NewCircuitHasher(5)
	if hasher.Hash(c1) != hasher.Hash(c2) {
		t.Error("mismatching hashes")
	}
}
