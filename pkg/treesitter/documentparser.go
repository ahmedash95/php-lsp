package treesitter

import (
	"ahmedash95/php-lsp-server/pkg/lsp"
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/php"
)

const (
	Kind_File          = 1
	Kind_Module        = 2
	Kind_Namespace     = 3
	Kind_Package       = 4
	Kind_Class         = 5
	Kind_Method        = 6
	Kind_Property      = 7
	Kind_Field         = 8
	Kind_Constructor   = 9
	Kind_Enum          = 10
	Kind_Interface     = 11
	Kind_Function      = 12
	Kind_Variable      = 13
	Kind_Constant      = 14
	Kind_String        = 15
	Kind_Number        = 16
	Kind_Boolean       = 17
	Kind_Array         = 18
	Kind_Object        = 19
	Kind_Key           = 20
	Kind_Null          = 21
	Kind_EnumMember    = 22
	Kind_Struct        = 23
	Kind_Event         = 24
	Kind_Operator      = 25
	Kind_TypeParameter = 26
)

const (
	Kind_Method_Label        = "Method"
	Kind_Function_Label      = "Function"
	Kind_Constructor_Label   = "Constructor"
	Kind_Field_Label         = "Field"
	Kind_Variable_Label      = "Variable"
	Kind_Class_Label         = "Class"
	Kind_Interface_Label     = "Interface"
	Kind_Module_Label        = "Module"
	Kind_Property_Label      = "Property"
	Kind_Unit_Label          = "Unit"
	Kind_Value_Label         = "Value"
	Kind_Enum_Label          = "Enum"
	Kind_Keyword_Label       = "Keyword"
	Kind_Snippet_Label       = "Snippet"
	Kind_Color_Label         = "Color"
	Kind_File_Label          = "File"
	Kind_Reference_Label     = "Reference"
	Kind_Folder_Label        = "Folder"
	Kind_EnumMember_Label    = "EnumMember"
	Kind_Constant_Label      = "Constant"
	Kind_Struct_Label        = "Struct"
	Kind_Event_Label         = "Event"
	Kind_Operator_Label      = "Operator"
	Kind_TypeParameter_Label = "TypeParameter"
)

var Kind_Labels = map[int]string{
	Kind_Method:        Kind_Method_Label,
	Kind_Function:      Kind_Function_Label,
	Kind_Constructor:   Kind_Constructor_Label,
	Kind_Field:         Kind_Field_Label,
	Kind_Variable:      Kind_Variable_Label,
	Kind_Class:         Kind_Class_Label,
	Kind_Interface:     Kind_Interface_Label,
	Kind_Module:        Kind_Module_Label,
	Kind_Property:      Kind_Property_Label,
	Kind_Enum:          Kind_Enum_Label,
	Kind_File:          Kind_File_Label,
	Kind_EnumMember:    Kind_EnumMember_Label,
	Kind_Constant:      Kind_Constant_Label,
	Kind_Struct:        Kind_Struct_Label,
	Kind_Event:         Kind_Event_Label,
	Kind_Operator:      Kind_Operator_Label,
	Kind_TypeParameter: Kind_TypeParameter_Label,
}

type TextDocumentItem struct {
	Uri             string               `json:"uri"`
	LanguageId      string               `json:"languageId"`
	Version         int                  `json:"version"`
	Text            string               `json:"text"`
	DocumentSymbols []lsp.DocumentSymbol `json:"documentSymbols"`
}

func ParseDocument(content string) (*sitter.Tree, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(php.GetLanguage())

	return parser.ParseCtx(context.Background(), nil, []byte(content))
}

func GetNodeText(content string, node *sitter.Node) string {
	return content[node.StartByte():node.EndByte()]
}

func FindNodesByType(node *sitter.Node, nodeType string) []*sitter.Node {
	nodes := []*sitter.Node{}

	if node.Type() == nodeType {
		nodes = append(nodes, node)
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		nodes = append(nodes, FindNodesByType(child, nodeType)...)
	}

	return nodes
}
