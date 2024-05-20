package completor

import (
	"ahmedash95/php-lsp-server/pkg/logger"
	"ahmedash95/php-lsp-server/pkg/lsp"
	"ahmedash95/php-lsp-server/pkg/treesitter"

	sitter "github.com/smacker/go-tree-sitter"
)

type Match struct {
	Text string
	Kind int
}

type CompletorInterface interface {
	CanComplete(doc *treesitter.TextDocumentItem, node *sitter.Node) bool
	Complete(doc *treesitter.TextDocumentItem, node *sitter.Node) []Match
}

type Completor struct {
	registers []CompletorInterface
}

func NewCompletor() Completor {
	return Completor{
		registers: []CompletorInterface{
			&VariablesCompletor{},
			&ObjectAccessing{},
			&InstanceAccess{},
		},
	}
}

func (c *Completor) GetCompletions(doc *treesitter.TextDocumentItem, pos lsp.Position) []Match {
	node := treesitter.GetNodeAtPosition(doc.Text, pos)
	if node == nil {
		logger.GetLogger().Printf("No node found at position: %v", pos)
		panic("No node found at position")
	}

	if len(c.registers) == 0 {
		panic("No completor registered")
	}

	var matches []Match
	for _, register := range c.registers {
		if register.CanComplete(doc, node) {
			matches = append(matches, register.Complete(doc, node)...)
		}
	}

	return matches
}
