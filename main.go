package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/yeefea/gointerface/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

var (
	inputFile  = flag.String("i", "", "input file")
	outputFile = flag.String("o", "", "output file")
)

func main() {

	var fileInfoList []*parser.SourceFileInfo

	flag.Parse()
	if inputFile == nil || *inputFile == "" || *inputFile == "-" {
		raw, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input := antlr.NewInputStream(string(raw))
		fileInfo, err := analyze(input)
		if err != nil {
			fmt.Println(123123)
			panic(err)
		}
		fileInfoList = []*parser.SourceFileInfo{fileInfo}

	} else {

		info, err := os.Stat(*inputFile)
		if err != nil {
			panic(err)
		}
		if info.IsDir() { // is package
			files, err := ioutil.ReadDir(*inputFile)
			if err != nil {
				panic(err)
			}
			if len(files) == 0 {
				panic("package is empty")
			}

			for _, f := range files {
				if f.IsDir() {
					continue
				}
				if filepath.Ext(f.Name()) != ".go" {
					continue
				}
				filename := filepath.Join(*inputFile, f.Name())
				input, err := antlr.NewFileStream(filename)
				if err != nil {
					panic(err)
				}
				fileInfo, err := analyze(input)
				if err != nil {
					panic(err)
				}
				fileInfoList = append(fileInfoList, fileInfo)
			}
		} else { // is go file
			input, err := antlr.NewFileStream(*inputFile)
			if err != nil {
				panic(err)
			}
			fileInfo, err := analyze(input)
			if err != nil {
				panic(err)
			}
			fileInfoList = []*parser.SourceFileInfo{fileInfo}

		}
	}

	gen := parser.InterfaceGenerator{Files: fileInfoList}
	code, err := gen.GenerateCode()
	if err != nil {
		panic(err)
	}
	if outputFile == nil || *outputFile == "" {
		fmt.Println(code)
	} else {
		f, err := os.Create(*outputFile)
		if err != nil {
			panic(err)
		}
		// remember to close the file
		defer f.Close()
		_, err = f.WriteString(code)
		if err != nil {
			panic(err)
		}
	}

}

func analyze(input antlr.CharStream) (*parser.SourceFileInfo, error) {
	lexer := parser.NewGoLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)
	p := parser.NewGoParser(stream)
	p.BuildParseTrees = true
	p.AddErrorListener(parser.NewErrorListener())
	tree := p.SourceFile()
	listener := parser.NewMethodListener()
	walker := antlr.ParseTreeWalkerDefault
	walker.Walk(listener, tree)
	return listener.GetResult(), nil
}
