package lsp

type DocumentSymbolRequest struct {
	Request
	Params DocumentSymbolParams `json:"params"`
}

type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type DocumentSymbolResponse struct {
	Response
	Result []DocumentSymbol `json:"result"`
}

type DocumentSymbol struct {
	Name           string           `json:"name"`
	Detail         string           `json:"detail,omitempty"`
	Kind           int              `json:"kind"`
	Deprecated     bool             `json:"deprecated,omitempty"`
	Range          Range            `json:"range"`
	SelectionRange Range            `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}
