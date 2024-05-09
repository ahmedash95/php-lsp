package main

import (
	"ahmedash95/php-lsp-server/pkg/logger"
	"ahmedash95/php-lsp-server/pkg/lsp"
	"ahmedash95/php-lsp-server/pkg/rpc"
	"ahmedash95/php-lsp-server/pkg/treesitter"
	"bufio"
	"encoding/json"
	"io"
	"os"
)

func main() {
	l := logger.CreateLogFile("/tmp/php-lsp-server.log")
	l.Println("Starting PHP LSP Server")

	logger.SetLogger(l)

	state := treesitter.NewState()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()

		// Using a deferred function to recover from panics within this iteration of the loop
		defer func() {
			if r := recover(); r != nil {
				l.Printf("Recovered in f: %v\n", r)
			}
		}()

		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			l.Println("Error decoding message: ", err)
			continue
		}

		handleMessage(os.Stdout, state, method, contents)
	}
}

func handleMessage(writer io.Writer, state treesitter.State, method string, contents []byte) {
	logger := logger.GetLogger()
	logger.Printf("Recived message: [%s]", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Println("Error unmarshalling initialize request: ", err)
			return
		}

		logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		message := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(message)

		writer := os.Stdout
		writer.Write([]byte(reply))
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Println("Error unmarshalling didOpen request: ", err)
			return
		}

		state.Put(request.Params.TextDocument.Uri, request.Params.TextDocument.Text)
		logger.Printf("Opened file: %s", request.Params.TextDocument.Uri)

	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Println("Error unmarshalling didChange request: ", err)
			return
		}

		state.Update(request.Params.TextDocument.Uri, request.Params.ContentChanges)
		logger.Printf("Changed file: %s", request.Params.TextDocument.Uri)

	case "textDocument/documentSymbol":
		var request lsp.DocumentSymbolRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Println("Error unmarshalling documentSymbol request: ", err)
			return
		}

		response := state.TextDocumentDocumentSymbols(request.ID, request.Params.TextDocument.Uri)
		writeResponse(writer, response)
	}
}

func writeResponse(writer io.Writer, message any) {
	reply := rpc.EncodeMessage(message)
	writer.Write([]byte(reply))
}
