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

func findTooLongMethod(contents string, f *ast.File, fset *token.FileSet) []tooLongMethodStorer {
	var list []tooLongMethodStorer
	for _, decl := range f.Decls {
		var v tooLongMethodStorer
		if fn, ok := decl.(*ast.FuncDecl); ok {
			functionBody := contents[fn.Pos()-1 : fn.End()]

			if numOfLines := findNewLine(functionBody); numOfLines > LONG_METHOD_THRESHOLD {
				v.FunctionName = fn.Name.Name
				v.TooLongMethodLength = numOfLines
				list = append(list, v)
			}
		}
	}
	return list
}

var maxBlockDepth int

func findMaxNestingDepth(contents string, f *ast.File, fset *token.FileSet, path string) []nestingDepthStorer {
	var list []nestingDepthStorer
	for _, decl := range f.Decls {
		var vnds nestingDepthStorer
		if fn, ok := decl.(*ast.FuncDecl); ok {

			// if fn.Name.Name != "ReadPrefixed" {
			// 	continue
			// }

			// a := astutil.AddImport(fset, f, path)

			// v := &blockNestingVisitor{
			// 	contents: contents,
			// }
			// ast.Walk(v, fn)
			// println("depth", v.maxNesting) /////////////////////

			maxBlockDepth = 0
			v := &visitor2{
				depth:        0,
				blockCounter: 0,
			}
			ast.Walk(v, fn)

			if maxBlockDepth >= NESTING_DEPTH_THRESHOLD {
				numOfLines := findNewLine(contents[fn.Pos()-1 : fn.End()])
				vnds.MaxNestingDepth = numOfLines
				vnds.FunctionName = fn.Name.Name
				list = append(list, vnds)
			}

		}
	}
	return list
}

type visitor2 struct {
	depth        int
	blockCounter int
}

func (v visitor2) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}

	if v.blockCounter > maxBlockDepth {
		maxBlockDepth = v.blockCounter
	}

	// fmt.Printf("%d %s%T\n", v.blockCounter, strings.Repeat("\t", v.depth), n)

	switch n.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.SwitchStmt: //*ast.FuncDecl, *ast.IfStmt:
		v.blockCounter++
	}

	v.depth++
	return v
}

// type blockNestingVisitor struct {
// 	blocks       []*ast.BlockStmt
// 	maxNesting   int
// 	totalNesting int
// 	contents     string
// }

// func (v *blockNestingVisitor) Visit(node ast.Node) ast.Visitor {
// 	if v.blocks == nil {
// 		v.blocks = make([]*ast.BlockStmt, 0)
// 		// fmt.Printf("%v", v.block)
// 	}
// 	if node != nil {
// 		// depthCounter = 0

// 		if b, is := node.(*ast.BlockStmt); is {
// 			v.calcMaxNesting(b)
// 			v.calcTotalNestingLines(b)
// 		}
// 	}
// 	return v
// }

// func (v *blockNestingVisitor) calcMaxNesting(b *ast.BlockStmt) {
// 	depth := 0
// 	for _, previous := range v.blocks {
// 		if previous.Pos() < b.Pos() && b.End() < previous.End() {
// 			depth++
// 			if depth > v.maxNesting {
// 				v.maxNesting = depth
// 			}
// 		}
// 	}
// 	v.blocks = append(v.blocks, b)
// }

// func (v *blockNestingVisitor) calcTotalNestingLines(b *ast.BlockStmt) {
// 	body := v.contents[b.Pos()-1 : b.End()]
// 	body = strings.TrimSpace(strings.Trim(strings.TrimSpace(body), "{}"))
// 	c := findNewLine(body)
// 	v.totalNesting += c
// }

func findWordMatch(cms []string, names []string) float64 {
	cms = uniqueList(cms)
	names = uniqueList(names)
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

func findComments(contents string, f *ast.File, fset *token.FileSet) (int, []badCommentStorer, duplicateCommentStorer) {

	var allComments []string
	var badComments []badCommentStorer

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

			allComments = append(allComments, comment)

			coherence := findWordMatch(cms, nam)

			if coherence == 0 || coherence > .5 {
				var bc badCommentStorer
				bc.FunctionName = funcName
				bc.comment = comment
				badComments = append(badComments, bc)
			}
			total++
		}
	}

	dcs := countCommentSimilarity(allComments)

	return total, badComments, dcs
}

func findCommentCoherence(badComments int, totalComments int) float64 {
	if totalComments == 0 {
		totalComments++
	}
	ans := float64(badComments) / float64(totalComments)
	return ans
}

func findDuplicateComments(commentDuplicates int, totalComments int) float64 {
	if totalComments == 0 {
		totalComments++
	}
	return float64(commentDuplicates) / float64(totalComments)
}

//I am ripon video
func countCommentSimilarity(comments []string) duplicateCommentStorer {

	var duplicates duplicateCommentStorer
	var uniqs []string

	if len(comments) < 2 {
		return duplicates
	}

	uniqs = append(uniqs, comments[0])

	for _, c := range comments[1:] {
		var d dup
		matched := 0
		for _, u := range uniqs {
			if commentCompare(c, u) {
				matched = 1
				d.d1 = c
				d.d2 = u
				duplicates.duplicates = append(duplicates.duplicates, d)
				break
			}
		}
		if matched == 0 {
			uniqs = append(uniqs, c)
		}
	}

	duplicates.count = len(comments) - len(uniqs)
	return duplicates
}

func commentCompare(c1 string, c2 string) bool {
	c1Tokens := splitComment(c1)
	c2Tokens := splitComment(c2)

	if len(c1Tokens) != len(c2Tokens) {
		return false
	}

	for w := range c1Tokens {
		if levenshteinDistance(c1Tokens[w], c2Tokens[w]) >= 2 {
			return false
		}
	}
	return true
}
