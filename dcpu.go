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

	stopFlag bool
}

type opcodeFun func (a,b *word, cpu dcpu)

func (cpu dcpu) read(address word) *word {
	// todo wait 1 cycle
	return &cpu.mem[address]
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
	cpu.stopFlag = false
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
func (cpu dcpu) getRegister(exp word) *word {
	return &cpu.reg[exp & 0x7]
}

func (cpu dcpu) nextWord() *word {
	w := cpu.read(cpu.pc)
	cpu.pc++
	return w
}

func (cpu dcpu) evalOperand(n word, litteralContainer *word) *word {
	switch  {
	case n < 0x8: return cpu.getRegister(n)
	case n < 0x10: return cpu.read(*cpu.getRegister(n))
	case n < 0x18: return cpu.read(*cpu.nextWord() + *cpu.getRegister(n))
	case n == 0x18: val := cpu.read(cpu.sp) ; cpu.sp++; return val
	case n == 0x19: return cpu.read(cpu.sp)
	case n == 0x1a: cpu.sp-- ; return cpu.read(cpu.sp)
	case n == 0x1b: return &cpu.sp
	case n == 0x1c: return &cpu.pc
	case n == 0x1d: return &cpu.o
	case n == 0x1e: return cpu.read(*cpu.nextWord())
	case n == 0x1f: return cpu.nextWord()
	default: *litteralContainer = n - 0x20 ; return litteralContainer
	}
	panic("should not occur in evalOperand")
}

// * SET, AND, BOR and XOR take 1 cycle, plus the cost of a and b
// * ADD, SUB, MUL, SHR, and SHL take 2 cycles, plus the cost of a and b
// * DIV and MOD take 3 cycles, plus the cost of a and b
// * IFE, IFN, IFG, IFB take 2 cycles, plus the cost of a and b, plus 1 if the test fails

//     0x0: non-basic instruction - see below
func opNonBasic(nonBasicOp, a *word, cpu dcpu) {
	// Non-basic opcodes always have their lower four bits unset, have one value and a six bit opcode.
	// In binary, they have the format: aaaaaaoooooo0000
	// The value (a) is in the same six bit format as defined earlier.
	//
	// Non-basic opcodes: (6 bits)
	//          0x00: reserved for future expansion
	//          0x01: JSR a - pushes the address of the next instruction to the stack, then sets PC to a
	//     0x02-0x3f: reserved
	//
	// * JSR takes 2 cycles, plus the cost of a.

	// since there not many atm, they are directly implemented in the switch...
	switch  *nonBasicOp {
	case 0x01: // JSR a
		cpu.mem[cpu.sp] = cpu.pc
		cpu.sp--
		cpu.pc = *a
	}
}

//     0x1: SET a, b - sets a to b
func opSET(a,b *word, cpu dcpu) {
	*a = *b
}

//     0x2: ADD a, b - sets a to a+b, sets O to 0x0001 if there's an overflow, 0x0 otherwise
func opADD(a,b *word, cpu dcpu) {
	aVal, bVal := *a, *b
	newAVal := aVal + bVal
	*a = newAVal

	// since the dcpu values are unsigned, addition result must be
	// bigger (or equal) than the two operands, otherwise overflow
	if newAVal < aVal || newAVal < bVal {
		cpu.o = 0x1
	} else {
		cpu.o = 0x0
	}
}

//     0x3: SUB a, b - sets a to a-b, sets O to 0xffff if there's an underflow, 0x0 otherwise
func opSUB(a,b *word, cpu dcpu) {
	aVal, bVal := *a, *b
	newAVal := aVal - bVal
	*a = newAVal

	// since the dcpu values are unsigned, substraction result
	// must be smaller (or equal) than the two operands, otherwise
	// underflow
	if newAVal > aVal || newAVal > bVal {
		cpu.o = 0xffff
	} else {
		cpu.o = 0x0
	}	
}

//     0x4: MUL a, b - sets a to a*b, sets O to ((a*b)>>16)&0xffff
func opMUL(a,b *word, cpu dcpu) {
	*a = *a**b // lolz
	cpu.o = (*a>>0x10) & 0xffff
}

//     0x5: DIV a, b - sets a to a/b, sets O to ((a<<16)/b)&0xffff. if b==0, sets a and O to 0 instead.
func opDIV(a,b *word, cpu dcpu) {
	if *b == 0 {
		*a = 0x0
		cpu.o = 0x0
	} else {
		aVal, bVal := *a, *b
		*a = aVal / bVal
		cpu.o = ((aVal<<0x10) / bVal) & 0xffff
	}
}

//     0x6: MOD a, b - sets a to a%b. if b==0, sets a to 0 instead.
func opMOD(a,b *word, cpu dcpu) {
	bVal := *b
	if bVal == 0 {
		*a = 0
	} else {
		*a = *a % bVal
	}
}

//     0x7: SHL a, b - sets a to a<<b, sets O to ((a<<b)>>16)&0xffff
func opSHL(a,b *word, cpu dcpu) {
	*a = *a << *b
	cpu.o = (*a >> 0x10) & 0xffff
}

//     0x8: SHR a, b - sets a to a>>b, sets O to ((a<<16)>>b)&0xffff
func opSHR(a,b *word, cpu dcpu) {
	aVal, bVal := *a, *b
	*a = aVal >> bVal
	cpu.o = ((aVal << 0x10) >> bVal) & 0xffff
}

//     0x9: AND a, b - sets a to a&b
func opAND(a,b *word, cpu dcpu) {
	*a = *a & *b
}

//     0xa: BOR a, b - sets a to a|b
func opBOR(a,b *word, cpu dcpu) {
	*a = *a | *b
}

//     0xb: XOR a, b - sets a to a^b
func opXOR(a,b *word, cpu dcpu) {
	*a = *a ^ *b
}

//     0xc: IFE a, b - performs next instruction only if a==b
func opIFE(a,b *word, cpu dcpu) {
	if *a != *b {
		cpu.pc++ // skip if not equal
	}
}

//     0xd: IFN a, b - performs next instruction only if a!=b
func opIFN(a,b *word, cpu dcpu) {
	if *a == *b {
		cpu.pc++ // skip if equal
	}
}

//     0xe: IFG a, b - performs next instruction only if a>b
func opIFG(a,b *word, cpu dcpu) {
	if *a <= *b {
		cpu.pc++ // skip if greater or equal than
	}
}

//     0xf: IFB a, b - performs next instruction only if (a&b)!=0
func opIFB(a,b *word, cpu dcpu) {
	if (*a & *b) == 0 {
		cpu.pc++ // skip if (a&b)==0
	}
}

var opcodeTable = [0x10]opcodeFun{
	opNonBasic,
	opSET,
	opADD,
	opSUB,
	opMUL,
	opDIV,
	opMOD,
	opSHL,
	opSHR,
	opAND,
	opBOR,
	opXOR,
	opIFE,
	opIFN,
	opIFG,
	opIFB,
}

func (cpu dcpu) apply(op word, a, b *word) {
	opcodeTable[op](a,b,cpu)
}


func (cpu dcpu) Step() {
	// these are just used if litteral expression are given, since
	// we are using pointers to words, they must be contained ;p
	var aLitteralContainer, bLitteralContainer word;
	
	// note prob shouldnt use read here as reading the
	// inst shouldn't take an extra cycle?
	inst := *cpu.nextWord()

	op := inst & 0x7 // 4 lower bits
	a := (inst >> 4) & (1<<7 - 1) // bits 10-5 inclusively
	b := inst >>10 // 6 stronger bits

	// here a is evaluated *before* b, as specified in dcpu doc
	aPtr, bPtr := cpu.evalOperand(a, &aLitteralContainer), cpu.evalOperand(b, &bLitteralContainer)
	cpu.apply(op, aPtr, bPtr)
}

// assumes cpu is loaded with the code to run
func (cpu dcpu) Run() {
	cpu.init()

	for !cpu.stopFlag {
		cpu.Step()
	}
}

func (cpu dcpu) Stop() {
	cpu.stopFlag = true
}