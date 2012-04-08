%{
package parser

import "fmt"
import "dcpu"

type DCPUnative []dcpu.Word

var dcpuLexerReturnValue DCPUnative = DCPUnative{}
%}

%start program
%token <vvar> instruction
%token <vvar> register
%token <vvar> label
%token <numb> litteral

%union {
	vvar string
	num uint16
}

%%

program: expression | expression expression
expression: instruction operand ',' operand
{
  println("testing!")
  dcpuLexerReturnValue = append(dcpuLexerReturnValue, 0xF0F0)
}

expression: label instruction operand ',' operand  ;

operand: register | reference | litteral | label ;
reference: '[' referenceValue ']' ;
referenceValue: register | sum  | litteral ;
sum: litteral '+' register
