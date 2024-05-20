package completor

import (
	"ahmedash95/php-lsp-server/pkg/logger"
	"ahmedash95/php-lsp-server/pkg/lsp"
	"ahmedash95/php-lsp-server/pkg/treesitter"

	sitter "github.com/smacker/go-tree-sitter"
)

type ObjectAccessing struct{}

func (com *ObjectAccessing) CanComplete(doc *treesitter.TextDocumentItem, node *sitter.Node) bool {
	firstNamedChild := node.Parent().NamedChild(0)
	if firstNamedChild == nil {
		return false
	}

	firstNamedChildIsThis := treesitter.GetNodeText(doc.Text, firstNamedChild) == "$this"

	logger.GetLogger().Printf("Named Child content: %s", treesitter.GetNodeText(doc.Text, firstNamedChild))

	if node.Type() == "name" && node.Parent().Type() == "member_access_expression" && firstNamedChildIsThis {
		return true
	}

	return false
}

func (com *ObjectAccessing) Complete(doc *treesitter.TextDocumentItem, node *sitter.Node) []Match {
	var matches []Match

	parentScope := com.findParentScope(node)

	logger.GetLogger().Printf("Scope content: %s", doc.Text[parentScope.StartByte():parentScope.EndByte()])

	// get all variables from parent scope and append to matches
	com.findInSymbols(node, &matches, parentScope, doc.DocumentSymbols)
	return matches
}

func (com *ObjectAccessing) findInSymbols(node *sitter.Node, matches *[]Match, parentScope *sitter.Node, symbols []lsp.DocumentSymbol) {
	for _, symbol := range symbols {

		com.findInSymbols(node, matches, parentScope, symbol.Children)

		if uint32(symbol.Range.Start.Line) < parentScope.StartPoint().Row || uint32(symbol.Range.End.Line) > parentScope.EndPoint().Row || node.EndPoint().Row < uint32(symbol.Range.Start.Line) {
			logger.GetLogger().Printf("Symbol %s is out of scope", symbol.Name)
			continue
		}

		if symbol.Kind == treesitter.Kind_Property {
			logger.GetLogger().Printf("Property found: %s of kind %d", symbol.Name, symbol.Kind)
			*matches = append(*matches, Match{Text: symbol.Name, Kind: lsp.Symbol_Kind_Property})
		}

		if symbol.Kind == treesitter.Kind_Method {
			logger.GetLogger().Printf("Method found: %s of kind %d", symbol.Name, symbol.Kind)
			*matches = append(*matches, Match{Text: symbol.Name, Kind: lsp.Symbol_Kind_Method})
		}
	}
}

func (com *ObjectAccessing) findParentScope(node *sitter.Node) *sitter.Node {
	// iterate over parent chain to one of scopes body or program
	for node != nil {
		if node.Type() == "class_declaration" {
			return node
		}

		node = node.Parent()
	}

	return nil
}
