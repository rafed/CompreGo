package main

import (
	"bufio"
	"go/ast"
	"go/token"
	"strings"
	"go/parser"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"fmt"
)

func forEachMethod(dir string){
	file, err := os.Open("./ignore_dirs.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var ignoreList []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ignoreList = append(ignoreList, scanner.Text())
	}

	print(dir)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// print(dir)
		if info.IsDir() && stringInSlice(info.Name(), ignoreList) {
			return filepath.SkipDir
		}

		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") { /*  */
			fset := token.NewFileSet()
			f, _ := parser.ParseFile(fset, path, nil, parser.ParseComments)
			file, _ := ioutil.ReadFile(path)
			contents := string(file)
			contents += " "

			if len(contents) == 0 {
				return err
			}

			findSLOCandErrors(contents, f, fset)
		}
		return err
	})
}

func findSLOCandErrors(contents string, f *ast.File, fset *token.FileSet) {
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			functionBody := contents[fn.Pos()-1 : fn.End()]

			SLOC := findNewLine(functionBody)
			NilCount := strings.Count(functionBody, " nil ")

			// fmt.Printf("%s: %d, %d\n", fn.Name, SLOC, NilCount)
			fmt.Printf("%d,%d", SLOC, NilCount)
		}
	}
}
