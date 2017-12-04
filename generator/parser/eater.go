package parser

import (
	"fmt"
)

type eater struct {
	scanner    *Scanner
	lastTokens []*Token
	lastErr    error
}

func NewEater(scanner *Scanner) *eater {
	return &eater{scanner: scanner}
}

func (e *eater) Eat(expected ...TokenKind) bool {
	if e.lastErr != nil {
		return false
	}
	tkn := e.scanner.Scan()
	for _, kind := range expected {
		if tkn.tk == kind {
			e.lastTokens = append(e.lastTokens, tkn)
			return true
		}
	}
	e.lastErr = fmt.Errorf(
		"unexpected token: %+v, expected one in kinds: %+v",
		*tkn,
		expected,
	)
	return false
}

// eats each token in group, in order till the FIRST token in the group does not match
// a comma resets.  It only fails if a noncomplete group is eaten, or nothing is eaten
func (e *eater) EatCommaSeperatedGroupOf(group ...TokenKind) ([][]*Token, bool) {
	if e.lastErr != nil {
		return nil, false
	}
	tokenGroups := make([][]*Token, 0)
	scanOneGroup := func() (g []*Token) {
		for _, kind := range group {
			if e.Eat(kind) {
				g = append(g, e.Top())
			}
		}
		return
	}
	if e.scanner.Peek(1)[0].tk == group[0] {
		e.lastErr = fmt.Errorf("asked to eat a group of %+v, but none was found", group)
		return nil, false
	}
	for {
		group := scanOneGroup()
		if e.lastErr != nil {
			return nil, false
		}
		tokenGroups = append(tokenGroups, group)
		if e.scanner.Peek(1)[0].tk != COMMA {
			break
		}
		e.Eat(COMMA)
	}

	return tokenGroups, true
}

// eats the pattern: <(expected, expected, ...)> and returns the expected
func (e *eater) EatArrayOf(expected ...TokenKind) ([]*Token, bool) {
	// for returning the value tokens in the array
	var values []*Token

	e.Eat(OPEN_PARAN)
	e.Eat(expected...)
	for {
		if e.scanner.Peek(1)[0].tk == CLOSE_PARAN {
			break
		}
		if !e.Eat(COMMA) {
			return nil, false
		}
		if !e.Eat(expected...) {
			return nil, false
		}
		// populate values with the top of the array
		values = append(values, e.Top())
	}
	e.Eat(CLOSE_PARAN)

	return values, e.lastErr == nil
}

func (e *eater) Top() *Token {
	if len(e.lastTokens) == 0 {
		return &Token{tk: ILLEGAL, raw: "cannot get top of empty token slice"}
	}
	return e.lastTokens[len(e.lastTokens)-1]
}

func (e *eater) TakeTokens() []*Token {
	tkns := e.lastTokens
	e.lastTokens = nil
	return tkns
}

func (e *eater) TakeErr() error {
	err := e.lastErr
	e.lastErr = nil
	return err
}
