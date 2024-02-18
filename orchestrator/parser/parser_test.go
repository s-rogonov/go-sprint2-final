package parser

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTokenizer(t *testing.T) {
	tokens, err := tokenize("1*(2+ 33 ) - 4")
	if err != nil {
		t.Fatal(err)
	}

	expected := []token{
		{TkNumber, "1", 0},
		{TkOp2, "*", 1},
		{TkOpen, "(", 2},
		{TkNumber, "2", 3},
		{TkOp1, "+", 4},
		{TkNumber, "33", 6},
		{TkClose, ")", 9},
		{TkOp1, "-", 11},
		{TkNumber, "4", 13},
	}

	if !cmp.Equal(tokens, expected) {
		t.Error("different tokens sequence")
		t.Error(">>>", tokens)
		t.Error(">>>", expected)
	}
}

func TestInfix2Postfix(t *testing.T) {
	tokens := []token{
		{TkNumber, "1", 0},
		{TkOp2, "*", 1},
		{TkOpen, "(", 2},
		{TkNumber, "2", 3},
		{TkOp1, "+", 4},
		{TkNumber, "33", 6},
		{TkClose, ")", 9},
		{TkOp1, "-", 11},
		{TkNumber, "4", 13},
	}

	postfix, err := infix2postfix(tokens)
	if err != nil {
		t.Fatal(err)
	}

	expected := []token{
		{TkNumber, "1", 0},
		{TkNumber, "2", 3},
		{TkNumber, "33", 6},
		{TkOp1, "+", 4},
		{TkOp2, "*", 1},
		{TkNumber, "4", 13},
		{TkOp1, "-", 11},
	}

	if !cmp.Equal(postfix, expected) {
		t.Error("different postfix sequence")
		t.Error(">>>", postfix)
		t.Error(">>>", expected)
	}
}

func TestBuildAST(t *testing.T) {
	tokens := []token{
		{TkNumber, "1", 0},
		{TkNumber, "2", 3},
		{TkNumber, "33", 6},
		{TkOp1, "+", 4},
		{TkOp2, "*", 1},
		{TkNumber, "4", 13},
		{TkOp1, "-", 11},
	}
	_, err := buildAST(tokens, map[string]time.Duration{
		"+": 1 * time.Second,
		"-": 1 * time.Second,
		"*": 1 * time.Second,
		"/": 1 * time.Second,
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestParseSimpleExpression(t *testing.T) {
	durations := map[string]time.Duration{
		"+": 1 * time.Second,
		"-": 1 * time.Second,
		"*": 1 * time.Second,
		"/": 1 * time.Second,
	}
	q, err := ParseExpression("1+2", durations)
	if err != nil {
		t.Fatal(err)
	}

	if len(q.Tasks) != 3 {
		t.Fatalf(`expected 3 tasks, got %v`, len(q.Tasks))
	}
	t1, t2, t3 := q.Tasks[0], q.Tasks[1], q.Tasks[2]
	if t1.Result != 1 || len(t1.Subtasks) > 0 {
		t.Errorf(`expected plain number "1"`)
		t.Error(">>>", t1)
	}
	if t2.Result != 2 || len(t2.Subtasks) > 0 {
		t.Errorf(`expected plain number "2"`)
		t.Error(">>>", t2)
	}
	if t3.Operation != "+" || t3.Duration != 1*time.Second || len(t3.Subtasks) != 2 {
		t.Errorf(`expected operation "+"`)
		t.Error(">>>", t3)
	}
	t1, t2 = t3.Subtasks[0], t3.Subtasks[1]
	if t1.Result != 1 || len(t1.Subtasks) > 0 {
		t.Errorf(`first argument expected to be a plain number "1"`)
		t.Error(">>>", t1)
	}
	if t2.Result != 2 || len(t1.Subtasks) > 0 {
		t.Errorf(`second argument expected to be a plain number "2"`)
		t.Error(">>>", t2)
	}
}

func TestParseExpression(t *testing.T) {
	durations := map[string]time.Duration{
		"+": 1 * time.Second,
		"-": 1 * time.Second,
		"*": 1 * time.Second,
		"/": 1 * time.Second,
	}

	successParseExpression("1", t, durations)
	successParseExpression("1+2+3", t, durations)
	successParseExpression("1+2+3 * 4 / 5 + 6 + 7 * 8 / 9", t, durations)
	successParseExpression("(1+2+3)*(6* 7+8)", t, durations)

	failParseExpression("", errors.New("there is no expression"), t, durations)
	failParseExpression("-4", errors.New("not enough arguments for \"-\" at 0"), t, durations)
	failParseExpression("1+1 2+2", errors.New("unused arguments are left in expression"), t, durations)
	failParseExpression("1+2+", errors.New("not enough arguments for \"+\" at 3"), t, durations)
	failParseExpression("(1+2", errors.New("unclosed \"(\" at 0"), t, durations)
	failParseExpression("3+(1+2)*4)", errors.New("extra \")\" at 9"), t, durations)
}

func successParseExpression(expr string, t *testing.T, durations map[string]time.Duration) {
	_, err := ParseExpression(expr, durations)
	if err != nil {
		t.Error("expr: ", expr, "err: ", err)
	}
}

func failParseExpression(expr string, mustErr error, t *testing.T, durations map[string]time.Duration) {
	_, err := ParseExpression(expr, durations)
	if err.Error() != mustErr.Error() {
		t.Error("expr: ", expr)
		t.Error(">>>", mustErr)
		t.Error(">>>", err)
	}
}
