package parser

import (
	"fmt"
	"go/format"
	"sort"
	"strings"
)

type InterfaceGenerator struct {
	Files []*SourceFileInfo
	Types map[string]struct{}
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

	tps := make([]string, 0, len(structMap))
	for tp := range structMap {
		tps = append(tps, tp)
	}
	sort.Strings(tps)

	for _, typeName := range tps {
		repr := structMap[typeName]
		if gen.Types != nil {
			// check if typeName is in gen.Types
			if _, ok := gen.Types[typeName]; !ok {
				continue // skip
			}
		}
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

	sort.Slice(methods, func(i, j int) bool { return methods[i].Identifier < methods[j].Identifier })
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
