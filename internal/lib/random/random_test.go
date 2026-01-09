package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		Name string
		Size int
	}{
		{
			Name: "Size = 1",
			Size: 1,
		},
		{
			Name: "Size = 5",
			Size: 5,
		},

		{
			Name: "Size = 10",
			Size: 10,
		},
		{
			Name: "Zero size",
			Size: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := NewRandomString(tt.Size)
			assert.Len(t, got, tt.Size)
			t.Logf("Generated string: %s", got)
		})
	}
}
