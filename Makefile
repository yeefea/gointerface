.PHONY: build

PKG=github.com/yeefea/gointerface

OUTPUT_DIR=bin

GO=go

build:
	$(GO) build -o $(OUTPUT_DIR)/gointerface $(PKG)
