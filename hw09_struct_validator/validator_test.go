package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

type App2 struct {
	Version string `validate:"len:5"`
	Age     int    `validate:"min:18|max:50"`
	Ages    []int  `validate:"min:18|max:50"`
	I       float64
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          App2{"Hello", 20, []int{18, 29}, 0.0},
			expectedErr: ValidationErrors{{"Age", ErrUnsupportedTag}},
		},
	}

	for i, tt := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
