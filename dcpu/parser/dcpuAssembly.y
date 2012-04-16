// -*- mode: go -*-
%{
package parser

import "fmt"
import "reflect"
import "dcpu"
%}

%start program
%token <inst> instruction
%token <reg> register
%token <specialReg> specialRegister
%token <lab> label
%token <lit> litteral
%token <dat> dataInstruction
%token <str> stringData
%type <expr> expression
%type <nexpr> normalExpression
%type <dexpr> dataExpression
%type <operand> operand
%type <ref> reference
%type <ref> referenceValue
%type <sum> sum
%type <prog> program
%type <exprlst> expressionList

%union {
	prog DcpuProgram
	expr DcpuExpression
	nexpr DcpuNormalExpression
	dexpr DcpuDataExpression
	inst DcpuInstruction
	reg DcpuRegister
	lab DcpuLabel
	lit DcpuLitteral
	dat string
	str string
	operand DcpuOperand
	ref DcpuReference
	sum DcpuSum
	specialReg DcpuSpecialRegister
	exprlst []DcpuExpression
}

%%
	
program: expressionList
{
	if lexer, ok := yylex.(*DCPULex) ; ok {
		lexer.hack.ast = DcpuProgram{expressions: $1}
	} else {
		panic(fmt.Sprintf("unexected lexer type, got: %s", reflect.TypeOf(lexer)))
	}
}

expressionList: expression expressionList 
{
	if lexer, ok := yylex.(*DCPULex) ; ok {
		var expr DcpuExpression = $1
		var list []DcpuExpression = $2

		if _, ok := expr.(DcpuDataExpression) ; ok {
			// data expression should be appended *at the
			// end of the code segment*
			$$ = append(list, expr)
		} else {
			// here, since the regular expression solving will
			// match the last expressions first, we must append at
			// the begginning of the slice ;p
			$$ = append([]DcpuExpression{expr}, list...)
		}
		
		
	} else {
		panic(fmt.Sprintf("unexected lexer type, got: %s", reflect.TypeOf(lexer)))
	}
}

expressionList:
{
	$$ = []DcpuExpression{}
}

expression: normalExpression
{
	$$ = $1
}

expression: dataExpression
{
	$$ = $1
}


normalExpression: instruction operand ',' operand
{
	expr := new(DcpuNormalExpression)
	expr.inst = $1
	expr.a = $2
	expr.b = $4
	expr.label = ""
	$$ = *expr
}

normalExpression: label instruction operand ',' operand
{
	expr := new(DcpuNormalExpression)
	expr.inst = $2
	expr.a = $3
	expr.b = $5
	expr.label = $1
	$$ = *expr
}

operand: register
{
	$$ = DcpuRegister($1)
}

operand: specialRegister
{
	$$ = DcpuSpecialRegister($1)
}

operand: reference
{
	$$ = $1
}
operand: litteral
{
	$$ = DcpuLitteral($1)
}
operand: label
{
	$$ = DcpuLabel($1)
}

reference: '[' referenceValue ']'
{
	$$ = $2
}

referenceValue: register
{
	reference := new (DcpuReference)
	reference.ref = $1
	$$ = *reference
}
referenceValue: sum
{
	reference := new (DcpuReference)
	reference.ref = $1
	$$ = *reference
}
referenceValue: litteral
{
	reference := new (DcpuReference)
	reference.ref = $1
	$$ = *reference
}
referenceValue: label
{
	reference := new (DcpuReference)
	reference.ref = $1
	$$ = *reference
}

sum: litteral '+' register
{
	sum := new(DcpuSum)
	sum.lit = $1
	sum.reg = $3
	$$ = *sum
}

dataExpression: label dataInstruction litteral
{
	expr := new(DcpuDataExpression)
	expr.data = []dcpu.Word{dcpu.Word($3)}
	expr.label = $1
	$$ = *expr
}

dataExpression: label dataInstruction stringData
{
	expr := new(DcpuDataExpression)
	//expr.data = []dcpu.Word([]byte(string(label)))

	str := $3
	data := []dcpu.Word{}
	for _,char := range str {
		data = append(data, dcpu.Word(char))
	}
	expr.data = data
	
	expr.label = $1
	$$ = *expr
}
