package treesitter

import (
	sitter "github.com/smacker/go-tree-sitter"
)

type Position struct {
	LineStart   uint32
	LineEnd     uint32
	OffsetStart uint32
	OffsetEnd   uint32
}

type Symbol struct {
	Name     string
	Kind     uint32
	Position Position
	Children []Symbol
}

func GetDocumentSymbols(content string) []Symbol {
	tree, err := ParseDocument(content)
	if err != nil {
		return []Symbol{}
	}

	var symbols []Symbol

	WalkTree(content, tree.RootNode(), &symbols)
	return symbols
}

func WalkTree(content string, node *sitter.Node, symbols *[]Symbol) {

	switch node.Type() {
	case "variable_name":
		n := node.Child(1)
		*symbols = append(*symbols, getSymbolFromNode(content, Kind_Variable, n))
	case "function_definition":
		n := node.Child(1)
		*symbols = append(*symbols, getSymbolFromNode(content, Kind_Function, n))
	case "class_declaration":
		n := findNodeOfType(node, "name")
		*symbols = append(*symbols, getSymbolFromNode(content, Kind_Class, n))
	case "base_clause":
		var kind uint32
		if node.Parent().Type() == "interface_declaration" {
			kind = Kind_Interface
		} else {
			kind = Kind_Class
		}
		for i := 0; i < int(node.NamedChildCount()); i++ {
			*symbols = append(*symbols, getSymbolFromNode(content, kind, node.NamedChild(i)))
		}
	case "class_interface_clause":
		for i := 0; i < int(node.NamedChildCount()); i++ {
			n := node.NamedChild(i)
			*symbols = append(*symbols, getSymbolFromNode(content, Kind_Interface, n))
		}
	case "interface_declaration":
		n := node.Child(1)
		*symbols = append(*symbols, getSymbolFromNode(content, Kind_Interface, n))
	case "trait_declaration":
		n := node.Child(1)
		*symbols = append(*symbols, getSymbolFromNode(content, Kind_Class, n))
	case "use_declaration":
		for i := 0; i < int(node.NamedChildCount()); i++ {
			n := node.NamedChild(i)
			*symbols = append(*symbols, getSymbolFromNode(content, Kind_Class, n))
		}
	case "method_declaration":
		n := findNodeOfType(node, "name")
		*symbols = append(*symbols, getSymbolFromNode(content, Kind_Method, n))

	case "property_declaration":
		n := findNodeOfType(node, "name")
		*symbols = append(*symbols, getSymbolFromNode(content, Kind_Property, n))
		return

	case "const_declaration":
		n := findNodeOfType(node, "name")
		*symbols = append(*symbols, getSymbolFromNode(content, Kind_Constant, n))
		return

		// disabled for now because it can't handle deifine($key, $value) calls
		// case "function_call_expression":
		// 	if GetNodeText(content, node.Child(0)) == "define" {
		// 		n := findNodeOfType(node.Child(1), "string")
		// 		*symbols = append(*symbols, getSymbolFromNode(content, Kind_Constant, n))
		//
		// 	}
	}

	var childrenSymbols []Symbol

	for child := node.Child(0); child != nil; child = child.NextSibling() {
		WalkTree(content, child, &childrenSymbols)
	}

	canHaveChildren := node.Type() == "class_declaration" || node.Type() == "interface_declaration" || node.Type() == "trait_declaration" || node.Type() == "function_definition"

	if canHaveChildren {
		parent := &(*symbols)[len(*symbols)-1]
		parent.Children = childrenSymbols
	} else {
		*symbols = append(*symbols, childrenSymbols...)
	}
}

func getSymbolFromNode(content string, kind uint32, node *sitter.Node) Symbol {
	return Symbol{
		Name: GetNodeText(content, node),
		Kind: kind,
		Position: Position{
			LineStart:   node.StartPoint().Row,
			LineEnd:     node.EndPoint().Row,
			OffsetStart: node.StartPoint().Column,
			OffsetEnd:   node.EndPoint().Column,
		},
	}
}

func findNodeOfType(node *sitter.Node, nodeType string) *sitter.Node {
	if node == nil {
		return nil
	}

	if node.Type() == nodeType {
		return node
	}

	for child := node.Child(0); child != nil; child = child.NextSibling() {
		n := findNodeOfType(child, nodeType)
		if n != nil {
			return n
		}
	}

	return nil
}
