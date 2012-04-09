package parser

import "fmt"
import "strconv"

var debugActivated bool = false
func debugf(fmtstr string, args ...interface{}) {
	if debugActivated {
		fmt.Printf(fmtstr, args...)
	}
}

func isAlpha(r byte) bool {
	return r >= 'A' && r <= 'z'
}

func isDigit(r byte) bool {
	debugf("isDigit '%c' >= '0': %t, '%c' <= '9': %t\n", r, (r >= '0'), r, (r <= '9'))
	return r >= '0' && r <= '9'
}

  
type DCPULex struct {
	index int
	currentRune byte
	code string
	initialized bool
}

func (lex *DCPULex) nextLetter() byte {
	if lex.index >= len(lex.code) {
		lex.currentRune = 0x0 // end of parsing
		return lex.currentRune
	} else {
		
		lex.currentRune = lex.code[lex.index] ;
		lex.index++ ;
		debugf("lex nextLetter: '%c'\n", lex.currentRune)
		return lex.currentRune
	}
	panic ("shouldnt occur in nextLetter")
}

func (lex *DCPULex) getRune() byte {
	if (!lex.initialized) {
		lex.nextLetter()
		lex.initialized = true
	}
	return lex.currentRune
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
		yylval.reg = DcpuRegister(symbol);
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
		yylval.inst = DcpuInstruction(symbol);
		debugf("lex: found instruction %s\n", symbol)
		return instruction
	}

	panic ("couldnt lex symbol")
}

func (lex *DCPULex) findLabel(yylval *yySymType) int {
	debugf("findLabel\n")
	r := lex.getRune()
	symbol := string(r)
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
	n, _ := strconv.ParseUint(symbol, 0, 16)
	// not sure whawt to do with err
	yylval.lit = DcpuLitteral(n)
	debugf("lex: found litteral %x\n", yylval.lit)
	return litteral
}

func (lex *DCPULex) Lex(yylval *yySymType) int {
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
		debugf("hmm should not happen?\n")
		lex.nextLetter()
		return int(r) // hmm...
	}
	panic ("should not occur in Lex")
}

func (DCPULex) Error(s string) {
	println("syntax error!")
}