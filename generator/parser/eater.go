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
	kindName := TokenNames[tkn.tk]
	var expectedNames []string
	var exampleVals []string
	for _, e := range expected {
		expectedNames = append(expectedNames, TokenNames[e])
		exampleVals = append(exampleVals, TokenKeyWords[e]...)
	}
	e.lastErr = fmt.Errorf(
		"Unexpected token: %#v, kind: %s.\n"+
			"Near Position: %d.\n Expected one in kinds: %+v\n"+
			"Examples: %+v",
		*tkn,
		kindName,
		e.scanner.pos,
		expectedNames,
		exampleVals,
	)
	return false
}

// eats each token in group, in order till the FIRST token in the group does not match
// a comma resets.  It only fails if a noncomplete group is eaten, or nothing is eaten
func (e *eater) EatCommaSeperatedGroupOf(group ...hasable) ([][]*Token, bool) {
	if e.lastErr != nil {
		return nil, false
	}
	tokenGroups := make([][]*Token, 0)
	scanOneGroup := func() (g []*Token) {
		for _, kinds := range group {
			if e.Eat(kinds.Values()...) {
				g = append(g, e.Top())
			}
		}
		return
	}
	if !group[0].Has(e.scanner.Peek(1)[0].tk) {
		groupNames := make([]string, len(group))
		for i, kind := range group[0].Values() {
			groupNames[i] = TokenNames[kind]
		}
		e.lastErr = fmt.Errorf("asked to eat a group of %+v, but none was found", groupNames)
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
	if e.scanner.Peek(1)[0].tk == CLOSE_PARAN {
		e.Eat(CLOSE_PARAN)
		return []*Token{}, e.lastErr == nil
	}
	for {
		if !e.Eat(expected...) {
			return nil, false
		}
		// populate values with the top of the array
		values = append(values, e.Top())
		peeked := e.scanner.Peek(1)[0].tk
		if peeked != COMMA {
			break
		}
		if !e.Eat(COMMA) {
			return nil, false
		}

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
