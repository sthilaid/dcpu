package parser

import "testing"

func TestParsing0(t *testing.T) {
	lex := new(DCPULex)
	lex.Init("SET A, 0x30")
	yyParse(lex)
	
	expr0 := DcpuExpression{inst: DcpuInstruction("SET"),
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
	expr0 := DcpuExpression{inst: DcpuInstruction("ADD"),
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
	expr0 := DcpuExpression{inst: DcpuInstruction("BOR"),
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
	expr0 := DcpuExpression{inst: DcpuInstruction("SET"),
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
	expr0 := DcpuExpression{inst: DcpuInstruction("SET"),
		                a:    DcpuSpecialRegister("PUSH"),
	                        b:    DcpuLitteral(0x10),
	}
	expr1 := DcpuExpression{inst: DcpuInstruction("ADD"),
		                a:    DcpuSpecialRegister("PEEK"),
	                        b:    DcpuLitteral(0x1),
	}


	program := DcpuProgram{expressions: []DcpuExpression{expr0,expr1}}

	if equal, err := program.IsEqualTo(lex.hack.ast) ; !equal {
		t.Errorf(err)
	}
}


/*
func TestParsing1(t *testing.T) {
	lex := new(DCPULex)
	lex.code = "ADD [0xAAAA], 0xFF00 SET PUSH, [0xAAAA]"
	yyParse(lex)

	var expr0 DcpuExpression
	var expr1 DcpuExpression
	var ref0 DcpuReference

	expr0.inst = DcpuInstruction("ADD")
	ref0 = DcpuLitteral(0xaaaa)
	expr0.a = ref0
	expr0.b = DcpuLitteral(0xff00)

	expr1.inst = DcpuInstruction("SET")
	//...
	
	ast := new(DcpuProgram)
	ast.expressions = []DcpuExpression{expr0, expr1}
}
*/