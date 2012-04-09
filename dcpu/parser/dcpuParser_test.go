package parser

import "testing"

func TestParsing0(t *testing.T) {
	lex := new(DCPULex)
	lex.code = "SET A, 0x30"
	yyParse(lex)

	if len(ParsedProgram.expressions) != 1 {
		t.Errorf("TestParsing0 failed: didn't parse the right number of expressions! %d",
			len(ParsedProgram.expressions))
	} else {

		expr := ParsedProgram.expressions[0]
		
		if expr.inst != DcpuInstruction("SET") {
			t.Errorf("TestParsing0 failed: wrong instruction (%s)", expr.inst)
		}
		
		if reg, ok := expr.a.(DcpuRegister) ; !ok || reg != DcpuRegister("A") {
			t.Errorf("TestParsing0 failed: wrong 1st argument (%s != \"A\")", reg)
		}

		if lit, ok := expr.b.(DcpuLitteral) ; !ok || lit != DcpuLitteral(0x30) {
			t.Errorf("TestParsing0 failed: wrong 2nd argument (0x%x != 0x30)", lit)
		}
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