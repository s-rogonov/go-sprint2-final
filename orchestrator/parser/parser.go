package parser

import (
	"fmt"
)

const (
	TkNumber = 0
	TkOpen   = 3
	TkClose  = 4
	TkOp1    = 10001
	TkOp2    = 10002
)

func tkType2string(tkType int) string {
	switch tkType {
	case TkNumber:
		return "TkNumber"
	case TkOp1:
		return "TkOp1"
	case TkOp2:
		return "TkOp2"
	case TkOpen:
		return "TkOpen"
	case TkClose:
		return "TkClose"
	default:
		return "<Unknown>"
	}
}

type token struct {
	TkType     int
	TkText     string
	TkStartPos int
}

func (t token) String() string {
	return fmt.Sprintf(
		`{%v, "%v", %v}`,
		tkType2string(t.TkType), t.TkText, t.TkStartPos,
	)
}

func tokenize(expr string) (stream []token, err error) {
	runes := []rune(expr)
	var value []rune
	for pos, char := range runes {
		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			value = append(value, char)
		default:
			if len(value) > 0 {
				stream = append(stream, token{
					TkType:     TkNumber,
					TkText:     string(value),
					TkStartPos: pos - len(value),
				})
				value = []rune{}
			}
		}

		switch char {
		case ' ':
		case '+', '-':
			stream = append(stream, token{
				TkType:     TkOp1,
				TkText:     string([]rune{char}),
				TkStartPos: pos,
			})
		case '*', '/':
			stream = append(stream, token{
				TkType:     TkOp2,
				TkText:     string([]rune{char}),
				TkStartPos: pos,
			})
		case '(':
			stream = append(stream, token{
				TkType:     TkOpen,
				TkText:     string([]rune{char}),
				TkStartPos: pos,
			})
		case ')':
			stream = append(stream, token{
				TkType:     TkClose,
				TkText:     string([]rune{char}),
				TkStartPos: pos,
			})
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			return nil, fmt.Errorf("unknown symbol `%v` at %d", char, pos)
		}
	}

	if len(value) > 0 {
		stream = append(stream, token{
			TkType:     TkNumber,
			TkText:     string(value),
			TkStartPos: len(runes) - len(value),
		})
		value = []rune{}
	}

	return
}

func infix2postfix(tokens []token) ([]token, error) {
	var stream []token
	var stack []token

	for _, t := range tokens {
		switch t.TkType {
		case TkNumber:
			stream = append(stream, t)
		case TkOpen:
			stack = append(stack, t)
		case TkClose:
			for len(stack) > 0 && stack[len(stack)-1].TkType != TkOpen {
				stream = append(stream, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf(`extra ")" at %v`, t.TkStartPos)
			}
			stack = stack[:len(stack)-1] // remove open brace
		case TkOp1, TkOp2:
			for len(stack) > 0 && stack[len(stack)-1].TkType >= t.TkType {
				stream = append(stream, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, t)
		default:
			// skip unknown tokens
		}
	}

	for len(stack) > 0 && stack[len(stack)-1].TkType != TkOpen {
		stream = append(stream, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	if len(stack) > 0 {
		return nil, fmt.Errorf(`unclosed "(" at %v`, stack[len(stack)-1].TkStartPos)
	}
	return stream, nil
}

// func ParseExpression(expr string) (*models.Query, error) {
// 	tokens, err := tokenize(expr)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	_, err = infix2postfix(tokens)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return nil, nil
// }
