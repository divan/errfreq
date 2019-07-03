package main

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// ID returns unique ID for the function (just for grouping counts)
func fnID(pkgName string, f *ast.FuncDecl) string {
	if f.Recv != nil {
		recv := ""
		recvT := f.Recv.List[0].Type
		switch x := recvT.(type) {
		case *ast.Ident:
			recv = x.Name
		case *ast.StarExpr:
			recv = x.X.(*ast.Ident).String()
		}
		return fmt.Sprintf("%s.%s.%s", pkgName, recv, f.Name.Name)
	}
	return fmt.Sprintf("%s.%s", pkgName, f.Name.Name)
}

// isIfErr detects if condition is an error check of our interest or not.
// TODO: this is naive approach, and probably can be fooled, but
// it should be good enought for most of straightforward cases, which
// are, franckly, is our target anyway.
func isIfErr(pkg *packages.Package, is *ast.IfStmt) bool {
	binExpr, ok := is.Cond.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	if binExpr.Op.String() != "!=" {
		return false
	}
	switch x := binExpr.X.(type) {
	case *ast.Ident:
		o := pkg.TypesInfo.ObjectOf(x)
		named, ok := o.Type().(*types.Named)
		if !ok {
			return false
		}
		return named.String() == "error"
	}
	return false
}
