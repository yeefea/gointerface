package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/yeefea/gointerface/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

var (
	inputFile     string
	outputFile    string
	types         string
	private       bool
	interestTypes map[string]struct{}
)

func init() {
	flag.StringVar(&inputFile, "i", "", "Input file or package. By default, the program read from stdin.")
	flag.StringVar(&outputFile, "o", "", "Output file. By default, the program writes content to stdout.")
	flag.StringVar(&types, "t", "", "Specify the types. Multiple types are separated by comma(,). Extract all types if not specified.")
	flag.BoolVar(&private, "p", false, "Include private methods.")
}

func main() {
	var fileInfoList []*parser.SourceFileInfo
	flag.Parse()

	if types != "" {
		tmpTypes := strings.Split(types, ",")
		interestTypes = make(map[string]struct{})
		for _, t := range tmpTypes {
			interestTypes[t] = struct{}{}
		}
	}

	if inputFile == "" || inputFile == "-" { // read from stdin
		raw, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input := antlr.NewInputStream(string(raw))
		fileInfo, err := analyze(input)
		if err != nil {
			panic(err)
		}
		fileInfoList = []*parser.SourceFileInfo{fileInfo}

	} else {
		info, err := os.Stat(inputFile)
		if err != nil {
			panic(err)
		}
		if info.IsDir() { // is package
			files, err := ioutil.ReadDir(inputFile)
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
				filename := filepath.Join(inputFile, f.Name())
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
			input, err := antlr.NewFileStream(inputFile)
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

	gen := parser.InterfaceGenerator{Files: fileInfoList, Types: interestTypes}
	code, err := gen.GenerateCode()
	if err != nil {
		panic(err)
	}
	if outputFile == "" { // write to stdout
		fmt.Println(code)
	} else {
		f, err := os.Create(outputFile)
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
	listener := parser.NewMethodListener(private)
	walker := antlr.ParseTreeWalkerDefault
	walker.Walk(listener, tree)
	return listener.GetResult(), nil
}
