package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

// Analyzer holds state for analysis across packages.
type Analyzer struct {
	// counts stores number of 'if err !=' checks per function
	Counts map[string]int

	countZeros bool
}

// NewAnalyzer inits new analyzer.
func NewAnalyzer(zeros bool) *Analyzer {
	return &Analyzer{
		Counts:     map[string]int{},
		countZeros: zeros,
	}
}

func (a *Analyzer) ParsePackage(pkg *packages.Package) {
	fmt.Println("Analyzing package:", pkg)
	currentFunc := ""

	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			// Find func declarations
			switch x := n.(type) {
			case *ast.FuncDecl:
				currentFunc = fnID(pkg.PkgPath, x)
				if a.countZeros {
					a.Counts[currentFunc] = 0
				}
			case *ast.IfStmt:
				isErr := isIfErr(pkg, x)
				if isErr {
					a.Counts[currentFunc]++
				}
			}
			return true
		})
	}
}
