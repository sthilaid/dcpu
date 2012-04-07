package dcpu

import "testing"

func TestSetRegister(t *testing.T) {
	cpu := new(DCPU)
	cpu.Init()
	// SET A, 0x30 ; 7c01 0030
	cpu.Loadprogram([]Word{0x7c01, 0x0030})
	cpu.Step()
	if cpu.reg[A] != 0x30 {
		t.Errorf("Simple test failed: cpu.reg[A]: %x != 0x30", cpu.reg[A])
	}
}

