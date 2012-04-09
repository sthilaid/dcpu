package parser

import "fmt"
import "dcpu"

//*****************************************************************************

//--------------------
// program description
//--------------------

type DcpuComparable interface {
	IsEqualTo(comp DcpuComparable) (isEqual bool, errStr string)
}

type DcpuProgram struct {
	expressions []DcpuExpression
	labels map[DcpuLabel]byte
}

func (prog DcpuProgram)processLabels() {
	var currentSize byte = 0
	prog.labels = make(map[DcpuLabel]byte)

	for _, expr := range prog.expressions {
		if expr.label != "" {
			prog.labels[expr.label] = currentSize
			currentSize += expr.Size(prog)
		}
	}
}

func (prog DcpuProgram)Code() []dcpu.Word {
	code := []dcpu.Word{}
	for _, expr := range prog.expressions {
		code = append(code, expr.Code(prog)...)
	}
	return code
}

func (prog DcpuProgram)String() string {
	str := "Program {"
	for _,expr := range prog.expressions {
		str += expr.String() + ", "
	}
	return str + "}\n"
}

func (prog DcpuProgram)IsEqualTo(prog1 DcpuProgram) (bool, string) {
	if len(prog.expressions) != len(prog1.expressions) {
		return false, "Programs are not of same length!"
	}
	for i,expr := range prog.expressions {
		if equal,err := expr.IsEqualTo(prog1.expressions[i]) ; !equal {
			return false, fmt.Sprintf("expressions %d are not equal: %s", i, err)
		}
	}
	return true, ""
}

//-----------------------
// Expression description
//-----------------------

type DcpuExpression struct {
	inst DcpuInstruction
	a DcpuOperand
	b DcpuOperand
	label DcpuLabel
}

func (exp DcpuExpression)Code(prog DcpuProgram) []dcpu.Word {
	opCode := exp.inst.Code()
	a := exp.a.Code(prog)
	b := exp.b.Code(prog)

	binaryExpr := b[0] << 10 + a[0] << 4 + opCode

	intermediateResult := append([]dcpu.Word{binaryExpr}, a[1:]...)
	return append(intermediateResult, b[1:]...)
}

func (exp DcpuExpression)Size(prog DcpuProgram) byte {
	return byte(1) + exp.a.Size(prog) + exp.b.Size(prog)
}

func (exp DcpuExpression)String() string {
	return string(exp.inst) +" "+ exp.a.String() +", "+ exp.b.String()
}

func (expr DcpuExpression)IsEqualTo(expr1 DcpuExpression) (bool, string) {
	if (expr.inst != expr1.inst) {
		return false, fmt.Sprintf("Instructions are different (%s, %s)", expr.inst, expr1.inst)
	} else if equal, str := expr.a.IsEqualTo(expr1.a) ; !equal {
		return false, fmt.Sprintf("'a' operands are different (%s)", str)
	} else if equal, str := expr.b.IsEqualTo(expr1.b) ; !equal {
		return false, fmt.Sprintf("'b' operands are different (%s)", str)
	}
	return true, ""
}

//*****************************************************************************

//-----------------------------
// instruction node description
//-----------------------------

type DcpuInstruction string

func (inst DcpuInstruction)Code() dcpu.Word {
	switch inst {
	case "SET": return 0x1
	case "ADD": return 0x2
	case "SUB": return 0x3
	case "MUL": return 0x4
	case "DIV": return 0x5
	case "MOD": return 0x6
	case "SHL": return 0x7
	case "SHR": return 0x8
	case "AND": return 0x9
	case "BOR": return 0xa
	case "XOR": return 0xb
	case "IFE": return 0xc
	case "IFN": return 0xd
	case "IFG": return 0xe
	case "IFB": return 0xf
	case "JSR": return 0x10
	}
	panic(fmt.Sprintf("unknown instruction: %s", inst))
}

//*****************************************************************************

//---------------------------
// operand node meta description
//---------------------------

type DcpuOperand interface {
	Code(prog DcpuProgram) []dcpu.Word
	Size(prog DcpuProgram) byte
	String() string
	IsEqualTo(op DcpuComparable) (isEqual bool, errStr string)
}

//---------------------------
// reference node description
//---------------------------

type DcpuReference struct {
	ref interface {
		ReferenceCode(prog DcpuProgram) []dcpu.Word
		ReferenceSize(prog DcpuProgram) byte
		String() string
		IsEqualTo(ref DcpuComparable) (bool,string)
	}
}
func (refNode DcpuReference)Code(prog DcpuProgram) []dcpu.Word {
	return refNode.ref.ReferenceCode(prog)
}

func (refNode DcpuReference)Size(prog DcpuProgram) byte {
	return refNode.ref.ReferenceSize(prog)
}

func (refNode DcpuReference)String() string {
	return refNode.ref.String()
}

func (refNode DcpuReference)IsEqualTo(op DcpuComparable) (bool, string) {
	if ref1, ok := op.(DcpuReference) ; !ok {
		return false, fmt.Sprintf("different types of operands (%s, %s)", refNode, op)
	} else if equal, err := refNode.ref.IsEqualTo(ref1.ref) ; !equal {
		return false, fmt.Sprintf("Incompabible references: %s", err)
	}
	return true, ""
}


//---------------------------
// register node description
//---------------------------

type DcpuRegister string

func (reg DcpuRegister) Code(prog DcpuProgram) []dcpu.Word {
	registers := []DcpuRegister{"A", "B", "C", "X", "Y", "Z", "I", "J"}
	for i, r := range registers {
		if reg == r {
			return []dcpu.Word{dcpu.Word(i)}
		}
	}
	panic(fmt.Sprintf("Couldn't find register: %s", reg))
}
func (reg DcpuRegister) Size(prog DcpuProgram) byte {
	return 0
}

func (reg DcpuRegister) ReferenceCode(prog DcpuProgram) []dcpu.Word {
	// register Code always returns one Word
	return []dcpu.Word{reg.Code(prog)[0] + 0x8}
}

func (reg DcpuRegister) ReferenceSize(prog DcpuProgram) byte {
	return 0
}

func (reg DcpuRegister) String() string {
	return string(reg)
}

func (reg DcpuRegister) IsEqualTo(op DcpuComparable) (bool,string) {
	if reg1, ok := op.(DcpuRegister) ; !ok {
		return false, fmt.Sprintf("different types of operands (%s, %s)", reg, op)
	} else if reg != reg1 {
		return false, fmt.Sprintf("Different registers (%s, %s)", reg, reg1)
	}
	return true, ""
}

//--------------------------------
// litteral value node description
//--------------------------------

type DcpuLitteral dcpu.Word

func (lit DcpuLitteral)isEmbeded() bool {
	return lit < 0x20
}

func (lit DcpuLitteral)Code(prog DcpuProgram) []dcpu.Word {
	value := dcpu.Word(lit)
	if  lit.isEmbeded() {
		// embeded litteral
		return []dcpu.Word{value + 0x20}
	}
	// else next word
	return []dcpu.Word{0x1f, value}
}

func (lit DcpuLitteral)Size(prog DcpuProgram) byte {
	if lit.isEmbeded() {
		return 0
	}
	return 1
}

func (lit DcpuLitteral)ReferenceCode(prog DcpuProgram) []dcpu.Word {
	return []dcpu.Word{0x1e, dcpu.Word(lit)}
}

func (lit DcpuLitteral)ReferenceSize(prog DcpuProgram) byte {
	return 1
}

func (lit DcpuLitteral)String() string {
	return fmt.Sprintf("0x%x", uint16(lit))
}

func (lit DcpuLitteral)IsEqualTo(op DcpuComparable) (bool,string) {
	if lit1, ok := op.(DcpuLitteral) ; !ok {
		return false, fmt.Sprintf("different types of operands (%s, %s)", lit, op)
	} else if lit != lit1 {
		return false, fmt.Sprintf("different litteral values (%p, %p)", lit, lit1)
	}
	return true, ""
}

//--------------------------------
// sum regerence node description
//--------------------------------

type DcpuSum struct {
 	reg DcpuRegister
 	lit DcpuLitteral
}

func (sum DcpuSum)ReferenceCode(prog DcpuProgram) []dcpu.Word {
	// registers Code always returns 1 word
	return []dcpu.Word{sum.reg.Code(prog)[0] + 0x10, dcpu.Word(sum.lit)}
}

func (sum DcpuSum)ReferenceSize(prog DcpuProgram) byte {
	return 1
}

func (sum DcpuSum)String() string {
	return sum.lit.String() +" + "+ sum.reg.String()
}

func (sum DcpuSum)IsEqualTo(op DcpuComparable) (bool,string) {
	if sum1, ok := op.(DcpuSum) ; !ok {
		return false, fmt.Sprintf("different types of operands (%s, %s)", sum, op)
	} else if equal, err := sum.reg.IsEqualTo(sum1.reg) ; !equal {
		return false, fmt.Sprintf("sum registers are different: %s", err)
	} else if equal, err := sum.lit.IsEqualTo(sum1.lit) ; !equal {
		return false, fmt.Sprintf("sum litterals are different: %s", err)
	}
	return true, ""
}

//--------------------------------
// label node description
//--------------------------------

type DcpuLabel string

func (label DcpuLabel)Code(prog DcpuProgram) []dcpu.Word {
	value := prog.labels[label]
	lit := DcpuLitteral(value)
	return lit.Code(prog)
}

func (label DcpuLabel)Size(prog DcpuProgram) byte {
	value := prog.labels[label]
	lit := DcpuLitteral(value)
	return lit.Size(prog)
}

func (label DcpuLabel)String() string {
	return string(label)
}

func (label DcpuLabel)IsEqualTo(op DcpuComparable) (bool,string) {
	if label1, ok := op.(DcpuLabel) ; !ok {
		return false, fmt.Sprintf("different types of operands (%s, %s)", label, op)
	} else if label != label1 {
		return false, fmt.Sprintf("different labels (%s, %s)", label, label1)
	}
	return true, ""
}
