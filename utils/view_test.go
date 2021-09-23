package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSiteTitle(t *testing.T) {
	type Test struct {
		get  []string
		want string
	}

	tests := []Test{
		Test{[]string{"A", "B"}, "A - B"},
		Test{[]string{"A"}, "A"},
	}

	for _, test := range tests {
		assert.Equal(t, SiteTitle(test.get...), test.want)
	}
}
