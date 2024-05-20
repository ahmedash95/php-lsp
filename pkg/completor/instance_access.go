package completor

import (
	"ahmedash95/php-lsp-server/pkg/logger"
	"ahmedash95/php-lsp-server/pkg/lsp"
	"ahmedash95/php-lsp-server/pkg/treesitter"

	sitter "github.com/smacker/go-tree-sitter"
)

type InstanceAccess struct{}

func (com *InstanceAccess) CanComplete(doc *treesitter.TextDocumentItem, node *sitter.Node) bool {
	firstNamedChild := node.Parent().NamedChild(0)
	if firstNamedChild == nil {
		return false
	}

	firstNamedChildIsThis := treesitter.GetNodeText(doc.Text, firstNamedChild) != "$this"

	logger.GetLogger().Printf("Named Child content: %s", treesitter.GetNodeText(doc.Text, firstNamedChild))

	if node.Type() == "name" && node.Parent().Type() == "member_access_expression" && firstNamedChildIsThis {
		return true
	}

	return false
}

func (com *InstanceAccess) Complete(doc *treesitter.TextDocumentItem, node *sitter.Node) []Match {
	// we need to first find the object name
	name := com.findObjectName(doc, node)
	if name == "" {
		logger.GetLogger().Printf("Failed to extract object name")
		return []Match{}
	}
	// then find the class of that object
	className := com.findClassName(doc, name)
	if className == "" {
		logger.GetLogger().Printf("Failed to extract class name for object: %s", name)
	}
	// then find in doc symbols the class and get all properties and methods
	var classNode *lsp.DocumentSymbol
	for _, symbol := range doc.DocumentSymbols {
		if symbol.Name == className && symbol.Kind == treesitter.Kind_Class {
			classNode = &symbol
			break
		}
	}

	if classNode == nil {
		logger.GetLogger().Printf("Failed to find class node for class: %s", className)
		return []Match{}
	}

	var matches []Match
	com.findInSymbols(&matches, classNode.Children)
	return matches
}

func (com *InstanceAccess) findInSymbols(matches *[]Match, symbols []lsp.DocumentSymbol) {
	for _, symbol := range symbols {
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

func (com *InstanceAccess) findObjectName(doc *treesitter.TextDocumentItem, node *sitter.Node) string {
	// iterate over parent chain to find the object name
	for node != nil {
		if node.Type() == "member_access_expression" {
			return treesitter.GetNodeText(doc.Text, node.NamedChild(0))
		}
		node = node.Parent()
	}
	return ""
}

func (com *InstanceAccess) findClassName(doc *treesitter.TextDocumentItem, objectName string) string {
	root, err := treesitter.ParseDocument(doc.Text)
	if err != nil {
		logger.GetLogger().Printf("Failed to parse document")
		return ""
	}

	// find node of assignment expression and if left variable name is the object name, then return the right side string
	nodes := treesitter.FindNodesByType(root.RootNode(), "assignment_expression")

	for _, n := range nodes {
		left := n.Child(0)
		right := n.Child(2)

		if left == nil || right == nil {
			continue
		}

		if treesitter.GetNodeText(doc.Text, left) == objectName {
			return treesitter.GetNodeText(doc.Text, right.Child(1))
		}
	}

	logger.GetLogger().Printf("Failed to find class name for object: %s", objectName)
	return ""
}
