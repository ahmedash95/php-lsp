package lsp

type WorkspaceSymbolRequest struct {
	Request
	Params WorkspaceSymbolParams `json:"params"`
}

type WorkspaceSymbolParams struct {
	Query string `json:"query"`
}

type WorkspaceSymbolResponse struct {
	Response
	Result []WorkSpaceSymbol `json:"result"`
}

type WorkSpaceSymbol struct {
	Name     string   `json:"name"`
	Kind     int      `json:"kind"`
	Location Location `json:"location"`
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}
