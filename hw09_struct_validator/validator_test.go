package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/Vadim-Govorukhin/otus-hw/hw09_struct_validator/tags"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:6"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11|regexp:^\\d{11}$"`
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

func TestValidate(t *testing.T) {
	testCases := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "App valid",
			in:   App{"v1.09"},
			expectedErr: ValidationErrors{
				{"Version", nil},
			},
		},
		{
			name: "App invalid",
			in:   App{"VERSION"},
			expectedErr: ValidationErrors{
				{"Version", tags.ErrInvaildByTag},
			},
		},
		{
			name:        "Token always valid",
			in:          Token{[]byte("Header"), []byte("Payload"), []byte("Signature")},
			expectedErr: ValidationErrors{},
		},
		{
			name: "Responce valid",
			in:   Response{200, "anything"},
			expectedErr: ValidationErrors{
				{"Code", nil},
			},
		},
		{
			name: "Responce invalid",
			in:   Response{505, "anything"},
			expectedErr: ValidationErrors{
				{"Code", tags.ErrInvaildByTag},
			},
		},
		{
			name: "User valid",
			in: User{"123456", "Vadim", 18, "valid@example.com", "admin",
				[]string{"12345678901", "10987654321"}, make(json.RawMessage, 2)},
			expectedErr: ValidationErrors{
				{"ID", nil},
				{"Age", nil},
				{"Email", nil},
				{"Role", nil},
				{"Phones", nil},
			},
		},
		{
			name: "User invalid",
			in: User{"123456789", "TooYoung", 15, "invalid@open.me.com", "awesomeFishing",
				[]string{"12345678901", "1098no54321"}, make(json.RawMessage, 2)},
			expectedErr: ValidationErrors{
				{"ID", tags.ErrInvaildByTag},
				{"Age", tags.ErrInvaildByTag},
				{"Email", tags.ErrInvaildByTag},
				{"Role", tags.ErrInvaildByTag},
				{"Phones", tags.ErrInvaildByTag},
			},
		},
	}

	for i, tt := range testCases {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			tt := tt
			fmt.Printf("============= Start test %s =============\n", tt.name)
			//t.Parallel()

			err := Validate(tt.in)
			fmt.Printf("%#v\n", err)
			fmt.Printf("%#v\n", tt.expectedErr)
			fmt.Println(errors.Is(err, tt.expectedErr))
			require.True(t, errors.Is(err, tt.expectedErr))
		})
	}
}
