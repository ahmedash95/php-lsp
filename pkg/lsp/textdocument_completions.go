package lsp

const (
	Symbol_Kind_Text          = 1
	Symbol_Kind_Method        = 2
	Symbol_Kind_Function      = 3
	Symbol_Kind_Constructor   = 4
	Symbol_Kind_Field         = 5
	Symbol_Kind_Variable      = 6
	Symbol_Kind_Class         = 7
	Symbol_Kind_Interface     = 8
	Symbol_Kind_Module        = 9
	Symbol_Kind_Property      = 10
	Symbol_Kind_Unit          = 11
	Symbol_Kind_Value         = 12
	Symbol_Kind_Enum          = 13
	Symbol_Kind_Keyword       = 14
	Symbol_Kind_Snippet       = 15
	Symbol_Kind_Color         = 16
	Symbol_Kind_File          = 17
	Symbol_Kind_Reference     = 18
	Symbol_Kind_Folder        = 19
	Symbol_Kind_EnumMember    = 20
	Symbol_Kind_Constant      = 21
	Symbol_Kind_Struct        = 22
	Symbol_Kind_Event         = 23
	Symbol_Kind_Operator      = 24
	Symbol_Kind_TypeParameter = 25
)

type CompletionRequest struct {
	Request
	Params CompletionParams `json:"params"`
}

type CompletionParams struct {
	TextDocumentPositionParams
}

type CompletionResponse struct {
	Response
	Result []CompletionItem `json:"result"`
}

type CompletionItem struct {
	Label         string `json:"label"`
	Kind          int    `json:"kind"`
	Detail        string `json:"detail,omitempty"`
	Documentation string `json:"documentation,omitempty"`
}
