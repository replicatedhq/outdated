package outdated

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNamespaceExcluded(t *testing.T) {
	tests := []struct {
		name       string
		namespace  string
		exclusions []string
		expected   bool
	}{
		{
			name:       "no exclusions",
			namespace:  "foo",
			exclusions: []string{},
			expected:   false,
		},
		{
			name:       "exact match",
			namespace:  "foo",
			exclusions: []string{"foo"},
			expected:   true,
		},
		{
			name:       "exact match in list",
			namespace:  "foo",
			exclusions: []string{"one", "foo", "two"},
			expected:   true,
		},
		{
			name:       "not in list",
			namespace:  "foo",
			exclusions: []string{"one", "two"},
			expected:   false,
		},
		{
			name:       "wildcard match",
			namespace:  "foo_one",
			exclusions: []string{"foo_*"},
			expected:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := isNamespaceExcluded(test.namespace, test.exclusions)
			assert.Equal(t, test.expected, actual)
		})
	}
}
