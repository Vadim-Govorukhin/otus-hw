package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty cmd and env", func(t *testing.T) {
		r := RunCmd([]string{}, Environment{})
		require.Equal(t, -1, r)
	})

	t.Run("empty env", func(t *testing.T) {
		r := RunCmd([]string{"ls"}, Environment{})
		require.Equal(t, 0, r)
	})
}
