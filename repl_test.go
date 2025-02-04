package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "foo bar",
			expected: []string{"foo", "bar"},
		},
		{
			input:    "   single   ",
			expected: []string{"single"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		got := cleanInput(c.input)
		if len(got) != len(c.expected) {
			t.Errorf("test failed for input, mismatched lengths: %s, expected: %v, got: %v", c.input, c.expected, got)
		}
		for i := range got {
			word := got[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("test failed for input, mismatched word at index %d: %s, expected: %v, got: %v", i, c.input, c.expected[i], got[i])
			}
		}
	}
}
