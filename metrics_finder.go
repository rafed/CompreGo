package main

import (
	"go/ast"
	"go/token"
	"strings"
)

func findNewLine(contents string) int {
	count := 0
	for _, b := range contents {
		if b == '\n' {
			count++
		}
	}
	return count
}

func findFileLength(contents string) int {
	return findNewLine(contents) + 1
}

func findTooLongMethod(contents string, f *ast.File, fset *token.FileSet, threshold int) int {
	longMethodLines := 0
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			functionBody := contents[fn.Pos()-1 : fn.End()]

			if numOfLines := findNewLine(functionBody); numOfLines > threshold {
				longMethodLines += numOfLines
			}
		}
	}
	return longMethodLines
}

func findMaxNestingDepth(contents string, f *ast.File, fset *token.FileSet) int {
	nestedMethodLines := 0
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			v := &blockNestingVisitor{
				contents: contents,
			}
			ast.Walk(v, fn)

			if v.maxNesting == 0 {
				v.maxNesting = 0
			}

			// println(fn.Name.Name, v.maxNesting)

			if v.maxNesting > 3 {
				numOfLines := findNewLine(contents[fn.Pos()-1 : fn.End()])
				nestedMethodLines += numOfLines
			}
		}
	}
	return nestedMethodLines
}

type blockNestingVisitor struct {
	blocks       []*ast.BlockStmt
	maxNesting   int
	totalNesting int
	contents     string
}

func (v *blockNestingVisitor) Visit(node ast.Node) ast.Visitor {
	if v.blocks == nil {
		v.blocks = make([]*ast.BlockStmt, 0)
	}
	if node != nil {
		if b, is := node.(*ast.BlockStmt); is {
			v.calcMaxNesting(b)
			v.calcTotalNesting(b)
		}
	}
	return v
}

func (v *blockNestingVisitor) calcTotalNesting(b *ast.BlockStmt) {
	body := v.contents[b.Pos()-1 : b.End()]
	body = strings.TrimSpace(strings.Trim(strings.TrimSpace(body), "{}"))
	c := findNewLine(body)
	v.totalNesting += c
}

func (v *blockNestingVisitor) calcMaxNesting(b *ast.BlockStmt) {
	depth := 0
	for _, previous := range v.blocks {
		if previous.Pos() < b.Pos() && b.End() < previous.End() {
			depth++
			if depth > v.maxNesting {
				v.maxNesting = depth
			}
		}
	}
	v.blocks = append(v.blocks, b)
}

func findCommentIncompleteness() float32 {
	return 0
}
