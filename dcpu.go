package dcpu

type word uint16
type memory [0x10000]word



type registers [8]word
const (
	A = iota
	B = iota
	C = iota
	X = iota
	Y = iota
	Z = iota
	I = iota
	J = iota
)


type dcpu struct {
	reg registers
	mem memory
	pc word
	sp word
	o word // overflow
}

func (cpu dcpu) read(address word) word {
	// todo wait 1 cycle
	return cpu.mem[address]
}

func (cpu dcpu) Loadprogram(program []word) dcpu {
	if len(program) > len(cpu.mem) {
		panic("program doesn't fit in dcpu memory space...");
	}

	// load the program into memory
	for i, inst := range program {
		cpu.mem[i] = inst
	}

	return cpu
}

func (reg registers) init() {
	for i, _ := range reg {
		reg[i] = 0
	}
}

func (cpu dcpu) init() {
	cpu.reg.init()
	cpu.pc = 0
	cpu.sp = 0
	cpu.o = 0
}

// 0x00-0x07: register (A, B, C, X, Y, Z, I or J, in that order)
// 0x08-0x0f: [register]
// 0x10-0x17: [next word + register]
//      0x18: POP / [SP++]
//      0x19: PEEK / [SP]
//      0x1a: PUSH / [--SP]
//      0x1b: SP
//      0x1c: PC
//      0x1d: O
//      0x1e: [next word]
//      0x1f: next word (literal)
// 0x20-0x3f: literal value 0x00-0x1f (literal)
func (cpu dcpu) getRegister(exp word) word {
	return cpu.reg[exp & 0x7]
}

func (cpu dcpu) nextWord() word {
	w := cpu.read(cpu.pc)
	cpu.pc++
	return w
}

func (cpu dcpu) evalOperand(n word) word {
	switch  {
	case n < 0x8: return cpu.getRegister(n)
	case n < 0x10: return cpu.read(cpu.getRegister(n))
	case n < 0x18: return cpu.read(cpu.nextWord() + cpu.getRegister(n))
	case n == 0x18: val := cpu.read(cpu.sp) ; cpu.sp++; return val
	case n == 0x19: return cpu.read(cpu.sp)
	case n == 0x1a: cpu.sp-- ; return cpu.read(cpu.sp)
	case n == 0x1b: return cpu.sp
	case n == 0x1c: return cpu.pc
	case n == 0x1d: return cpu.o
	case n == 0x1e: return cpu.read(cpu.nextWord())
	case n == 0x1f: return cpu.nextWord()
	default: return n - 0x20
	}
	panic("should not occur in evalOperand")
	return 0
}

func (cpu dcpu) apply(op, a, b word) {

}

// assumes cpu is loaded with the code to run
func (cpu dcpu) Run() {
	cpu.init()

	for {
		// note prob shouldnt use read here as reading the
		// inst shouldn't take an extra cycle?
		inst := cpu.nextWord()

		op := inst & 0x7 // 4 lower bits
		a := (inst >> 4) & (1<<7 - 1) // bits 10-5 inclusively
		b := inst >>10 // 6 stronger bits

		// here a is evaluated *before* b, as specified in dcpu doc
		cpu.apply(op, cpu.evalOperand(a), cpu.evalOperand(b))
	}
}