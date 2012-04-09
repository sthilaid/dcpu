// -*- mode: go -*-
%{
package parser

import "fmt"

var ParsedProgram DcpuProgram
%}

%start program
%token <inst> instruction
%token <reg> register
%token <lab> label
%token <lit> litteral
%type <expr> expression
%type <operand> operand
%type <ref> reference
%type <ref> referenceValue
%type <sum> sum

%union {
	expr DcpuExpression
	inst DcpuInstruction
	reg DcpuRegister
	lab DcpuLabel
	lit DcpuLitteral
	operand DcpuOperand
	ref DcpuReference
	sum DcpuSum
}

%%
	
program: expression
{
	ParsedProgram.expressions = append(ParsedProgram.expressions, $1)
}

program: expression expression
{
	ParsedProgram.expressions = append(ParsedProgram.expressions, $1, $2)
}

expression: instruction operand ',' operand
{
	expr := new(DcpuExpression)
	expr.inst = $1
	expr.a = $2
	expr.b = $4
	expr.label = ""
	$$ = *expr
}

expression: label instruction operand ',' operand
{
	expr := new(DcpuExpression)
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
sum: litteral '+' register
{
	sum := new(DcpuSum)
	sum.lit = $1
	sum.reg = $3
	$$ = *sum
}