package main

import (
	"ahmedash95/php-lsp-server/pkg/logger"
	"ahmedash95/php-lsp-server/pkg/lsp"
	"ahmedash95/php-lsp-server/pkg/rpc"
	"ahmedash95/php-lsp-server/pkg/treesitter"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
)

func main() {
	fmt.Println("Starting PHP LSP Server")

	l := logger.CreateLogFile("/tmp/php-lsp-server.log")
	l.Println("Starting PHP LSP Server")

	logger.SetLogger(l)

	workspace := treesitter.NewWorkspace("")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()

		// Using a deferred function to recover from panics within this iteration of the loop
		defer func() {
			if r := recover(); r != nil {
				l.Printf("Recovered in f: %v\n", r)
				// Print the stack trace
				fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			}
		}()

		logger.GetLogger().Printf("Recived request: [%s]", msg)
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			l.Println("Error decoding message: ", err)
			continue
		}

		handleMessage(os.Stdout, workspace, method, contents)
	}
}

func handleMessage(writer io.Writer, workspace *treesitter.Workspace, method string, contents []byte) {
	logger := logger.GetLogger()
	logger.Printf("Recived message: [%s]", method)

	if method != "initialize" && workspace == nil {
		logger.Println("Error: Initialize request must be sent first before any other request")
		return
	}

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Println("Error unmarshalling initialize request: ", err)
			return
		}

		logger.Printf("Initializing workspace: %s", request.Params.RootPath)
		workspace.RootPath = request.Params.RootPath

		message := lsp.NewInitializeResponse(request.ID)
		writeResponse(writer, message)

		progressStartRequest := lsp.CreateProgressBeginRequest("indexing", "Indexing workspace")
		writeResponse(writer, progressStartRequest)

		update := func(path string, percent int) {
			logger.Printf("Indexing file: %s", path)
			progressUpdateRequest := lsp.CreateProgressUpdateRequest(progressStartRequest.Params.Token, path, percent)
			writeResponse(writer, progressUpdateRequest)
		}
		workspace.StartIndex(update, func() {
			progressEndRequest := lsp.CreateProgressEndRequest(progressStartRequest.Params.Token, "Indexing complete")
			writeResponse(writer, progressEndRequest)
		})
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Println("Error unmarshalling didOpen request: ", err)
			return
		}

		workspace.Put(request.Params.TextDocument.Uri, request.Params.TextDocument.Text)
		logger.Printf("Opened file: %s", request.Params.TextDocument.Uri)

	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Println("Error unmarshalling didChange request: ", err)
			return
		}

		workspace.Update(request.Params.TextDocument.Uri, request.Params.ContentChanges)
		logger.Printf("Changed file: %s", request.Params.TextDocument.Uri)

	case "textDocument/documentSymbol":
		var request lsp.DocumentSymbolRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Println("Error unmarshalling documentSymbol request: ", err)
			return
		}

		response := workspace.TextDocumentDocumentSymbols(request.ID, request.Params.TextDocument.Uri)
		writeResponse(writer, response)
	case "workspace/symbol":
		var request lsp.WorkspaceSymbolRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Println("Error unmarshalling workspaceSymbol request: ", err)
			return
		}

		response := workspace.WorkspaceSymbols(request.ID, request.Params.Query)
		writeResponse(writer, response)

		logger.Printf("Returning workspace symbols response: %v", response)
	}
}

func writeResponse(writer io.Writer, message any) {
	reply := rpc.EncodeMessage(message)
	writer.Write([]byte(reply))

	logger.GetLogger().Printf("Sending response: [%s]", reply)
}
