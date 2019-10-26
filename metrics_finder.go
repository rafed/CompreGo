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

			if v.maxNesting > NESTING_DEPTH_THRESHOLD {
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

// cALC Max Nesting boo
func (v *blockNestingVisitor) calcMaxNesting(b *ast.BlockStmt) {
	/* damn nigga */
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

//hbdj shjdbdj skhdbdsjhb
func findWordMatch(cms []string, names []string) float64 {
	var matches []string
	for _, c := range cms {
		for _, w := range names {
			c = strings.ToLower(c)
			w = strings.ToLower(w)
			if levenshteinDistance(c, w) < 2 {
				matches = append(matches, w)

			}

		}
	}
	ans := float64(len(matches)) / float64(len(cms))
	return ans
}

// tui akta , coments sf12312 12412 asd
func findComments(contents string, f *ast.File, fset *token.FileSet) (int, int) {
	badComment := 0
	total := 0
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			funcName := fn.Name.Name
			comment := fn.Doc.Text()

			nam := splitCamelCase(funcName)
			cms := splitComment(comment)

			if len(cms) == 0 {
				continue
			}

			coherence := findWordMatch(cms, nam)

			if coherence == 0 || coherence > .5 {
				badComment++
			}
			total++
		}
	}

	return badComment, total
}

func findCommentCoherence(badComments int, totalComments int) float64 {
	if totalComments == 0 {
		totalComments++
	}
	ans := float64(badComments) / float64(totalComments)
	return ans
}
