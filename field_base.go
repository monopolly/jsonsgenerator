package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"path"
	"path/filepath"
	"strings"

	"github.com/monopolly/cast"
)

// source file: imports and dir for custom types detection
var source struct {
	Dir      string
	Imports  []string
	importer types.ImporterFrom
	packages map[string]*types.Package
	locals   map[string]localType
}

// golang type of field: go{type=""} or original
func (a *Field) GoType() string {
	if a.Go.Type != "" {
		return a.Go.Type
	}
	return a.Type
}

// type for converters and sql: role.Role > int
func (a *Field) BaseType() string {
	if a.Go.Base != "" {
		return a.Go.Base
	}
	return a.GoType()
}

// original type with custom type replaced by base: role.Role > int
func (a *Field) PlainType() string {
	if a.Go.Base != "" && a.Go.Type == "" {
		return a.Go.Base
	}
	return a.Type
}

// struct value for sql/clickhouse: a.Role > int(a.Role)
func (a *Field) GoValue() string {
	if a.Go.Base != "" {
		return fmt.Sprintf("%s(a.%s)", a.Go.Base, a.Go.Name)
	}
	return "a." + a.Go.Name
}

// custom type: role role.Role > Base = int, BaseImport = github.com/x/role
// local type: role Role (type Role int, type Role = role.Role) > Base = int
func (a *Field) parseBaseType() {

	t := a.GoType()

	// uuid.UUID, []uuid.UUID: builtin support, only import
	if strings.Contains(t, uuidType) {
		a.Go.BaseImport = uuidImport
		return
	}

	// time.Time, time.Duration, date.Int...
	if a.Go.Base == "" && (goConvertSupported(t) || cast.Supported(t) == nil) {
		return
	}

	expr, err := parser.ParseExpr(t)
	if err != nil {
		return
	}

	switch e := expr.(type) {
	case *ast.SelectorExpr:
		// role.Role: generated code needs import
		if id, ok := e.X.(*ast.Ident); ok {
			if pkg := sourcePackage(id.Name, source.Imports); pkg != nil {
				a.Go.BaseImport = pkg.Path()
			}
		}
	case *ast.Ident:
		// Role: same package
		if typeBuiltin(e.Name) {
			return
		}
	default:
		return
	}

	// go{base="int"}
	if a.Go.Base != "" {
		return
	}

	base, err := resolveBaseType(expr, source.Imports, 0)
	if err != nil {
		panic(fmt.Sprintf("%s %s: %s. Set base type: go{base=\"int\"}", a.Name, t, err))
	}
	if !goConvertSupported(base) && cast.Supported(base) != nil {
		panic(fmt.Sprintf("%s %s: cast not support underlying type %s", a.Name, t, base))
	}

	a.Go.Base = base
}

// underlying builtin type: Role > role.Role > int
func resolveBaseType(expr ast.Expr, imports []string, depth int) (string, error) {

	if depth > 10 {
		return "", fmt.Errorf("type loop")
	}

	switch e := expr.(type) {

	// int or local type Role
	case *ast.Ident:
		if typeBuiltin(e.Name) {
			return e.Name, nil
		}
		local, ok := sourceLocalTypes()[e.Name]
		if !ok {
			return "", fmt.Errorf("type %s not found in package", e.Name)
		}
		return resolveBaseType(local.Expr, local.Imports, depth+1)

	// role.Role
	case *ast.SelectorExpr:
		id, ok := e.X.(*ast.Ident)
		if !ok {
			return "", fmt.Errorf("type %s is not supported", types.ExprString(e))
		}
		pkg := sourcePackage(id.Name, imports)
		if pkg == nil {
			return "", fmt.Errorf("package %s not found in imports", id.Name)
		}
		obj, ok := pkg.Scope().Lookup(e.Sel.Name).(*types.TypeName)
		if !ok {
			return "", fmt.Errorf("type %s not found", types.ExprString(e))
		}
		base := types.TypeString(obj.Type().Underlying(), nil)
		if strings.ContainsAny(base, ".{") {
			return "", fmt.Errorf("underlying type %s is not supported", base)
		}
		return base, nil

	// []string, map[string]int
	default:
		base := types.ExprString(e)
		if strings.ContainsAny(base, ".{") {
			return "", fmt.Errorf("underlying type %s is not supported", base)
		}
		return base, nil
	}
}

// int, string, bool...
func typeBuiltin(name string) bool {
	_, ok := types.Universe.Lookup(name).(*types.TypeName)
	return ok
}

type localType struct {
	Expr    ast.Expr // int, role.Role
	Imports []string // file imports
}

// type declarations of source package, without generated files
func sourceLocalTypes() map[string]localType {

	if source.locals != nil {
		return source.locals
	}
	source.locals = map[string]localType{}

	files, _ := filepath.Glob(filepath.Join(source.Dir, "*.go"))
	for _, name := range files {
		switch {
		case strings.HasSuffix(name, "GO.go"), strings.HasSuffix(name, "_test.go"), strings.HasSuffix(name, ".pb.go"):
			continue
		}
		f, err := parser.ParseFile(token.NewFileSet(), name, nil, parser.SkipObjectResolution)
		if err != nil {
			continue
		}
		var imports []string
		for _, x := range f.Imports {
			imports = append(imports, strings.Trim(x.Path.Value, `"`))
		}
		for _, d := range f.Decls {
			g, ok := d.(*ast.GenDecl)
			if !ok || g.Tok != token.TYPE {
				continue
			}
			for _, sp := range g.Specs {
				ts := sp.(*ast.TypeSpec)
				source.locals[ts.Name.Name] = localType{Expr: ts.Type, Imports: imports}
			}
		}
	}
	return source.locals
}

// package from imports by name: role > github.com/x/role
func sourcePackage(name string, imports []string) *types.Package {

	if source.importer == nil {
		build.Default.Dir = source.Dir // module resolution from the source dir
		source.importer = importer.ForCompiler(token.NewFileSet(), "source", nil).(types.ImporterFrom)
		source.packages = map[string]*types.Package{}
	}

	// by path: github.com/x/role, github.com/x/role/v2
	for _, x := range imports {
		base := path.Base(x)
		if strings.HasPrefix(base, "v") && cast.Int(base[1:]) > 0 {
			base = path.Base(path.Dir(x))
		}
		if base != name {
			continue
		}
		pkg, ok := source.packages[x]
		if !ok {
			var err error
			pkg, err = source.importer.ImportFrom(x, source.Dir, 0)
			if err != nil {
				fmt.Println("custom type:", x, err)
			}
			source.packages[x] = pkg
		}
		if pkg != nil && pkg.Name() == name {
			return pkg
		}
	}
	return nil
}
