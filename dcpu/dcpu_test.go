package dcpu

import "testing"

func TestSetRegister0(t *testing.T) {
	cpu := new(DCPU)
	cpu.Init()
	// SET A, 0x30 ; 7c01 0030
	cpu.Loadprogram([]Word{0x7c01, 0x0030})
	cpu.Step()
	if cpu.reg[A] != 0x30 {
		t.Errorf("Simple test 0 failed: cpu.reg[A]: %x != 0x30", cpu.reg[A])
	}
}

func TestSetRegister1(t *testing.T) {
	cpu := new(DCPU)
	cpu.Init()
	// SET A, 0x30 ; 7c01 0030
	// SET PC, 0x0 ; 81c1
	cpu.Loadprogram([]Word{0x7c01, 0x0030, 0x81c1})
	cpu.Step()
	cpu.Step()
	if cpu.reg[A] != 0x30 {
		t.Errorf("Simple test 1 failed: cpu.reg[A]: %x != 0x30", cpu.reg[A])
	}

	if cpu.pc != 0 {
		t.Errorf("Simple test 1 failed: cpu.pc = %x != 0", cpu.pc)
	}
}

func BenchmarkSetRegister(b *testing.B) {
	cpu := new(DCPU)
	cpu.Init()
	// SET A, 0x30 ; 7c01 0030
	// SET PC, 0x0 ; 81c1
	cpu.Loadprogram([]Word{0x7c01, 0x0030, 0x81c1})

	for i := 0 ; i < b.N ; i++ {
		cpu.Step()
	}
}

