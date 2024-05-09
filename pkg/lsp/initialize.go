package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
	RootPath   string      `json:"rootPath"` // is null if no folder is open
	RootUri    string      `json:"rootUri"`  // is null if no folder is open
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync       int            `json:"textDocumentSync"`
	CompletionProvider     map[string]any `json:"completionProvider"`
	DocumentSymbolProvider bool           `json:"documentSymbolProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:       1,   // Full sync
				CompletionProvider:     nil, // map[string]any{},
				DocumentSymbolProvider: true,
			},
			ServerInfo: ServerInfo{
				Name:    "PHP Language Server",
				Version: "0.0.1",
			},
		},
	}
}
