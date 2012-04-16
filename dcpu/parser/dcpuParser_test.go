package parser

import "dcpu"
import "fmt"
import "testing"

//-----------------------------------------------------------------------------
// testing correct contruction of AST from assembly

func TestParsing0(t *testing.T) {
	lex := new(DCPULex)
	lex.Init("SET A, 0x30")
	yyParse(lex)
	
	expr0 := DcpuNormalExpression{inst: DcpuInstruction("SET"),
		                      a:    DcpuRegister("A"),
	                              b:    DcpuLitteral(0x30)}

	program := DcpuProgram{expressions: []DcpuExpression{expr0}}

	if equal, err := program.IsEqualTo(lex.hack.ast) ; !equal {
		t.Errorf(err)
	}
}

func TestParsing1(t *testing.T) {
	lex := new(DCPULex)
	lex.Init("ADD [X], [0x1]")
	yyParse(lex)
	expr0 := DcpuNormalExpression{inst: DcpuInstruction("ADD"),
		                      a:    DcpuReference{ref: DcpuRegister("X")},
	                              b:    DcpuReference{ref: DcpuLitteral(0x1)},
	}

	program := DcpuProgram{expressions: []DcpuExpression{expr0}}

	if equal, err := program.IsEqualTo(lex.hack.ast) ; !equal {
		t.Errorf(err)
	}
}

func TestParsing2(t *testing.T) {
	lex := new(DCPULex)
	lex.Init("BOR [0xab +I], [0x1]")
	yyParse(lex)
	expr0 := DcpuNormalExpression{inst: DcpuInstruction("BOR"),
	                              a:    DcpuReference{DcpuSum{lit: DcpuLitteral(0xab), reg: DcpuRegister("I")}},
	                              b:    DcpuReference{DcpuLitteral(0x1)},
	}

	program := DcpuProgram{expressions: []DcpuExpression{expr0}}

	if equal, err := program.IsEqualTo(lex.hack.ast) ; !equal {
		t.Errorf(err)
	}
}

func TestParsing3(t *testing.T) {
	lex := new(DCPULex)
	lex.Init("SET PC, 0x0")
	yyParse(lex)
	expr0 := DcpuNormalExpression{inst: DcpuInstruction("SET"),
		                      a:    DcpuSpecialRegister("PC"),
	                              b:    DcpuLitteral(0x0),
	}

	program := DcpuProgram{expressions: []DcpuExpression{expr0}}

	if equal, err := program.IsEqualTo(lex.hack.ast) ; !equal {
		t.Errorf(err)
	}
}

func TestParsing4(t *testing.T) {
	lex := new(DCPULex)
	lex.Init("SET PUSH, 0x10 ADD PEEK, 0x1")
	yyParse(lex)
	expr0 := DcpuNormalExpression{inst: DcpuInstruction("SET"),
		                      a:    DcpuSpecialRegister("PUSH"),
	                              b:    DcpuLitteral(0x10),
	}
	expr1 := DcpuNormalExpression{inst: DcpuInstruction("ADD"),
		                      a:    DcpuSpecialRegister("PEEK"),
	                              b:    DcpuLitteral(0x1),
	}


	program := DcpuProgram{expressions: []DcpuExpression{expr0,expr1}}

	if equal, err := program.IsEqualTo(lex.hack.ast) ; !equal {
		t.Errorf(err)
	}
}

func TestParsing5(t *testing.T) {
	lex := new(DCPULex)
	lex.Init(":loop SET PC, loop")
	yyParse(lex)
	expr0 := DcpuNormalExpression{inst: DcpuInstruction("SET"),
		                      a:    DcpuSpecialRegister("PC"),
		                      b:    DcpuLabel("loop"),
	                              label: "loop",
	}

	program := DcpuProgram{expressions: []DcpuExpression{expr0}}

	if equal, err := program.IsEqualTo(lex.hack.ast) ; !equal {
		t.Errorf(err)
	}
}

func TestParsing6(t *testing.T) {
	lex := new(DCPULex)
	lex.Init(":dat DAT 0xbadd")
	yyParse(lex)
	expr0 := DcpuDataExpression{label: "dat",
	                            data: []dcpu.Word{0xbadd},
	}

	program := DcpuProgram{expressions: []DcpuExpression{expr0}}

	if equal, err := program.IsEqualTo(lex.hack.ast) ; !equal {
		t.Errorf(err)
	}
}

func TestParsing7(t *testing.T) {
	lex := new(DCPULex)
	lex.Init(":dat DAT \"toto\"")
	yyParse(lex)
	expr0 := DcpuDataExpression{label: "dat",
		                    data: []dcpu.Word{0x74, 0x6f, 0x74, 0x6f, 0x0},
	}

	program := DcpuProgram{expressions: []DcpuExpression{expr0}}

	if equal, err := program.IsEqualTo(lex.hack.ast) ; !equal {
		t.Errorf(err)
	}
}

//-----------------------------------------------------------------------------
// Parsing failure checks

func TestParsingFailure0(t *testing.T) {
	defer func() {
		if panic := recover(); panic == nil {
			t.Errorf("Should not parse code with invalid syntax (missing ',')...")
		}
	}()
			
	lex := new(DCPULex)
	// missing a "," after 'a'
	lex.Init("SET PUSH 0x10")
	yyParse(lex)
}

func TestParsingFailure1(t *testing.T) {
	defer func() {
		if panic := recover(); panic == nil {
			t.Errorf("Should not permit addition outside ref...")
		}
	}()
			
	lex := new(DCPULex)
	// not allowed to have addition outside a reference
	lex.Init("SET 0x10+B, 0x10")
	yyParse(lex)
}

func TestParsingFailure2(t *testing.T) {
	defer func() {
		if panic := recover(); panic == nil {
			t.Errorf("Should not be allowed to referece special registers...")
		}
	}()
			
	lex := new(DCPULex)
	// not allowed to have addition outside a reference
	lex.Init("SET [SP], 0x0")
	yyParse(lex)
}

func TestParsingFailure3(t *testing.T) {
	defer func() {
		if panic := recover(); panic == nil {
			t.Errorf("Should not parse invalid instructions...")
		}
	}()
			
	lex := new(DCPULex)
	lex.Init("SOT A, 0x0")
	yyParse(lex)
}

//-----------------------------------------------------------------------------
// Validing AST -> binary conversion

// returns false if the 2 slices are different
func compareBinaries(bin1, bin2 []dcpu.Word) bool {
	size := len(bin1)
	if size != len(bin2) {
		return false
	}
	
	for i := 0 ; i < size ; i++ {
		if bin1[i] != bin2[i] {
			return false
		}
	}

	return true
}

func dump(mem *[]dcpu.Word) string {
	str := []byte("[")
	size := len(*mem)
	for i := 0 ; i < size ; i++ {
		str = append(str, fmt.Sprintf("0x%x, ", (*mem)[i])...)
	}
	return string(append(str, "]"...))
}


// simple instruction
func TestBinary0(t *testing.T) {
	lex := new(DCPULex)
	lex.Init("SET A, 0x30")
	yyParse(lex)

	code := lex.hack.ast.Code()
	expectedCode := []dcpu.Word{0x7c01, 0x0030}

	if !compareBinaries(code, expectedCode) {
		t.Errorf("Binaries does not correspond, got %s expcted %s", dump(&code), dump(&expectedCode))
	}
}

// register ref and next word ref
func TestBinary1(t *testing.T) {
	lex := new(DCPULex)
	lex.Init("ADD [X], [0x1]")
	yyParse(lex)
	code := lex.hack.ast.Code()
	expectedCode := []dcpu.Word{0x78b2, 0x1}

	if !compareBinaries(code, expectedCode) {
		t.Errorf("Binaries does not correspond, got %s expcted %s", dump(&code), dump(&expectedCode))
	}
}

// indirect dereferencing
func TestBinary2(t *testing.T) {
	lex := new(DCPULex)
	lex.Init("BOR [0xab +I], [0x1]")
	yyParse(lex)
	code := lex.hack.ast.Code()
	expectedCode := []dcpu.Word{0x796a, 0xab, 0x1}

	if !compareBinaries(code, expectedCode) {
		t.Errorf("Binaries does not correspond, got %s expcted %s", dump(&code), dump(&expectedCode))
	}
}

// special registers
func TestBinary3(t *testing.T) {
	lex := new(DCPULex)
	lex.Init("SET PC, 0x0")
	yyParse(lex)
	code := lex.hack.ast.Code()
	expectedCode := []dcpu.Word{0x81c1}

	if !compareBinaries(code, expectedCode) {
		t.Errorf("Binaries does not correspond, got %s expcted %s", dump(&code), dump(&expectedCode))
	}
}

// more complex 2 instruction program
func TestBinary4(t *testing.T) {
	lex := new(DCPULex)
	lex.Init("SET PUSH, 0x10 ADD PEEK, 0x1")
	yyParse(lex)
	code := lex.hack.ast.Code()
	expectedCode := []dcpu.Word{0xc1a1, 0x8592}

	if !compareBinaries(code, expectedCode) {
		t.Errorf("Binaries does not correspond, got %s expcted %s", dump(&code), dump(&expectedCode))
	}
}

// label
func TestBinary5(t *testing.T) {
	lex := new(DCPULex)
	lex.Init(":loop SET PC, loop")
	yyParse(lex)
	code := lex.hack.ast.Code()
	expectedCode := []dcpu.Word{0x81c1}

	if !compareBinaries(code, expectedCode) {
		t.Errorf("Binaries does not correspond, got %s expcted %s", dump(&code), dump(&expectedCode))
	}
}

// testing ref to label
// note: also testing upper and lower case code parsing...
func TestBinary6(t *testing.T) {
	lex := new(DCPULex)
	lex.Init("set pc, next MUL sP, 0x1 :next xOr a, B")
	yyParse(lex)
	code := lex.hack.ast.Code()
	expectedCode := []dcpu.Word{0x89c1, 0x85b4, 0x040b}

	if !compareBinaries(code, expectedCode) {
		t.Errorf("Binaries does not correspond, got %s expcted %s", dump(&code), dump(&expectedCode))
	}
}

// code using 'dat', making sure that the data is at the end of the
// code and that the labels are set correctly
func TestBinary7(t *testing.T) {
	lex := new(DCPULex)
	lex.Init(":message DAT 'hello' SET I, message SET A, [0x1+I]")
	yyParse(lex)
	code := lex.hack.ast.Code()
	expectedCode := []dcpu.Word{0x8c61, 0x5801, 0x1, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x0}

	if !compareBinaries(code, expectedCode) {
		t.Errorf("Binaries does not correspond, got %s expcted %s", dump(&code), dump(&expectedCode))
	}
}
