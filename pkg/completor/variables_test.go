package completor_test

import (
	complitor "ahmedash95/php-lsp-server/pkg/completor"
	"ahmedash95/php-lsp-server/pkg/lsp"
	"ahmedash95/php-lsp-server/pkg/treesitter"
	"testing"
)

func TestVariablesComplitor(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		position lsp.Position
		expected []complitor.Match
	}{
		{
			name: "Complete variables in a file",
			code: `<?php
$name = "Ahmed";
$num = 5;
echo $n
`,
			position: lsp.Position{
				Line:      3,
				Character: 6,
			},
			expected: []complitor.Match{
				{
					Text: "$name",
				},
				{
					Text: "$num",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node := treesitter.GetNodeAtPosition(test.code, test.position)

			complitor := complitor.VariablesCompletor{}
			matches := complitor.Complete(node)

			if len(matches) != len(test.expected) {
				t.Errorf("Expected %d matches, but got %d", len(test.expected), len(matches))
			}

			for i, match := range matches {
				if match.Text != test.expected[i].Text {
					t.Errorf("Expected %s, but got %s", test.expected[i].Text, match.Text)
				}
			}
		})
	}
}
