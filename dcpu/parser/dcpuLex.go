package dcpu

import "strconv"

func isAlpha(r byte) bool {
	return r >= 'A' && r <= 'z'
}

func isDigit(r byte) bool {
	return r <= '0' && r >= '9'
}

  
type DCPULex struct {
	index int
	code string
	prog DCPUnative
}

func (lex *DCPULex) nextLetter() byte {
	r := lex.code[lex.index] ;
	lex.index++ ;
	return r
}

func (lex *DCPULex) Lex(yylval *yySymType) int {
	var symbol string
loop:
	r := lex.nextLetter()
	switch {
	case r == ':':
		goto label
	case isAlpha(r):
		goto sym
	case r == '0' && lex.nextLetter() == 'x':
		goto num
	case r == ' ' || r == '\t' || r == '\n':
		goto loop
	default:
		return int(r) // hmm...
		
	}
sym:
	symbol = string(r)
	for isAlpha(r) {
		symbol += string(r)
		r = lex.nextLetter()
	}

	switch symbol {
	case "A":
	case "B":
	case "C":
	case "X":
	case "Y":
	case "Z":
	case "I":
	case "J":
		yylval.vvar = symbol;
		return register
	case "SET":
	case "ADD":
	case "SUB":
	case "MUL":
	case "DIV":
	case "MOD":
	case "SHL":
	case "SHR":
	case "AND":
	case "BOR":
	case "XOR":
	case "IFE":
	case "IFN":
	case "IFG":
	case "IFB":
		yylval.vvar = symbol;
		return instruction
		
	}
label:
	symbol = string(r)
	for isAlpha(r) {
		symbol += string(r)
		r = lex.nextLetter()
	}
	yylval.vvar = symbol
	return label

num:
	symbol = "0x" + string(r)
	for isDigit(r) {
		symbol += string(r)
		r = lex.nextLetter()
	}
	n, _ := strconv.ParseUint(symbol, 0, 16)
	// not sure whawt to do with err
	yylval.num = uint16(n)
	return litteral
}

func (DCPULex) Error(s string) {
	println("syntax error!")
}