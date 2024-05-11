# PHP LSP

This is a language server for PHP that adheres to the [Language Server Protocol (LSP)](https://microsoft.github.io/language-server-protocol/). 

The LSP is implemented in Go and is designed to be fast and efficient. and as Go does not have a good PHP parser library, I decided to use [Tree-sitter](https://tree-sitter.github.io/tree-sitter/) to parse PHP code. I'm not sure if this approach is good enough to build a powerful language server, but I'm trying to make it work and it seems to be working well so far.

## Features
- [x] Text Document Sync (full sync)
- [x] Document Symbols
- [x] Workspace Symbols
    - [ ] Support for query param
- [ ] Completion
- [ ] Code Actions
- [ ] Hover
- [ ] Signature Help
- [ ] Goto Definition
- [ ] Find References
- [ ] Diagnostics
- [ ] Formatting

## Installation
TBD

## Build from source


```bash
make build
```
Then you can run the server with the following command:

```bash
./php-lsp-server
```

## Testing
```bash
make test
```
