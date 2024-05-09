# build go cmd/lsp_php_server/main.go

.PHONY: build_langauge_server
build_langauge_server:
	go build -o php-lsp-server cmd/php-lsp-server/languageserver.go


.PHONY: build
build: build_langauge_server

.PHONY: test
test:
	go test ./...

