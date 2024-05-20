package treesitter

import (
	"ahmedash95/php-lsp-server/pkg/logger"
	"ahmedash95/php-lsp-server/pkg/lsp"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

func LinePositionToIndex(content string, pos lsp.Position) int {
	// split content by newline
	lines := strings.Split(content, "\n")
	// iterate over each line and count the number of characters
	charCount := 0
	for i, line := range lines {
		if i == pos.Line {
			break
		}

		charCount += len(line)
	}
	// once we reach the line number, we can calculate the offset
	offset := charCount + pos.Character
	// then we get the index of the character in the string

	return offset
}

func GetNodeAtPosition(content string, pos lsp.Position) *sitter.Node {
	ast, err := ParseDocument(content)
	if err != nil {
		return nil
	}
	root := ast.RootNode()

	node := walkTreeToPosition(root, pos)

	if node == nil {
		return nil
	}

	fmt.Printf("GetNodeAtPosition: %s\n", node.Type())

	switch node.Type() {
	}

	str := GetNodeText(content, node)

	logger.GetLogger().Printf("Node of type %s at position: line %d at char %d for content %s", node.Type(), pos.Line, pos.Character, str)

	return node
}

// walkTreeToPosition finds the node at the specified position in the source code.
func walkTreeToPosition(node *sitter.Node, pos lsp.Position) *sitter.Node {
	// Convert LSP's 0-based position to Tree-sitter's 0-based index.
	targetRow := uint32(pos.Line)
	targetColumn := uint32(pos.Character)

	// Check if the target position is outside of the current node's range.
	if !nodeContainsPosition(node, targetRow, targetColumn) {
		return nil // The position is outside of this node.
	}

	// Recursively search children nodes.
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if foundNode := walkTreeToPosition(child, pos); foundNode != nil {
			return foundNode // Return the first matching child node.
		}
	}

	// If no matching child is found, and this node contains the position, return this node.
	return node
}

// nodeContainsPosition checks if the node contains the given position.
func nodeContainsPosition(node *sitter.Node, targetRow, targetColumn uint32) bool {
	startRow := node.StartPoint().Row
	startColumn := node.StartPoint().Column
	endRow := node.EndPoint().Row
	endColumn := node.EndPoint().Column

	// Check if the target position is within the node's vertical range.
	if targetRow < startRow || targetRow > endRow {
		return false
	}

	// If on the same start or end row, check the horizontal range.
	if (targetRow == startRow && targetColumn < startColumn) || (targetRow == endRow && targetColumn >= endColumn) {
		return false
	}

	return true
}
