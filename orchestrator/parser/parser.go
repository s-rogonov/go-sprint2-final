package parser

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"dbprovider/models"
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

func buildAST(tokens []token, durations map[string]time.Duration) (*models.Query, error) {
	var stack []*models.Task
	q := &models.Query{}

	for _, t := range tokens {
		switch t.TkType {
		case TkNumber:
			val, err := strconv.ParseFloat(t.TkText, 64)
			if err != nil {
				return nil, fmt.Errorf(`cannot parse number "%v" at %v`, t.TkText, t.TkStartPos)
			}
			task := &models.Task{
				Result: val,
			}
			stack = append(stack, task)
			q.Tasks = append(q.Tasks, task)
		default:
			d, ok := durations[t.TkText]
			if !ok {
				return nil, fmt.Errorf(`no duration for operation "%v" at %v`, t.TkText, t.TkStartPos)
			}
			n := len(stack)
			if n < 2 {
				return nil, fmt.Errorf(`not enough arguments for "%v" at %v`, t.TkText, t.TkStartPos)
			}
			arg1, arg2 := stack[n-2], stack[n-1]
			stack = stack[:n-2]
			task := &models.Task{
				Operation: t.TkText,
				Duration:  d,
				Subtasks:  []*models.Task{arg1, arg2},
			}
			stack = append(stack, task)
			q.Tasks = append(q.Tasks, task)
		}
	}

	if len(stack) > 1 {
		return nil, errors.New("unused arguments are left in expression")
	}
	if len(stack) == 0 {
		return nil, errors.New("there is no expression")
	}

	return q, nil
}

func ParseExpression(expr string, durations map[string]time.Duration) (*models.Query, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return nil, err
	}

	postfix, err := infix2postfix(tokens)
	if err != nil {
		return nil, err
	}

	return buildAST(postfix, durations)
}
