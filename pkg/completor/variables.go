package completor

import (
	"ahmedash95/php-lsp-server/pkg/logger"
	"ahmedash95/php-lsp-server/pkg/lsp"
	"ahmedash95/php-lsp-server/pkg/treesitter"

	sitter "github.com/smacker/go-tree-sitter"
)

type VariablesCompletor struct{}

func (com *VariablesCompletor) CanComplete(doc *treesitter.TextDocumentItem, node *sitter.Node) bool {
	if node.Type() == "name" && node.Parent().Type() == "variable_name" {
		return true
	}

	if node.Type() == "variable" {
		return true
	}

	logger.GetLogger().Printf("No variable at -  node type: %s and parent %s", node.Type(), node.Parent().Type())

	return false
}

func (com *VariablesCompletor) Complete(doc *treesitter.TextDocumentItem, node *sitter.Node) []Match {
	var matches []Match

	parentScope := findParentScope(node)

	logger.GetLogger().Printf("Scope content: %s", doc.Text[parentScope.StartByte():parentScope.EndByte()])

	// get all variables from parent scope and append to matches
	findInSymbols(node, &matches, parentScope, doc.DocumentSymbols)
	return matches
}

func findInSymbols(node *sitter.Node, matches *[]Match, parentScope *sitter.Node, symbols []lsp.DocumentSymbol) {
	for _, symbol := range symbols {

		findInSymbols(node, matches, parentScope, symbol.Children)

		if uint32(symbol.Range.Start.Line) < parentScope.StartPoint().Row || uint32(symbol.Range.End.Line) > parentScope.EndPoint().Row || node.EndPoint().Row < uint32(symbol.Range.Start.Line) {
			logger.GetLogger().Printf("Symbol %s is out of scope", symbol.Name)
			continue
		}

		if symbol.Kind == treesitter.Kind_Variable {
			logger.GetLogger().Printf("Variable found: %s of kind %d", symbol.Name, symbol.Kind)
			*matches = append(*matches, Match{Text: symbol.Name, Kind: lsp.Symbol_Kind_Variable})
		}

	}
}

func findParentScope(node *sitter.Node) *sitter.Node {
	// iterate over parent chain to one of scopes body or program
	for node != nil {
		if node.Type() == "function_declaration" || node.Type() == "class_declaration" || node.Type() == "interface_declaration" || node.Type() == "compound_statement" || node.Type() == "program" {
			return node
		}

		node = node.Parent()
	}

	return node
}
