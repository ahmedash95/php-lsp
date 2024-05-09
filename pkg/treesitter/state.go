package treesitter

import (
	"ahmedash95/php-lsp-server/pkg/logger"
	"ahmedash95/php-lsp-server/pkg/lsp"
)

type TextDocumentItem struct {
	Uri        string `json:"uri"`
	LanguageId string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

type State struct {
	Uris map[string]TextDocumentItem
}

func NewState() State {
	return State{
		Uris: make(map[string]TextDocumentItem),
	}
}

func (s *State) Get(uri string) TextDocumentItem {
	return s.Uris[uri]
}

func (s *State) Put(uri string, content string) {
	s.Uris[uri] = TextDocumentItem{
		Uri:        uri,
		LanguageId: "php",
		Version:    1,
		Text:       content,
	}
}

func (s *State) Update(uri string, contentChanges []lsp.TextDocumentContentChangeEvent) {
	s.Uris[uri] = TextDocumentItem{
		Uri:        uri,
		LanguageId: "php",
		Version:    1,
		Text:       contentChanges[0].Text,
	}
}

func (s *State) TextDocumentDocumentSymbols(id int, uri string) lsp.DocumentSymbolResponse {
	logger.GetLogger().Printf("Document symbol for textDocument: %s \n content: %s", uri, s.Get(uri).Text)
	symbols := GetDocumentSymbols(s.Get(uri).Text)

	items := []lsp.DocumentSymbol{}

	for _, symbol := range symbols {
		items = append(items, lsp.DocumentSymbol{
			Name: symbol.Name,
			Kind: int(symbol.Kind),
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      int(symbol.Position.LineStart),
					Character: int(symbol.Position.OffsetStart),
				},
				End: lsp.Position{
					Line:      int(symbol.Position.LineEnd),
					Character: int(symbol.Position.OffsetEnd),
				},
			},
			SelectionRange: lsp.Range{
				Start: lsp.Position{
					Line:      int(symbol.Position.LineStart),
					Character: int(symbol.Position.OffsetStart),
				},
				End: lsp.Position{
					Line:      int(symbol.Position.LineEnd),
					Character: int(symbol.Position.OffsetEnd),
				},
			},
		})
	}

	response := lsp.DocumentSymbolResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: items,
	}

	logger.GetLogger().Printf("Document symbol response: %v", response)

	return response
}
