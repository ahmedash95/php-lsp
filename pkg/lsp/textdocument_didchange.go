package lsp

type DidChangeTextDocumentNotification struct {
	Notification
	Params DidChangeTextDocumentParamsParams `json:"params"`
}

type DidChangeTextDocumentParamsParams struct {
	TextDocument   TextDocumentIdentifier           `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength int    `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}
