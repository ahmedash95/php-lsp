package workspace

import (
	"ahmedash95/php-lsp-server/internal/util"
	"ahmedash95/php-lsp-server/pkg/completor"
	"ahmedash95/php-lsp-server/pkg/logger"
	"ahmedash95/php-lsp-server/pkg/lsp"
	"ahmedash95/php-lsp-server/pkg/treesitter"
	workspacescanner "ahmedash95/php-lsp-server/pkg/workspace_scanner"
	"fmt"

	"github.com/sahilm/fuzzy"
)

type Workspace struct {
	Uris     map[string]*treesitter.TextDocumentItem
	RootPath string
}

func NewWorkspace(rootpath string) *Workspace {
	return &Workspace{
		Uris:     make(map[string]*treesitter.TextDocumentItem),
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

func (s *Workspace) Get(uri string) *treesitter.TextDocumentItem {
	return s.Uris[uri]
}

func (s *Workspace) Put(uri string, content string) {
	s.Uris[uri] = &treesitter.TextDocumentItem{
		Uri:        uri,
		LanguageId: "php",
		Version:    1,
		Text:       content,
	}

	s.FetchDocumentSymbols(uri, content)
}

// @todo implement incremental parsing
func (s *Workspace) Update(uri string, contentChanges []lsp.TextDocumentContentChangeEvent) {
	s.Uris[uri].Text = contentChanges[0].Text

	s.FetchDocumentSymbols(uri, contentChanges[0].Text)
}

func symbolToLspSymbol(symbol *treesitter.Symbol) lsp.DocumentSymbol {
	childs := []lsp.DocumentSymbol{}
	for _, child := range symbol.Children {
		childs = append(childs, symbolToLspSymbol(&child))
	}

	return lsp.DocumentSymbol{
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
		Children: childs,
	}
}

func (s *Workspace) FetchDocumentSymbols(uri string, text string) {
	symbols := treesitter.GetDocumentSymbols(text)

	items := []lsp.DocumentSymbol{}
	for _, symbol := range symbols {
		items = append(items, symbolToLspSymbol(&symbol))
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

func (s *Workspace) TextDocumentCompletion(id int, textDocumentPosition lsp.TextDocumentPositionParams) lsp.CompletionResponse {

	pos := textDocumentPosition.Position
	pos.Character = pos.Character - 1

	doc := s.Get(textDocumentPosition.TextDocument.Uri)

	completor := completor.NewCompletor()
	matches := completor.GetCompletions(doc, pos)

	completions := []lsp.CompletionItem{}
	for _, match := range matches {
		completions = append(completions, lsp.CompletionItem{
			Label: match.Text,
			Kind:  match.Kind,
		})
	}

	response := lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: completions,
	}

	return response
}
