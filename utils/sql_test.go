package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnPlaceholder(t *testing.T) {
	type Test struct {
		get  []string
		want string
	}
	tests := []Test{
		{[]string{"A"}, "`A`"},
		{[]string{"A", "B"}, "`A`, `B`"},
	}
	for _, test := range tests {
		assert.Equal(t, ColumnPlaceholder(test.get...), test.want)
	}
}
