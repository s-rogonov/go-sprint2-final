package parser

import (
	"testing"

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
