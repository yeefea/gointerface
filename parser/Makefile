include .env

.PHONY: lexer parser all

DEFAULT: all

TARGET_LANG=Go
FLAGES=-no-visitor

LEXER_SRC=GoLexer.g4
PARSER_SRC=GoParser.g4

all: lexer parser

lexer:
	$(ANTLR) -Dlanguage=$(TARGET_LANG) $(FLAGES) $(LEXER_SRC)

parser:
	$(ANTLR) -Dlanguage=$(TARGET_LANG) $(FLAGES) $(PARSER_SRC)
