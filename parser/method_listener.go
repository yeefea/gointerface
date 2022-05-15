package parser

import (
	"fmt"
	"go/format"
	"sort"
	"strings"

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
	fileInfo   *SourceFileInfo
	inMethod   bool
	inReceiver bool

	currentMethod *MethodDecl
}

func NewMethodListener() *MethodListener {
	return &MethodListener{BaseGoParserListener: &BaseGoParserListener{}}
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

// EnterParamResult is called when entering the paramResult production.
func (s *MethodListener) EnterParamResult(c *ParamResultContext) {

}

// EnterParamSimple is called when entering the paramSimple production.
func (s *MethodListener) EnterParamSimple(c *ParamSimpleContext) {

}

func (s *MethodListener) EnterMethodDecl(ctx *MethodDeclContext) {
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
	s.currentMethod = &MethodDecl{
		Recv:       &ReceiverDecl{},
		Identifier: ctx.IDENTIFIER().GetText(),
		Signature:  formatSignature(ctx.Signature()),
		Comment:    comment}
}

func formatSignature(sign ISignatureContext) string {
	// fmt.Println(start, stop)
	stream := sign.GetParser().GetInputStream().(*antlr.CommonTokenStream)

	return stream.GetTextFromTokens(sign.GetStart(), sign.GetStop())
}

func (s *MethodListener) ExitMethodDecl(ctx *MethodDeclContext) {
	s.fileInfo.Methods = append(s.fileInfo.Methods, s.currentMethod)
	s.inMethod = false
}

func (s *MethodListener) EnterReceiver(ctx *ReceiverContext) {
	s.inReceiver = true
}

func (s *MethodListener) ExitReceiver(ctx *ReceiverContext) {
	s.inReceiver = false
}

type InterfaceGenerator struct {
	Files []*SourceFileInfo
}

type interfaceRepr struct {
	ValueRecvMeth   []*MethodDecl
	PointerRecvMeth []*MethodDecl
}

func (gen *InterfaceGenerator) GenerateCode() (string, error) {
	if len(gen.Files) == 0 {
		return "", nil
	}
	sb := strings.Builder{}
	// emit comment
	sb.WriteString(COMMENT)

	// check the package name and emit package statement
	pkgName := gen.Files[0].PkgName
	for _, f := range gen.Files[1:] {
		if f.PkgName != pkgName {
			return "", fmt.Errorf("package name not same, %s != %s", f.PkgName, pkgName)
		}
	}
	sb.WriteString("package ")
	sb.WriteString(pkgName)
	sb.WriteString("\n")

	// emit imports
	imports := make([]*ImportStmt, 0)
	for _, f := range gen.Files {
		imports = append(imports, f.Imports...)
	}
	emitImports(&sb, imports)

	// emit interfaces
	structMap := map[string]*interfaceRepr{}
	for _, f := range gen.Files {
		for _, m := range f.Methods {
			recv := m.Recv
			tp := recv.StructType
			repr, ok := structMap[tp]
			if !ok {
				repr = &interfaceRepr{}
				structMap[tp] = repr
			}
			if recv.IsPointer {
				repr.PointerRecvMeth = append(repr.PointerRecvMeth, m)
			} else {
				repr.ValueRecvMeth = append(repr.ValueRecvMeth, m)
			}
		}
	}

	for typeName, repr := range structMap {
		if len(repr.PointerRecvMeth) != 0 && len(repr.ValueRecvMeth) != 0 {
			emitInterface(&sb, fmt.Sprintf("I%s", typeName), repr.PointerRecvMeth)
			emitInterface(&sb, fmt.Sprintf("I%sValue", typeName), repr.ValueRecvMeth)
		} else if len(repr.PointerRecvMeth) != 0 {
			emitInterface(&sb, fmt.Sprintf("I%s", typeName), repr.PointerRecvMeth)
		} else if len(repr.ValueRecvMeth) != 0 {
			emitInterface(&sb, fmt.Sprintf("I%s", typeName), repr.ValueRecvMeth)
		}
	}

	rawCode := sb.String()
	// fmt.Println(rawCode)
	// format interface code
	code, err := format.Source([]byte(rawCode))
	return string(code), err
}

func emitInterface(sb *strings.Builder, name string, methods []*MethodDecl) {
	sb.WriteString(fmt.Sprintf("type %s interface {\n\n", name))
	for _, m := range methods {
		emitMethod(sb, m)
	}
	sb.WriteString("}\n\n")
}

func emitImports(sb *strings.Builder, imp []*ImportStmt) {
	if len(imp) == 0 {
		return
	}

	// make unique
	importsMap := make(map[ImportStmt]struct{})
	for _, i := range imp {
		importsMap[*i] = struct{}{}
	}

	// sort
	sortedImports := make([]ImportStmt, 0, len(importsMap))
	for i := range importsMap {
		sortedImports = append(sortedImports, i)
	}

	sort.Slice(sortedImports, func(i, j int) bool {
		if sortedImports[i].Path == sortedImports[j].Path {
			return sortedImports[i].Alias < sortedImports[j].Alias
		}
		return sortedImports[i].Path < sortedImports[j].Path
	})

	sb.WriteString("import (\n")
	for _, i := range imp {
		sb.WriteString(fmt.Sprintf("%s %s\n", i.Alias, i.Path))
	}
	sb.WriteString(")\n")
}

func emitMethod(sb *strings.Builder, m *MethodDecl) {
	sb.WriteString("\n")
	sb.WriteString(m.Comment)
	sb.WriteString(m.Identifier)
	sb.WriteString(" ")
	sb.WriteString(m.Signature)
	sb.WriteString("\n")

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
