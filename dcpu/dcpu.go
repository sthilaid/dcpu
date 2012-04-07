package dcpu

import "fmt"

type Word uint16
type memory [0x10000]Word
type registers [8]Word

// facilitates access to the registers
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

type DCPU struct {
	reg registers
	mem memory
	pc Word
	sp Word
	o Word // overflow

	stopFlag bool
}

type opcodeFun func (a,b *Word, cpu *DCPU)

// ----------------------------------------------------------------------------
// dcpu Initialization

func (reg *registers) init() {
	for i, _ := range reg {
		reg[i] = 0x0
	}
}

func (cpu *DCPU) Init() {
	cpu.reg.init()
	cpu.pc = 0x0
	cpu.sp = 0x0
	cpu.o = 0x0
	cpu.stopFlag = false
}

// ----------------------------------------------------------------------------
// instructions implementation

// 0x00-0x07: register (A, B, C, X, Y, Z, I or J, in that order)
// 0x08-0x0f: [register]
// 0x10-0x17: [next Word + register]
//      0x18: POP / [SP++]
//      0x19: PEEK / [SP]
//      0x1a: PUSH / [--SP]
//      0x1b: SP
//      0x1c: PC
//      0x1d: O
//      0x1e: [next Word]
//      0x1f: next Word (literal)
// 0x20-0x3f: literal value 0x00-0x1f (literal)
func (cpu *DCPU) getRegister(exp Word) *Word {
	debugf(highdebug, "getting regiester: %x", exp & 0x7)
	return &cpu.reg[exp & 0x7]
}

func (cpu *DCPU) read(address Word) *Word {
	// todo wait 1 cycle
	debugf(highdebug, "reading: %x => %x (@ %x)", address, cpu.mem[address], &cpu.mem[address])
	return &cpu.mem[address]
}

func (cpu *DCPU) nextWord() *Word {
	debugf(highdebug, "nextWord pc: %x", cpu.pc)
	w := cpu.read(cpu.pc)
	cpu.pc++
	return w
}

func (cpu *DCPU) evalOperand(n Word, litteralContainer *Word) *Word {
	debugf(highdebug, "evalOperand(%x, %x)", n, litteralContainer)
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
func opNonBasic(nonBasicOp, a *Word, cpu *DCPU) {
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
func opSET(a,b *Word, cpu *DCPU) {
	*a = *b
}

//     0x2: ADD a, b - sets a to a+b, sets O to 0x0001 if there's an overflow, 0x0 otherwise
func opADD(a,b *Word, cpu *DCPU) {
	aVal, bVal := *a, *b
	newAVal := aVal + bVal
	*a = newAVal

	// since the DCPU values are unsigned, addition result must be
	// bigger (or equal) than the two operands, otherwise overflow
	if newAVal < aVal || newAVal < bVal {
		cpu.o = 0x1
	} else {
		cpu.o = 0x0
	}
}

//     0x3: SUB a, b - sets a to a-b, sets O to 0xffff if there's an underflow, 0x0 otherwise
func opSUB(a,b *Word, cpu *DCPU) {
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
func opMUL(a,b *Word, cpu *DCPU) {
	*a = *a**b // lolz
	cpu.o = (*a>>0x10) & 0xffff
}

//     0x5: DIV a, b - sets a to a/b, sets O to ((a<<16)/b)&0xffff. if b==0, sets a and O to 0 instead.
func opDIV(a,b *Word, cpu *DCPU) {
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
func opMOD(a,b *Word, cpu *DCPU) {
	bVal := *b
	if bVal == 0 {
		*a = 0
	} else {
		*a = *a % bVal
	}
}

//     0x7: SHL a, b - sets a to a<<b, sets O to ((a<<b)>>16)&0xffff
func opSHL(a,b *Word, cpu *DCPU) {
	*a = *a << *b
	cpu.o = (*a >> 0x10) & 0xffff
}

//     0x8: SHR a, b - sets a to a>>b, sets O to ((a<<16)>>b)&0xffff
func opSHR(a,b *Word, cpu *DCPU) {
	aVal, bVal := *a, *b
	*a = aVal >> bVal
	cpu.o = ((aVal << 0x10) >> bVal) & 0xffff
}

//     0x9: AND a, b - sets a to a&b
func opAND(a,b *Word, cpu *DCPU) {
	*a = *a & *b
}

//     0xa: BOR a, b - sets a to a|b
func opBOR(a,b *Word, cpu *DCPU) {
	*a = *a | *b
}

//     0xb: XOR a, b - sets a to a^b
func opXOR(a,b *Word, cpu *DCPU) {
	*a = *a ^ *b
}

//     0xc: IFE a, b - performs next instruction only if a==b
func opIFE(a,b *Word, cpu *DCPU) {
	if *a != *b {
		cpu.pc++ // skip if not equal
	}
}

//     0xd: IFN a, b - performs next instruction only if a!=b
func opIFN(a,b *Word, cpu *DCPU) {
	if *a == *b {
		cpu.pc++ // skip if equal
	}
}

//     0xe: IFG a, b - performs next instruction only if a>b
func opIFG(a,b *Word, cpu *DCPU) {
	if *a <= *b {
		cpu.pc++ // skip if greater or equal than
	}
}

//     0xf: IFB a, b - performs next instruction only if (a&b)!=0
func opIFB(a,b *Word, cpu *DCPU) {
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

func (cpu *DCPU) apply(op Word, a, b *Word) {
	debugf(mediumdebug, "applying opcode %d with a/*a: %x / %x b/*b: %x / %x", op, a, *a, b, *b)
	opcodeTable[op](a,b,cpu)
}

// ----------------------------------------------------------------------------
// dcpu management

func (cpu *DCPU) Loadprogram(program []Word) {
	if len(program) > len(cpu.mem) {
		panic("program doesn't fit in dcpu memory space...");
	}

	// load the program into memory
	for i, inst := range program {
	 	cpu.mem[i] = inst
	}

	debugf(lowdebug, "after loading program: %s", cpu.MemDump(5))
}

func (cpu *DCPU) Step() {
	// these are just used if litteral expression are given, since
	// we are using pointers to Words, they must be contained ;p
	var aLitteralContainer, bLitteralContainer Word;
	
	// note prob shouldnt use read here as reading the
	// inst shouldn't take an extra cycle?

	inst := *cpu.nextWord()
	
	op := inst & 0x7              // 4 lower bits
	a := (inst >> 4) & (1<<6 - 1) // bits 10-5 inclusively
	b := inst >>10                // 6 stronger bits

	debugf(mediumdebug, "inst: %x, op: %x a: %x b: %x", inst, op, a, b)

	// here a is evaluated *before* b, as specified in dcpu doc
	aPtr, bPtr := cpu.evalOperand(a, &aLitteralContainer), cpu.evalOperand(b, &bLitteralContainer)
	cpu.apply(op, aPtr, bPtr)
}

// assumes cpu is loaded with the code to run
func (cpu *DCPU) Run() {
	cpu.Init()

	for !cpu.stopFlag {
		cpu.Step()
	}
}

func (cpu *DCPU) Stop() {
	cpu.stopFlag = true
}

//-----------------------------------------------------------------------------
// Debugging

func (mem *memory) Dump(size Word) string {
	str := []byte("[")
	var i Word
	for i=0 ; i<size ; i++ {
		str = append(str, fmt.Sprintf("%x, ", mem[i])...)
	}
	return string(append(str, "]"...))
}

func (cpu *DCPU) MemDump(n Word) string {
	return fmt.Sprintf("cpu mem (first %d bytes): %s mem0: %x", n, cpu.mem.Dump(n), cpu.mem[0])
}


func (reg *registers) String() string {
	return fmt.Sprintf("{A: %x, B: %x, C: %x, X: %x, Y: %x, Z: %x, I: %x, J: %x}",
		reg[A], reg[B], reg[C], reg[X], reg[Y], reg[Z], reg[I], reg[J])
}

func (cpu *DCPU) String() string {
	return fmt.Sprintf("DCPU reg: %s pc: %x sp: %x, o: %x",
		cpu.reg.String(),
		cpu.pc,
		cpu.sp,
		cpu.o)
}

const (
        nodebug = iota
        lowdebug = iota
        mediumdebug = iota
        highdebug = iota
)

var CurrentDebugLevel = nodebug

func debugf(level int, fmtstr string, args ...interface{}) {
	if level <= CurrentDebugLevel {
		println(fmt.Sprintf(fmtstr, args...))
	}
}

