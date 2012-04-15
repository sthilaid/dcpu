
//line dcpuAssembly.y:3
package parser

import "fmt"
import "reflect"
import "dcpu"

//line dcpuAssembly.y:28
type yySymType struct {
	yys int
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

const instruction = 57346
const register = 57347
const specialRegister = 57348
const label = 57349
const litteral = 57350
const dataInstruction = 57351
const stringData = 57352

var yyToknames = []string{
	"instruction",
	"register",
	"specialRegister",
	"label",
	"litteral",
	"dataInstruction",
	"stringData",
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

const yyNprod = 21
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 32

var yyAct = []int{

	9, 10, 11, 14, 13, 29, 28, 21, 15, 30,
	25, 18, 26, 20, 31, 23, 22, 24, 16, 27,
	19, 6, 2, 17, 7, 12, 8, 5, 4, 3,
	1, 32,
}
var yyPact = []int{

	17, -1000, -1000, 17, -1000, -1000, -4, 14, -1000, 0,
	-1000, -1000, -1000, -1000, -1000, 8, -4, 2, -4, -7,
	-1000, -1000, -9, -1000, -2, -1000, -1000, -1000, -1000, 9,
	-4, -1000, -1000,
}
var yyPgo = []int{

	0, 30, 29, 28, 27, 0, 25, 20, 7, 22,
}
var yyR1 = []int{

	0, 1, 9, 9, 2, 2, 3, 3, 5, 5,
	5, 5, 5, 6, 7, 7, 7, 7, 8, 4,
	4,
}
var yyR2 = []int{

	0, 1, 2, 0, 1, 1, 4, 5, 1, 1,
	1, 1, 1, 3, 1, 1, 1, 1, 3, 3,
	3,
}
var yyChk = []int{

	-1000, -1, -9, -2, -3, -4, 4, 7, -9, -5,
	5, 6, -6, 8, 7, 12, 4, 9, 11, -7,
	5, -8, 8, 7, -5, 8, 10, -5, 13, 14,
	11, 5, -5,
}
var yyDef = []int{

	3, -2, 1, 3, 4, 5, 0, 0, 2, 0,
	8, 9, 10, 11, 12, 0, 0, 0, 0, 0,
	14, 15, 16, 17, 0, 19, 20, 6, 13, 0,
	0, 18, 7,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 14, 11, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 12, 3, 13,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10,
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
		//line dcpuAssembly.y:49
		{
		if lexer, ok := yylex.(*DCPULex) ; ok {
			lexer.hack.ast = DcpuProgram{expressions: yyS[yypt-0].exprlst}
		} else {
			panic(fmt.Sprintf("unexected lexer type, got: %s", reflect.TypeOf(lexer)))
		}
	}
	case 2:
		//line dcpuAssembly.y:58
		{
		if lexer, ok := yylex.(*DCPULex) ; ok {
			var expr DcpuExpression = yyS[yypt-1].expr
			var list []DcpuExpression = yyS[yypt-0].exprlst
			// here, since the regular expression solving will
		// match the last expressions first, we must append at
		// the begginning of the slice ;p
		yyVAL.exprlst = append([]DcpuExpression{expr}, list...)
		} else {
			panic(fmt.Sprintf("unexected lexer type, got: %s", reflect.TypeOf(lexer)))
		}
	}
	case 3:
		//line dcpuAssembly.y:72
		{
		yyVAL.exprlst = []DcpuExpression{}
	}
	case 4:
		//line dcpuAssembly.y:77
		{
		yyVAL.expr = yyS[yypt-0].nexpr
	}
	case 5:
		//line dcpuAssembly.y:82
		{
		yyVAL.expr = yyS[yypt-0].dexpr
	}
	case 6:
		//line dcpuAssembly.y:88
		{
		expr := new(DcpuNormalExpression)
		expr.inst = yyS[yypt-3].inst
		expr.a = yyS[yypt-2].operand
		expr.b = yyS[yypt-0].operand
		expr.label = ""
		yyVAL.nexpr = *expr
	}
	case 7:
		//line dcpuAssembly.y:98
		{
		expr := new(DcpuNormalExpression)
		expr.inst = yyS[yypt-3].inst
		expr.a = yyS[yypt-2].operand
		expr.b = yyS[yypt-0].operand
		expr.label = yyS[yypt-4].lab
		yyVAL.nexpr = *expr
	}
	case 8:
		//line dcpuAssembly.y:108
		{
		yyVAL.operand = DcpuRegister(yyS[yypt-0].reg)
	}
	case 9:
		//line dcpuAssembly.y:113
		{
		yyVAL.operand = DcpuSpecialRegister(yyS[yypt-0].specialReg)
	}
	case 10:
		//line dcpuAssembly.y:118
		{
		yyVAL.operand = yyS[yypt-0].ref
	}
	case 11:
		//line dcpuAssembly.y:122
		{
		yyVAL.operand = DcpuLitteral(yyS[yypt-0].lit)
	}
	case 12:
		//line dcpuAssembly.y:126
		{
		yyVAL.operand = DcpuLabel(yyS[yypt-0].lab)
	}
	case 13:
		//line dcpuAssembly.y:131
		{
		yyVAL.ref = yyS[yypt-1].ref
	}
	case 14:
		//line dcpuAssembly.y:136
		{
		reference := new (DcpuReference)
		reference.ref = yyS[yypt-0].reg
		yyVAL.ref = *reference
	}
	case 15:
		//line dcpuAssembly.y:142
		{
		reference := new (DcpuReference)
		reference.ref = yyS[yypt-0].sum
		yyVAL.ref = *reference
	}
	case 16:
		//line dcpuAssembly.y:148
		{
		reference := new (DcpuReference)
		reference.ref = yyS[yypt-0].lit
		yyVAL.ref = *reference
	}
	case 17:
		//line dcpuAssembly.y:154
		{
		reference := new (DcpuReference)
		reference.ref = yyS[yypt-0].lab
		yyVAL.ref = *reference
	}
	case 18:
		//line dcpuAssembly.y:161
		{
		sum := new(DcpuSum)
		sum.lit = yyS[yypt-2].lit
		sum.reg = yyS[yypt-0].reg
		yyVAL.sum = *sum
	}
	case 19:
		//line dcpuAssembly.y:169
		{
		expr := new(DcpuDataExpression)
		expr.data = []dcpu.Word{dcpu.Word(yyS[yypt-0].lit)}
		expr.label = yyS[yypt-2].lab
		yyVAL.dexpr = *expr
	}
	case 20:
		//line dcpuAssembly.y:177
		{
		expr := new(DcpuDataExpression)
		//expr.data = []dcpu.Word([]byte(string(label)))

		str := yyS[yypt-0].str
		data := []dcpu.Word{}
		for _,char := range str {
			data = append(data, dcpu.Word(char))
		}
		expr.data = data
		
		expr.label = yyS[yypt-2].lab
		yyVAL.dexpr = *expr
	}
	}
	goto yystack /* stack new state and value */
}
