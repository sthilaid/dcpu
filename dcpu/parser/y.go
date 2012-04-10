
//line dcpuAssembly.y:3
package parser

import "fmt"
import "reflect"

//line dcpuAssembly.y:21
type yySymType struct {
	yys int
	expr DcpuExpression
	inst DcpuInstruction
	reg DcpuRegister
	lab DcpuLabel
	lit DcpuLitteral
	operand DcpuOperand
	ref DcpuReference
	sum DcpuSum
	specialReg DcpuSpecialRegister
}

const instruction = 57346
const register = 57347
const specialRegister = 57348
const label = 57349
const litteral = 57350

var yyToknames = []string{
	"instruction",
	"register",
	"specialRegister",
	"label",
	"litteral",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 15
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 25

var yyAct = []int{

	6, 7, 8, 11, 10, 22, 12, 23, 21, 16,
	14, 24, 18, 3, 19, 20, 4, 13, 2, 17,
	15, 5, 9, 1, 25,
}
var yyPact = []int{

	9, -1000, 9, -4, 13, -1000, 1, -1000, -1000, -1000,
	-1000, -1000, 4, -4, -4, -3, -1000, -1000, -7, -2,
	-1000, -1000, 6, -4, -1000, -1000,
}
var yyPgo = []int{

	0, 23, 18, 0, 22, 20, 19,
}
var yyR1 = []int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	4, 5, 5, 5, 6,
}
var yyR2 = []int{

	0, 1, 2, 4, 5, 1, 1, 1, 1, 1,
	3, 1, 1, 1, 3,
}
var yyChk = []int{

	-1000, -1, -2, 4, 7, -2, -3, 5, 6, -4,
	8, 7, 10, 4, 9, -5, 5, -6, 8, -3,
	-3, 11, 12, 9, 5, -3,
}
var yyDef = []int{

	0, -2, 1, 0, 0, 2, 0, 5, 6, 7,
	8, 9, 0, 0, 0, 0, 11, 12, 13, 0,
	3, 10, 0, 0, 14, 4,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 12, 9, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 10, 3, 11,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c > 0 && c <= len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return fmt.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return fmt.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		fmt.Printf("lex %U %s\n", uint(char), yyTokname(c))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		fmt.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				fmt.Printf("%s", yyStatname(yystate))
				fmt.Printf("saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					fmt.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				fmt.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		fmt.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line dcpuAssembly.y:36
		{
		if lexer, ok := yylex.(*DCPULex) ; ok {
			lexer.hack.ast = DcpuProgram{expressions: append(lexer.hack.ast.expressions, yyS[yypt-0].expr)}
		} else {
			panic(fmt.Sprintf("unexected lexer type, got: %s", reflect.TypeOf(lexer)))
		}
	}
	case 2:
		//line dcpuAssembly.y:45
		{
		if lexer, ok := yylex.(*DCPULex) ; ok {
			lexer.hack.ast = DcpuProgram{expressions: append(lexer.hack.ast.expressions, yyS[yypt-1].expr, yyS[yypt-0].expr)}
		} else {
			panic(fmt.Sprintf("unexected lexer type, got: %s", reflect.TypeOf(lexer)))
		}	
	}
	case 3:
		//line dcpuAssembly.y:54
		{
		expr := new(DcpuExpression)
		expr.inst = yyS[yypt-3].inst
		expr.a = yyS[yypt-2].operand
		expr.b = yyS[yypt-0].operand
		expr.label = ""
		yyVAL.expr = *expr
	}
	case 4:
		//line dcpuAssembly.y:64
		{
		expr := new(DcpuExpression)
		expr.inst = yyS[yypt-3].inst
		expr.a = yyS[yypt-2].operand
		expr.b = yyS[yypt-0].operand
		expr.label = yyS[yypt-4].lab
		yyVAL.expr = *expr
	
	}
	case 5:
		//line dcpuAssembly.y:75
		{
		yyVAL.operand = DcpuRegister(yyS[yypt-0].reg)
	}
	case 6:
		//line dcpuAssembly.y:80
		{
		yyVAL.operand = DcpuSpecialRegister(yyS[yypt-0].specialReg)
	}
	case 7:
		//line dcpuAssembly.y:85
		{
		yyVAL.operand = yyS[yypt-0].ref
	}
	case 8:
		//line dcpuAssembly.y:89
		{
		yyVAL.operand = DcpuLitteral(yyS[yypt-0].lit)
	}
	case 9:
		//line dcpuAssembly.y:93
		{
		yyVAL.operand = DcpuLabel(yyS[yypt-0].lab)
	}
	case 10:
		//line dcpuAssembly.y:98
		{
		yyVAL.ref = yyS[yypt-1].ref
	}
	case 11:
		//line dcpuAssembly.y:103
		{
		reference := new (DcpuReference)
		reference.ref = yyS[yypt-0].reg
		yyVAL.ref = *reference
	}
	case 12:
		//line dcpuAssembly.y:109
		{
		reference := new (DcpuReference)
		reference.ref = yyS[yypt-0].sum
		yyVAL.ref = *reference
	}
	case 13:
		//line dcpuAssembly.y:115
		{
		reference := new (DcpuReference)
		reference.ref = yyS[yypt-0].lit
		yyVAL.ref = *reference
	}
	case 14:
		//line dcpuAssembly.y:121
		{
		sum := new(DcpuSum)
		sum.lit = yyS[yypt-2].lit
		sum.reg = yyS[yypt-0].reg
		yyVAL.sum = *sum
	}
	}
	goto yystack /* stack new state and value */
}
