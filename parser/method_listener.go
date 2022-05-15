package parser

import (
	"strings"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

const (
	COMMENT = "// Code generated from gointerface\n"
)

type SourceFileInfo struct {
	PkgName string
	Imports []*ImportStmt
	Methods []*MethodDecl
}

type ImportStmt struct {
	Alias string
	Path  string
}

type ReceiverDecl struct {
	StructType string
	IsPointer  bool
}

type MethodDecl struct {
	Recv       *ReceiverDecl
	Identifier string
	Signature  string
	Comment    string
}

type MethodListener struct {
	*BaseGoParserListener
	IncludePrivate bool
	fileInfo       *SourceFileInfo
	inMethod       bool
	inReceiver     bool

	currentMethod *MethodDecl
}

func NewMethodListener(includePrivate bool) *MethodListener {
	return &MethodListener{BaseGoParserListener: &BaseGoParserListener{}, IncludePrivate: includePrivate}
}

func (s *MethodListener) GetResult() *SourceFileInfo {
	return s.fileInfo
}

// 注意这里的参数很容易写错，一旦写错就无法实现接口
func (s *MethodListener) EnterPackageClause(c *PackageClauseContext) {
	// fmt.Println(c.GetStart().GetLine(), c.packageName) // get line no
	s.fileInfo = &SourceFileInfo{PkgName: c.packageName.GetText()}

}

// EnterImportSpec is called when entering the importSpec production.
func (s *MethodListener) EnterImportSpec(c *ImportSpecContext) {
	var aliasStr string
	alias := c.alias
	if alias == nil {
		aliasStr = ""
	} else {
		aliasStr = alias.GetText()
	}
	// import alias可以位空 _ . identifier
	if aliasStr == "_" {
		// ignore
		return
	}
	importPath := c.ImportPath().GetText()
	imp := &ImportStmt{Alias: aliasStr, Path: importPath}
	s.fileInfo.Imports = append(s.fileInfo.Imports, imp)
}

func (s *MethodListener) EnterParameterDecl(ctx *ParameterDeclContext) {
	if !s.inReceiver {
		return
	}

	tp := ctx.Type_().GetText()
	recv := s.currentMethod.Recv
	if len(tp) > 0 && tp[0] == '*' {
		recv.StructType = tp[1:]
		recv.IsPointer = true
	} else {
		recv.StructType = tp
	}
}

func (s *MethodListener) EnterMethodDecl(ctx *MethodDeclContext) {
	ident := ctx.IDENTIFIER().GetText()
	if !s.IncludePrivate && unicode.IsLower(rune(ident[0])) {
		// skip private method
		return
	}

	startToken := ctx.GetStart()
	i := startToken.GetTokenIndex()
	stream := ctx.parser.GetInputStream().(*antlr.CommonTokenStream)
	tokens := stream.GetHiddenTokensToLeft(i, antlr.TokenHiddenChannel)
	comments := make([]string, 0, len(tokens))
	for _, t := range tokens {
		comments = append(comments, t.GetText())
	}
	comment := strings.Join(comments, "")
	s.inMethod = true
	// set current method
	s.currentMethod = &MethodDecl{
		Recv:       &ReceiverDecl{},
		Identifier: ident,
		Signature:  formatSignature(ctx.Signature()),
		Comment:    comment}
}

func formatSignature(sign ISignatureContext) string {
	stream := sign.GetParser().GetInputStream().(*antlr.CommonTokenStream)
	return stream.GetTextFromTokens(sign.GetStart(), sign.GetStop())
}

func (s *MethodListener) ExitMethodDecl(ctx *MethodDeclContext) {
	if !s.inMethod {
		return
	}
	s.fileInfo.Methods = append(s.fileInfo.Methods, s.currentMethod)
	s.inMethod = false
	s.currentMethod = nil
}

func (s *MethodListener) EnterReceiver(ctx *ReceiverContext) {
	if !s.inMethod {
		return
	}
	s.inReceiver = true
}

func (s *MethodListener) ExitReceiver(ctx *ReceiverContext) {
	s.inReceiver = false
}

type ErrorListener struct {
	// *antlr.DefaultErrorListener
}

// func (l *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
// 	panic(msg)
// }

func (l *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	panic(msg)
}
func (l *ErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
}
func (l *ErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
}
func (l *ErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
}

func NewErrorListener() *ErrorListener {
	return &ErrorListener{}
}
