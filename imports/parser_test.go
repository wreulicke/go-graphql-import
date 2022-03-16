package imports_test

import (
	"testing"

	"github.com/wreulicke/go-graphql-import/imports"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    `# import X from "x.graphql"`,
			expected: []string{"X"},
		},
		{
			input:    `# import X, Y from "x.graphql"`,
			expected: []string{"X", "Y"},
		},
		{
			input:    `# import X, Mutation.* from "x.graphql"`,
			expected: []string{"X", "Mutation"},
		},
		{
			input:    `# import X, Y, Z from "x.graphql"`,
			expected: []string{"X", "Y", "Z"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			p := imports.NewParser(tt.input)
			r, errors := p.Parse()
			if len(errors) > 0 {
				t.Error("errors is found", errors)
				return
			}
			if r.From == "x.graphql" {
				t.Error("from is unexpected")
			}
			if len(r.Imports) == 0 {
				t.Error("imports should has positve length")
			}
			for index, e := range tt.expected {
				x := r.Imports[index]
				if x.Name != e {
					t.Errorf("expected %s found %s", e, x.Name)
				}
			}
		})
	}
}
