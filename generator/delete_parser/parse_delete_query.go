package delete_parser

import (
	"fmt"
	"strconv"
	"unicode"
)

type Token struct {
	Type  string
	Value string
}

func (t *Token) String() string {
	return fmt.Sprintf("{Type: \"%s\" Value: \"%s\"}", t.Type, t.Value)
}

type ParsedKeyRange struct {
	Start []*Token
	End   []*Token
	Kind  string
	Table string
}

type Parser struct {
	State        string
	Query        string
	Pos          int
	Char         rune
	CurrentToken *Token
	KeyRange     *ParsedKeyRange
	EOF          bool
}

func NewParser(q string) *Parser {
	kr := &ParsedKeyRange{
		Start: make([]*Token, 0),
		End:   make([]*Token, 0),
	}
	return &Parser{Query: q, KeyRange: kr}
}

func (in *Parser) GetNextToken() (*Token, error) {
	if in.Pos >= len(in.Query) {
		return &Token{Type: "EOF", Value: ""}, nil
	}

	char := rune(in.Query[in.Pos])

	if unicode.IsSpace(char) {
		if in.CurrentToken != nil && in.State != "DEFINING_STRING" {
			in.SkipWhitespace()
		}
	}
	tokStr := in.Scan()
	//fmt.Printf("got from scan: %s\n", tokStr)
	switch tokStr {
	case "START(":
		return &Token{Type: "START_DECLARATION", Value: tokStr}, nil
	case "END(":
		return &Token{Type: "END_DECLARATION", Value: tokStr}, nil
	case "KIND(":
		return &Token{Type: "KIND_DECLARATION", Value: tokStr}, nil
	case "FROM":
		return &Token{Type: "TABLE_DECLARATION", Value: tokStr}, nil
	case "DELETE":
		return &Token{Type: "DELETE", Value: tokStr}, nil
	case "CC":
		return &Token{Type: "KIND_TYPE", Value: "ClosedClosed"}, nil
	case "CO":
		return &Token{Type: "KIND_TYPE", Value: "ClosedOpen"}, nil
	case "OC":
		return &Token{Type: "KIND_TYPE", Value: "OpenClosed"}, nil
	case "OO":
		return &Token{Type: "KIND_TYPE", Value: "OpenOpen"}, nil
	case "'":
		if in.State != "DEFINING_STRING" {
			in.State = "DEFINING_STRING"
			return &Token{Type: "STRING_BEGIN", Value: tokStr}, nil
		} else {
			in.State = ""
			return &Token{Type: "STRING_END", Value: tokStr}, nil
		}
	case ",":
		return &Token{Type: ",", Value: tokStr}, nil
	case ")":
		return &Token{Type: ")", Value: tokStr}, nil
	case ".":
		return &Token{Type: ".", Value: tokStr}, nil
	case "?":
		return &Token{Type: "?", Value: tokStr}, nil
	default:
		return &Token{Type: "IDENTIFIER", Value: tokStr}, nil
	}
}

func (in *Parser) Scan() string {
	if in.State == "DEFINING_STRING" {
		if in.Char == '\'' {
			in.Advance()
			return "'"
		}
		str := ""
		for {
			if in.Char == '\'' || in.EOF {
				break
			}
			str += string(in.Char)
			//fmt.Printf("str now: %s\n", str)
			in.Advance()
		}
		//fmt.Printf("returning from scan: %s\n", str)
		return str
	}
	if in.IsSeperator(in.Char) {
		//fmt.Printf("Cur char is sep: %s\n", string(in.Char))
		str := string(in.Char)
		in.Advance()
		return str
	}
	str := ""
	for {
		if in.IsKeyword(str) || in.EOF || in.IsSeperator(in.Char) {
			//fmt.Printf("Cur char is something important: %s\n", str)
			break
		}
		str += string(in.Char)
		//fmt.Printf("str now: %s\n", str)
		in.Advance()
	}
	return str
}

func (in *Parser) IsSeperator(c rune) bool {
	if c == ' ' || c == ')' || c == ',' || c == '\'' || c == '.' {
		return true
	}
	return false
}

func (in *Parser) IsKeyword(st string) bool {
	if contains([]string{"DELETE", "FROM", "START(", "END(", "KIND(", ",", ")", "'", ".", "CC", "CO", "OC", "OO", "?"}, st) {
		return true
	}
	return false
}

func (in *Parser) Advance() {
	in.Pos += 1
	if in.Pos >= len(in.Query) {
		in.EOF = true
	} else {
		in.Char = rune(in.Query[in.Pos])
	}
}

func (in *Parser) SkipWhitespace() {
	for {
		if in.EOF || !unicode.IsSpace(in.Char) {
			break
		}
		in.Advance()
	}
}

// compare current token type to t, if they match
// eat the current token, assign the current token to the
// next token, otherwise return an error
func (in *Parser) Eat(t string) error {
	//fmt.Printf("eat: %+v, expected: %+v\n", in.CurrentToken, ts)
	if t == in.CurrentToken.Type {
		tok, err := in.GetNextToken()
		if err != nil {
			return err
		}
		in.CurrentToken = tok
		//fmt.Printf("current token now: %s\n", tok)
	} else {
		return fmt.Errorf("expected token of type: %s, but got: %s at pos: %d", t, in.CurrentToken, in.Pos)
	}
	return nil
}

// eats a ' token, identifier token, and ' token (' IDENTIFIER ') and returns a new token
// of type STRING  with the token's value being the string value of the identifier
func (in *Parser) GetStringIdentifier() (*Token, error) {
	if err := in.Eat("STRING_BEGIN"); err != nil {
		return nil, err
	}
	strVal := ""
	for {
		if in.CurrentToken.Type == "STRING_END" {
			break
		}
		ident := in.CurrentToken
		strVal += string(ident.Value)
		if err := in.Eat("IDENTIFIER"); err != nil {
			return nil, err
		}
	}
	if err := in.Eat("STRING_END"); err != nil {
		return nil, err
	}
	return &Token{Type: "STRING", Value: "\"" + strVal + "\""}, nil
}

// eats an identifier token, and if there is a . token, it eats that
// and parses another identifier.  It returns either an INT token,
// or a FLOAT token based on if it parsed a . or not
// It leaves the value as a string
func (in *Parser) GetNumberIdentifier() (*Token, error) {
	first := in.CurrentToken
	if err := in.Eat("IDENTIFIER"); err != nil {
		return nil, err
	}
	cur := in.CurrentToken
	if cur.Type == "." {
		if err := in.Eat("."); err != nil {
			return nil, err
		}
		second := in.CurrentToken
		if err := in.Eat("IDENTIFIER"); err != nil {
			return nil, err
		}
		if isNumeric(first.Value + "." + second.Value) {
			return &Token{Type: "FLOAT", Value: first.Value + "." + second.Value}, nil
		} else {
			return nil, fmt.Errorf("not a valid FLOAT near position: %d got: %s.%s", in.Pos, first.Value, second.Value)
		}
	} else {
		if isNumeric(first.Value) {
			return &Token{Type: "INT", Value: first.Value}, nil
		} else {
			return nil, fmt.Errorf("not a valid INT near position: %d got %s", in.Pos, first.Value)
		}
	}
}

func (in *Parser) DeclareTable() error {
	if err := in.Eat("TABLE_DECLARATION"); err != nil {
		return err
	}
	ident := in.CurrentToken
	if err := in.Eat("IDENTIFIER"); err != nil {
		return err
	}
	if isNumeric(ident.Value) {
		return fmt.Errorf("invalid table declaration near position: %d  table name cannot be numeric: %s", in.Pos, ident.Value)
	}
	// WRAP TABLE IN ""
	in.KeyRange.Table = fmt.Sprintf("\"%s\"", ident.Value)
	return nil
}

//Gets the
func (in *Parser) DeclareKind() error {
	if err := in.Eat("KIND_DECLARATION"); err != nil {
		return err
	}
	ident := in.CurrentToken
	in.KeyRange.Kind = ident.Value
	if err := in.Eat("KIND_TYPE"); err != nil {
		return err
	}
	if err := in.Eat(")"); err != nil {
		return err
	}
	return nil
}

// we should be have our current token pointint to an identifier
func (in *Parser) DeclareIdentList(tokens *[]*Token) error {
	for {
		cur := in.CurrentToken
		//fmt.Printf("what is cur for list: %s\n", cur)
		switch cur.Type {
		case ")":
			break
		case "STRING_BEGIN":
			tok, err := in.GetStringIdentifier()
			if err != nil {
				return err
			}
			*tokens = append(*tokens, tok)
		case "IDENTIFIER":
			tok, err := in.GetNumberIdentifier()
			if err != nil {
				return err
			}
			*tokens = append(*tokens, tok)
		case "?":
			tok := in.CurrentToken
			in.Eat("?")
			*tokens = append(*tokens, tok)
		default:
			return fmt.Errorf("not a valid token in identifier list")
		}
		if cur.Type == ")" {
			break
		}
		in.Eat(",")
	}
	in.Eat(")")
	return nil
}

func (in *Parser) DeclareStart() error {
	if err := in.Eat("START_DECLARATION"); err != nil {
		return err
	}
	if err := in.DeclareIdentList(&in.KeyRange.Start); err != nil {
		return err
	}
	return nil
}

func (in *Parser) DeclareEnd() error {
	if err := in.Eat("END_DECLARATION"); err != nil {
		return err
	}
	if err := in.DeclareIdentList(&in.KeyRange.End); err != nil {
		return err
	}
	return nil
}

//def: DELETE TABLE_DECLARATION START_DECLARATION END_DELCARTION KIND_DECLARATION
//DELETE => DELETE
// TABLE_DECLARATION => FROM CHARS
// START_DECLARATION => Start(IDENT_LIST)
// END_DECLARATION => End(IDENT_LIST)
// KIND_DECLARATION => Kind(STRING)
// IDENT_LIST => IDENT | IDENT, IDENT_LIST
// IDENT => STRING | FLOAT | INTS | ?
// FLOAT => INTS.INTS | .INTS
// INTS => INT | INTS INT
// STRING => "CHARS"
// CHARS => CHAR | CHARS CHAR

func (in *Parser) Expr() (*ParsedKeyRange, error) {
	if len(in.Query) > 0 {
		in.Char = rune(in.Query[0])
	}
	tok, err := in.GetNextToken()
	if err != nil {
		return nil, err
	}
	in.CurrentToken = tok
	in.Eat("DELETE")
	var hasTable, hasKind, hasStart, hasEnd bool
	for {
		next := in.CurrentToken
		if err != nil {
			//fmt.Println("next err: %s\n", err)
			return nil, err
		}
		//fmt.Printf("EXPR TOKEN %s\n", next)
		switch next.Type {
		case "TABLE_DECLARATION":
			if hasTable {
				return nil, fmt.Errorf("you already declared a table")
			}
			err := in.DeclareTable()
			if err != nil {
				return nil, err
			}
			hasTable = true

		case "KIND_DECLARATION":
			if hasKind {
				return nil, fmt.Errorf("you already declared a kind")
			}
			err := in.DeclareKind()
			if err != nil {
				return nil, err
			}
			hasKind = true
		case "START_DECLARATION":
			if hasStart {
				return nil, fmt.Errorf("you already declared a start")
			}
			err := in.DeclareStart()
			if err != nil {
				return nil, err
			}
			hasStart = true
		case "END_DECLARATION":
			if hasEnd {
				return nil, fmt.Errorf("you already declared a end")
			}
			err := in.DeclareEnd()
			if err != nil {
				return nil, err
			}
			hasEnd = true
		default:
			//fmt.Printf("what token is this? %s\n", next)
			return nil, fmt.Errorf("not a declaration: %s", next)
		}
		if (hasTable && hasKind && hasStart && hasEnd) || in.EOF {
			break
		}
	}
	return in.KeyRange, nil
}

func contains(ts []string, t string) bool {
	for _, a := range ts {
		if a == t {
			return true
		}
	}
	return false
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
