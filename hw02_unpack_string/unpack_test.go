package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: ""},
		{name: "charnum", input: "a4bc2d5e", expected: "aaaabccddddde"},
		{name: "two equal chars", input: "abccd", expected: "abccd"},
		{name: "char0", input: "aaa0b", expected: "aab"},
		{name: "char1", input: "a1b", expected: "ab"},
		{name: "rus", input: "б2ж3й1", expected: "ббжжжй"},
		{name: "space", input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{name: "tab", input: "d\t2abc", expected: "d\t\tabc"},
		{name: "slashspace", input: `d\n5abc`, expected: "d\n\n\n\n\nabc"},
		{name: "num", input: `qwe\4\5`, expected: `qwe45`},
		{name: "numrepeat", input: `qwe\45`, expected: `qwe44444`},
		{name: "slashrepeat", input: `qwe\\5`, expected: `qwe\\\\\`},
		{name: "slashnum", input: `qwe\\\3`, expected: `qwe\3`},
		{name: "slashnum", input: `\5qwe\\`, expected: `5qwe\`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `qwe\`, `qwe\\\`, `qwe\q`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
