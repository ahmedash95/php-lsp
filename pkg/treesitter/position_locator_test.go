package treesitter_test

import (
	"ahmedash95/php-lsp-server/pkg/lsp"
	"ahmedash95/php-lsp-server/pkg/treesitter"
	"fmt"
	"testing"
)

func TestLinePositionToIndex(t *testing.T) {

	var tests = []struct {
		content         string
		pos             lsp.Position
		expectedType    string
		expectedContent string
	}{
		{
			content: `<?php
$name = "Alice";
$ag
`,
			pos: lsp.Position{
				Line:      2,
				Character: 2,
			},
			expectedType:    "name",
			expectedContent: "ag",
		},
		{
			content: `<?php
				function hello() {
					echo $a
				}
				`,
			pos: lsp.Position{
				Line:      2,
				Character: 11,
			},
			expectedType:    "name",
			expectedContent: "a",
		}, {
			content: `<?php
				$obj = new stdClass();
				$obj->a
			`,
			pos: lsp.Position{
				Line:      2,
				Character: 10,
			},
			expectedType:    "name",
			expectedContent: "a",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test LinePositionToIndex for %s", test.expectedContent), func(t *testing.T) {

			content := test.content
			pos := test.pos
			node := treesitter.GetNodeAtPosition(content, pos)
			nodeContent := treesitter.GetNodeText(content, node)

			if node.Type() != test.expectedType {
				t.Errorf("Expected node type %s, got %s", test.expectedType, node.Type())
			}

			if nodeContent != test.expectedContent {
				t.Errorf("Expected node content %s, got %s", test.expectedContent, nodeContent)
			}
		})
	}
}
