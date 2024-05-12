package treesitter

import (
	"ahmedash95/php-lsp-server/internal/util"
	"ahmedash95/php-lsp-server/pkg/logger"
	"ahmedash95/php-lsp-server/pkg/lsp"
	workspacescanner "ahmedash95/php-lsp-server/pkg/workspace_scanner"
	"fmt"

	"github.com/sahilm/fuzzy"
	sitter "github.com/smacker/go-tree-sitter"
)

type TextDocumentItem struct {
	Uri             string               `json:"uri"`
	LanguageId      string               `json:"languageId"`
	Version         int                  `json:"version"`
	Text            string               `json:"text"`
	AST             *sitter.Tree         `json:"ast"`
	DocumentSymbols []lsp.DocumentSymbol `json:"documentSymbols"`
}

type Workspace struct {
	Uris     map[string]*TextDocumentItem
	RootPath string
}

func NewWorkspace(rootpath string) *Workspace {
	return &Workspace{
		Uris:     make(map[string]*TextDocumentItem),
		RootPath: rootpath,
	}
}

func (s *Workspace) StartIndex(update func(path string, percent int), end func()) {
	logger.GetLogger().Printf("Indexing workspace: %s", s.RootPath)

	scanner := workspacescanner.NewScanner(s.RootPath)
	files := scanner.Scan([]string{".php"})

	defer end()

	// @todo parse file in parallel to speed up indexing
	for i, file := range files {
		uri := fmt.Sprintf("file://%s/%s", s.RootPath, file)

		if _, ok := s.Uris[uri]; ok {
			continue
		}

		logger.GetLogger().Printf("Indexing file: %s", uri)

		content := scanner.GetFileContent(file)
		s.Put(uri, content)

		update(file, util.CalculatePercentage(i, len(files)))
	}
}

func (s *Workspace) Get(uri string) *TextDocumentItem {
	return s.Uris[uri]
}

func (s *Workspace) Put(uri string, content string) {
	s.Uris[uri] = &TextDocumentItem{
		Uri:        uri,
		LanguageId: "php",
		Version:    1,
		Text:       content,
	}

	// document is parsed and AST is generated
	var err error
	s.Uris[uri].AST, err = ParseDocument(content)
	if err != nil {
		logger.GetLogger().Printf("Error parsing document: %s", err)
	}

	s.FetchDocumentSymbols(uri, content)
}

// @todo implement incremental parsing
func (s *Workspace) Update(uri string, contentChanges []lsp.TextDocumentContentChangeEvent) {
	s.Uris[uri].Text = contentChanges[0].Text
	// document is re-parsed and AST is regenerated
	var err error
	s.Uris[uri].AST, err = ParseDocument(contentChanges[0].Text)
	if err != nil {
		logger.GetLogger().Printf("Error parsing document: %s", err)
	}
}

func (s *Workspace) FetchDocumentSymbols(uri string, text string) {
	symbols := GetDocumentSymbols(text)

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

	s.Uris[uri].DocumentSymbols = items
}

func (s *Workspace) TextDocumentDocumentSymbols(id int, uri string) lsp.DocumentSymbolResponse {
	response := lsp.DocumentSymbolResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: s.Uris[uri].DocumentSymbols,
	}

	logger.GetLogger().Printf("Document symbol response: %v", response)

	return response
}

type wsSymbols []struct {
	URI    string             `json:"uri"`
	Symbol lsp.DocumentSymbol `json:"symbol"`
}

// github.com/sahilm/fuzzy
func (s wsSymbols) Len() int {
	return len(s)
}

func (s wsSymbols) String(i int) string {
	return s[i].Symbol.Name
}

func (s *Workspace) WorkspaceSymbols(id int, query string) lsp.WorkspaceSymbolResponse {
	symbols := []lsp.WorkSpaceSymbol{}

	urisSymbols := wsSymbols{}
	for _, item := range s.Uris {
		for _, symbol := range item.DocumentSymbols {
			urisSymbols = append(urisSymbols, struct {
				URI    string             `json:"uri"`
				Symbol lsp.DocumentSymbol `json:"symbol"`
			}{
				URI:    item.Uri,
				Symbol: symbol,
			})
		}
	}

	for _, r := range fuzzy.FindFrom(query, urisSymbols) {
		ds := urisSymbols[r.Index]
		symbols = append(symbols, lsp.WorkSpaceSymbol{
			Name: ds.Symbol.Name,
			Kind: ds.Symbol.Kind,
			Location: lsp.Location{
				URI: ds.URI,
				Range: lsp.Range{
					Start: lsp.Position{
						Line:      ds.Symbol.Range.Start.Line,
						Character: ds.Symbol.Range.Start.Character,
					},
					End: lsp.Position{
						Line:      ds.Symbol.Range.End.Line,
						Character: ds.Symbol.Range.End.Character,
					},
				},
			},
		})
	}

	response := lsp.WorkspaceSymbolResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: symbols,
	}

	return response
}
