package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"reflect"
	"strings"
)

// begin transaction for TX methods
func (a *SQL) Begin() []byte {
	s := `
	//begin transaction for TX methods:
	//	tx, ctx, err := a.Begin()
	//	defer tx.Rollback(ctx)
	//	a.InsertTX(tx, ctx, item)
	//	err = tx.Commit(ctx)
	func (a *#CLASS) Begin() (tx pgx.Tx, ctx context.Context, err error) {
		ctx = context.Background()
		tx, err = a.pool.Begin(ctx)
		return
	}`
	return []byte(strings.ReplaceAll(s, "#CLASS", settings.SQL.Class))
}

// TX versions of generated sql methods:
// Get(id any) -> GetTX(tx pgx.Tx, ctx context.Context, id any)
// a.Conn(func(conn) {...}) body is inlined: conn -> tx, context.Background() -> ctx
func (a *SQL) TX(code []byte) []byte {

	src := append([]byte("package x\n"), code...)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.SkipObjectResolution)
	if err != nil {
		fmt.Println("sql tx: parse generated code:", err)
		return nil
	}

	// class methods
	var order []string
	methods := map[string]*ast.FuncDecl{}
	for _, d := range f.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || fn.Body == nil || txRecv(fn) != settings.SQL.Class {
			continue
		}
		order = append(order, fn.Name.Name)
		methods[fn.Name.Name] = fn
	}

	// methods with a.Conn or calls of other TX methods
	need := map[string]bool{}
	for changed := true; changed; {
		changed = false
		for name, fn := range methods {
			if need[name] || name == "Conn" || strings.HasSuffix(name, "TX") {
				continue
			}
			if txUses(fn, need) {
				need[name] = true
				changed = true
			}
		}
	}

	var res bytes.Buffer
	for _, name := range order {
		if !need[name] {
			continue
		}
		fn := methods[name]
		b, err := txMethod(src[fset.Position(fn.Pos()).Offset:fset.Position(fn.End()).Offset], need)
		if err != nil {
			fmt.Println("sql tx:", name, err)
			continue
		}
		fmt.Fprintf(&res, "\n\n//%sTX transaction version of %s\n", name, name)
		res.Write(b)
	}
	return res.Bytes()
}

// receiver type name: (a *NewsSQL) -> NewsSQL
func txRecv(fn *ast.FuncDecl) string {
	if len(fn.Recv.List) == 0 {
		return ""
	}
	t := fn.Recv.List[0].Type
	if s, ok := t.(*ast.StarExpr); ok {
		t = s.X
	}
	if id, ok := t.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// a.Method(...) -> Method
func txCall(e ast.Expr) (*ast.CallExpr, string) {
	call, ok := e.(*ast.CallExpr)
	if !ok {
		return nil, ""
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, ""
	}
	if id, ok := sel.X.(*ast.Ident); !ok || id.Name != "a" {
		return nil, ""
	}
	return call, sel.Sel.Name
}

func txUses(fn *ast.FuncDecl, need map[string]bool) (ok bool) {
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if e, is := n.(ast.Expr); is {
			if _, name := txCall(e); name == "Conn" || need[name] {
				ok = true
			}
		}
		return !ok
	})
	return
}

// context.Background()
func txIsBackground(e ast.Expr) bool {
	call, ok := e.(*ast.CallExpr)
	if !ok || len(call.Args) > 0 {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	return ok && id.Name == "context" && sel.Sel.Name == "Background"
}

// a.Conn(func(conn *pgxpool.Conn) (err error) {...}) -> {...}
func txConnBody(e ast.Expr) []ast.Stmt {
	call, name := txCall(e)
	if name != "Conn" || len(call.Args) != 1 {
		return nil
	}
	if lit, ok := call.Args[0].(*ast.FuncLit); ok {
		return lit.Body.List
	}
	return nil
}

func txMethod(code []byte, need map[string]bool) ([]byte, error) {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", append([]byte("package x\n"), code...), parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	fn := f.Decls[0].(*ast.FuncDecl)
	fn.Doc = nil
	fn.Name.Name += "TX"
	fn.Type.Params.List = append([]*ast.Field{
		{Names: []*ast.Ident{{Name: "tx"}}, Type: &ast.SelectorExpr{X: ast.NewIdent("pgx"), Sel: ast.NewIdent("Tx")}},
		{Names: []*ast.Ident{{Name: "ctx"}}, Type: &ast.SelectorExpr{X: ast.NewIdent("context"), Sel: ast.NewIdent("Context")}},
	}, fn.Type.Params.List...)

	// closure returned only err: return X -> err = X; return
	multi := fn.Type.Results != nil && fn.Type.Results.NumFields() > 1

	stmts := func(list []ast.Stmt) (res []ast.Stmt) {
		defer func() { res = txDedupReturn(res) }()
		for _, s := range list {
			switch x := s.(type) {
			case *ast.AssignStmt:
				// c := context.Background()
				if len(x.Rhs) == 1 && txIsBackground(x.Rhs[0]) {
					continue
				}
				// err = a.Conn(...)
				if len(x.Rhs) == 1 {
					if body := txConnBody(x.Rhs[0]); body != nil {
						res = append(res, body...)
						continue
					}
				}
			case *ast.ReturnStmt:
				// return a.Conn(...)
				if len(x.Results) == 1 {
					if body := txConnBody(x.Results[0]); body != nil {
						res = append(res, body...)
						continue
					}
				}
			}
			res = append(res, s)
		}
		if !multi {
			return
		}
		var out []ast.Stmt
		for _, s := range res {
			if r, ok := s.(*ast.ReturnStmt); ok && len(r.Results) == 1 {
				out = append(out,
					&ast.AssignStmt{Lhs: []ast.Expr{ast.NewIdent("err")}, Tok: token.ASSIGN, Rhs: r.Results},
					&ast.ReturnStmt{},
				)
				continue
			}
			out = append(out, s)
		}
		return out
	}

	exprs := func(e ast.Expr) ast.Expr {
		if txIsBackground(e) {
			return ast.NewIdent("ctx")
		}
		if id, ok := e.(*ast.Ident); ok {
			switch id.Name {
			case "conn":
				id.Name = "tx"
			case "c":
				id.Name = "ctx"
			}
		}
		// a.Get(...) -> a.GetTX(tx, ctx, ...)
		if call, name := txCall(e); need[name] {
			call.Fun.(*ast.SelectorExpr).Sel.Name = name + "TX"
			call.Args = append([]ast.Expr{ast.NewIdent("tx"), ast.NewIdent("ctx")}, call.Args...)
		}
		return e
	}

	txWalk(reflect.ValueOf(fn.Body), exprs, stmts)

	var b bytes.Buffer
	err = printer.Fprint(&b, fset, fn)
	return b.Bytes(), err
}

// inlined a.Conn body ends with return + outer return: drop the unreachable one
func txDedupReturn(list []ast.Stmt) (res []ast.Stmt) {
	for _, s := range list {
		if r, ok := s.(*ast.ReturnStmt); ok && len(r.Results) == 0 && len(res) > 0 {
			if _, prev := res[len(res)-1].(*ast.ReturnStmt); prev {
				continue
			}
		}
		res = append(res, s)
	}
	return
}

// ast rewrite by reflection: replaces ast.Expr values and []ast.Stmt lists
func txWalk(v reflect.Value, exprs func(ast.Expr) ast.Expr, stmts func([]ast.Stmt) []ast.Stmt) {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if !v.IsNil() {
			txWalk(v.Elem(), exprs, stmts)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			txField(v.Field(i), exprs, stmts)
		}
	case reflect.Slice:
		if list, ok := v.Interface().([]ast.Stmt); ok && v.CanSet() {
			v.Set(reflect.ValueOf(stmts(list)))
		}
		for i := 0; i < v.Len(); i++ {
			txField(v.Index(i), exprs, stmts)
		}
	}
}

func txField(f reflect.Value, exprs func(ast.Expr) ast.Expr, stmts func([]ast.Stmt) []ast.Stmt) {
	if !f.CanSet() {
		return
	}
	if f.Kind() == reflect.Interface && !f.IsNil() {
		if e, ok := f.Interface().(ast.Expr); ok {
			if r := exprs(e); r != e {
				f.Set(reflect.ValueOf(r))
				return
			}
		}
	}
	txWalk(f, exprs, stmts)
}
