package parser

import "fmt"
import "strconv"

var debugActivated bool = true
func debugf(fmtstr string, args ...interface{}) {
	if debugActivated {
		fmt.Printf(fmtstr, args...)
	}
}

func isAlpha(r byte) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func isDigit(r byte) bool {
	//debugf("isDigit '%c' >= '0': %t, '%c' <= '9': %t\n", r, (r >= '0'), r, (r <= '9'))
	return (r >= '0' && r <= '9') || (r >= 'A' && r <= 'F') || (r >= 'a' && r <= 'f')
}

type DCPULexHack struct {
	index int
	currentRune byte
	code string
	ast DcpuProgram
}

// Warning incoming HACK!
// this structure is a hack because Yacc is forcing us to implement
// the method (lex DCPULex) and not (lex *DCPULex). So with the hack
// container, the pointer to the data is copied, but the data is kept
// safe :)
type DCPULex struct {
	hack *DCPULexHack
}

func (lex *DCPULex)Init(code string) {
	lex.hack = new(DCPULexHack)
	lex.hack.code = code
	lex.hack.ast = DcpuProgram{expressions: []DcpuExpression{}}
	lex.nextLetter()
}

func (lex *DCPULex) nextLetter() byte {
	if lex.hack.index >= len(lex.hack.code) {
		lex.hack.currentRune = 0x0 // end of parsing
		return lex.hack.currentRune
	} else {
		
		lex.hack.currentRune = lex.hack.code[lex.hack.index] ;
		lex.hack.index++ ;
		debugf("lex nextLetter: '%c'\n", lex.hack.currentRune)
		return lex.hack.currentRune
	}
	panic ("shouldnt occur in nextLetter")
}

func (lex *DCPULex) getRune() byte {
	return lex.hack.currentRune
}

func (lex *DCPULex) findSym(yylval *yySymType) int {
	debugf("FindSym\n")
	r := lex.getRune()
	var symbol string = "";
	for isAlpha(r) {
		symbol += string(r)
		r = lex.nextLetter()
	}

	debugf("symbol: %s\n", symbol)

	switch symbol {
	case "A": fallthrough
	case "B": fallthrough
	case "C": fallthrough
	case "X": fallthrough
	case "Y": fallthrough
	case "Z": fallthrough
	case "I": fallthrough
	case "J": 
		yylval.reg = DcpuRegister(symbol)
		debugf("lex: found register %s\n", symbol)
		return register
	case "SET": fallthrough
	case "ADD": fallthrough
	case "SUB": fallthrough
	case "MUL": fallthrough
	case "DIV": fallthrough
	case "MOD": fallthrough
	case "SHL": fallthrough
	case "SHR": fallthrough
	case "AND": fallthrough
	case "BOR": fallthrough
	case "XOR": fallthrough
	case "IFE": fallthrough
	case "IFN": fallthrough
	case "IFG": fallthrough
	case "IFB":
		yylval.inst = DcpuInstruction(symbol)
		debugf("lex: found instruction %s\n", symbol)
		return instruction
	case "POP": fallthrough
	case "PEEK": fallthrough
	case "PUSH": fallthrough
	case "SP": fallthrough
	case "PC": fallthrough
	case "O":
		yylval.specialReg = DcpuSpecialRegister(symbol)
		debugf("lex: found special register %s\n", symbol)
		return specialRegister

	default:
		yylval.lab = DcpuLabel(symbol)
		debugf("lex: assuming label %s\n", symbol)
		return label
	}

	panic ("couldnt lex symbol (should not occur)")
}

func (lex *DCPULex) findLabel(yylval *yySymType) int {
	debugf("findLabel\n")
	r := lex.nextLetter() // dont keep the ":" in the label
	symbol := ""
	for isAlpha(r) {
		symbol += string(r)
		r = lex.nextLetter()
	}

	debugf("lex: found label %s\n", symbol)
	
	yylval.lab = DcpuLabel(symbol)
	return label
}

func (lex *DCPULex) findLitteral(yylval *yySymType) int {
	debugf("findLitteral\n")
	symbol := "0x"
	r := lex.nextLetter()
	for isDigit(r) {
		symbol += string(r)
		r = lex.nextLetter()
	}
	n, err := strconv.ParseUint(symbol, 0, 16)
	// not sure whawt to do with err
	if err != nil {
		debugf("parse err: %s\n", err.Error())
	}
	yylval.lit = DcpuLitteral(n)
	debugf("lex: found litteral %x\n", yylval.lit)
	return litteral
}

func (lex DCPULex) Lex(yylval *yySymType) int {
	r := lex.getRune()
loop:
	debugf("looping with '%c'\n", r)
	switch {
	case r == ':':
		return lex.findLabel(yylval)
		
	case isAlpha(r):
		return lex.findSym(yylval);
		
	case r == '0' && lex.nextLetter() == 'x':
		return lex.findLitteral(yylval)
		
	case r == ' ' || r == '\t' || r == '\n':
		r = lex.nextLetter()
		goto loop
		
	default:
		debugf("passing char '%c' to parser directly...\n", r)
		lex.nextLetter()
		return int(r) // hmm...
	}
	panic ("should not occur in Lex")
}

func (DCPULex) Error(s string) {
	panic("syntax error!")
}