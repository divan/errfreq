package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"golang.org/x/tools/go/packages"
)

func main() {
	var zeros = flag.Bool("zeros", false, "Include functions with zero 'if err' checks")
	var raw = flag.Bool("raw", false, "Don't create histogram, just dump raw data to stdout")
	var output = flag.String("o", "errs_histogram.png", "Filename of histogram PNG")
	var autoopen = flag.Bool("autoopen", true, "Autoopen histogram")
	flag.Parse()

	// parse package specified in args
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
	}
	pkgs, err := packages.Load(cfg, os.Args[1:]...)
	if err != nil {
		log.Fatalf("Load package: %v", err)
	}

	a := NewAnalyzer(*zeros)
	for _, pkg := range pkgs {
		a.ParsePackage(pkg)
	}

	if !*raw {
		err := a.PlotHistogram(*output)
		if err != nil {
			log.Fatal(err)
		}
		if *autoopen {
			OpenPlot(*output)
		}
		return
	}

	result := make([]int, 0, len(a.Counts))
	for _, v := range a.Counts {
		if v == 0 && !*zeros {
			continue
		}

		result = append(result, v)
	}

	sort.Ints(result)
	for i := range result {
		fmt.Println(result[i])
	}
}

// for dogfooding
type X struct {
}

func (X) Foo() {
	bar := func() {}
	_ = bar
}
func (x X) Foo1() {
}
func (*X) Bar() {
}
func (x *X) Bar2() {
}

func err1() error {
	return nil
}
func err2() (int, error) {
	return 42, nil
}

func errTest() {
	if err := err1(); err != nil {
		return
	}
	if _, err := err2(); err != nil {
		return
	}

	err := err1()
	if err != nil {
		return
	}
}
