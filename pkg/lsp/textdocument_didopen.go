package lsp

type DidOpenTextDocumentNotification struct {
	Notification
	Params DidOpenTextDocumentParamsParams `json:"params"`
}

type DidOpenTextDocumentParamsParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}
