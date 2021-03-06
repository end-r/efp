package efp

type lexer struct {
	buffer    []byte
	offset    int
	line      int
	column    int
	tokens    []token
	numTokens int
	length    int
	errors    []string
}

// processes the next token.
func (l *lexer) next() {
	if l.isEOF() {
		return
	}
	found := false
	for _, pt := range getProtoTokens() {
		if pt.identifier(l.buffer[l.offset]) {
			//fmt.Printf("offset: %d\n", l.offset)
			//fmt.Printf("found: %s\n", pt.name)
			t := pt.process(l)
			if t.tkntype != tknNone {
				l.tokens = append(l.tokens, t)
			} else {
				l.offset++
			}
			found = true
			break
		}
	}
	if !found {
		if l.errors == nil {
			l.errors = make([]string, 0)
		}
		l.errors = append(l.errors, errUnrecognisedToken)
		l.offset++
	}
	l.next()
}

func (l *lexer) isEOF() bool {
	return l.offset >= l.length
}

// creates a new string from the token's value
// TODO: escaped characters
func (l *lexer) tokenString(t token) string {
	data := make([]byte, t.end-t.start)
	copy(data, l.buffer[t.start:t.end])
	return string(data)
}

func (l *lexer) nextByte() byte {
	b := l.buffer[l.offset]
	l.offset++
	return b
}

func lexString(str string) *lexer {
	return lex([]byte(str))
}

func lex(bytes []byte) *lexer {
	l := new(lexer)
	l.buffer = bytes
	l.length = len(bytes)
	l.next()
	return l
}

func processNewLine(l *lexer) token {
	l.line++
	return token{
		tkntype: tknNone,
	}
}

func processIgnored(l *lexer) token {
	return token{
		tkntype: tknNone,
	}
}

func processNumber(l *lexer) (t token) {
	t.start = l.offset
	t.end = l.offset
	t.tkntype = tknNumber
	for '0' <= l.buffer[l.offset] && l.buffer[l.offset] <= '9' {
		l.offset++
		t.end++
		if l.isEOF() {
			return t
		}
	}
	return t
}

func processIdentifier(l *lexer) token {

	t := new(token)
	t.start = l.offset
	t.end = l.offset
	t.tkntype = tknValue
	if l.isEOF() {
		return *t
	}
	for isIdentifier(l.buffer[l.offset]) {
		//fmt.Printf("id: %c\n", l.buffer[l.offset])
		t.end++
		l.offset++
		if l.isEOF() {
			return *t
		}
	}
	return *t
}

// processes a string sequence to create a new token.
func processString(l *lexer) token {
	// the start - end is the value
	// it DOES include the enclosing quotation marks
	t := new(token)
	t.start = l.offset
	t.end = l.offset
	t.tkntype = tknString
	b := l.nextByte()
	b2 := l.nextByte()
	for b != b2 {
		t.end++
		b2 = l.nextByte()
		if l.isEOF() {
			t.end += 2
			return *t
		}
	}
	t.end += 2
	return *t
}
