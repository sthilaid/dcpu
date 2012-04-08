package parser

import "fmt"
import "strconv"

func isAlpha(r byte) bool {
	return r >= 'A' && r <= 'z'
}

func isDigit(r byte) bool {
	fmt.Printf("isDigit '%c' >= '0': %t, '%c' <= '9': %t\n", r, (r >= '0'), r, (r <= '9'))
	return r >= '0' && r <= '9'
}

  
type DCPULex struct {
	index int
	currentRune byte
	code string
	prog DCPUnative
	initialized bool
}

func (lex *DCPULex) nextLetter() byte {
	if lex.index >= len(lex.code) {
		lex.currentRune = 0x0 // end of parsing
		return lex.currentRune
	} else {
		
		lex.currentRune = lex.code[lex.index] ;
		lex.index++ ;
		fmt.Printf("lex nextLetter: '%c'\n", lex.currentRune)
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
	fmt.Printf("FindSym\n")
	r := lex.getRune()
	var symbol string = "";
	for isAlpha(r) {
		symbol += string(r)
		r = lex.nextLetter()
	}

	fmt.Printf("symbol: %s\n", symbol)

	switch symbol {
	case "A": fallthrough
	case "B": fallthrough
	case "C": fallthrough
	case "X": fallthrough
	case "Y": fallthrough
	case "Z": fallthrough
	case "I": fallthrough
	case "J": 
		yylval.vvar = symbol;
		fmt.Printf("lex: found register %s\n", symbol)
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
		yylval.vvar = symbol;
		fmt.Printf("lex: found instruction %s\n", symbol)
		return instruction
	}

	panic ("couldnt lex symbol")
}

func (lex *DCPULex) findLabel(yylval *yySymType) int {
	fmt.Printf("findLabel\n")
	r := lex.getRune()
	symbol := string(r)
	for isAlpha(r) {
		symbol += string(r)
		r = lex.nextLetter()
	}
	yylval.vvar = symbol
	fmt.Printf("lex: found label %s\n", symbol)
	return label

}

func (lex *DCPULex) findLitteral(yylval *yySymType) int {
	fmt.Printf("findLitteral\n")
	symbol := "0x"
	r := lex.nextLetter()
	for isDigit(r) {
		symbol += string(r)
		r = lex.nextLetter()
	}
	n, _ := strconv.ParseUint(symbol, 0, 16)
	// not sure whawt to do with err
	yylval.num = uint16(n)
	fmt.Printf("lex: found litteral %x\n", yylval.num)
	return litteral
}

func (lex *DCPULex) Lex(yylval *yySymType) int {
	r := lex.getRune()
loop:
	fmt.Printf("looping with '%c'\n", r)
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
		fmt.Printf("hmm should not happen?\n")
		lex.nextLetter()
		return int(r) // hmm...
	}
	panic ("should not occur in Lex")
}

func (DCPULex) Error(s string) {
	println("syntax error!")
}