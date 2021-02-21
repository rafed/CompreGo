package main

import (
	"bufio"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FileMetric contains metrics of a file
type FileMetric struct {
	FileName   string
	FileLength int

	tooLongMethods []tooLongMethodStorer
	nilsInMethods  []nilsInMethodStorer

	nestingDepth []nestingDepthStorer

	totalComments     int
	badComments       []badCommentStorer
	duplicateComments duplicateCommentStorer
}

type nilsInMethodStorer struct {
	FunctionCount int
	NilCount      int
	SLOC          int
}

type tooLongMethodStorer struct {
	FunctionName        string
	TooLongMethodLength int
	FunctionBody        string
}

type nestingDepthStorer struct {
	FunctionName      string
	MaxNestingDepth   int
	NestingDepthLines int
}

type badCommentStorer struct {
	FunctionName string
	comment      string
}

type dup struct {
	d1 string
	d2 string
}

type duplicateCommentStorer struct {
	count      int
	duplicates []dup
}

func findFileMetrics(filename string) []FileMetric {
	var metrics []FileMetric

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

	filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
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

			// fmt.Printf("Filepath: %s\n", path)

			var metric FileMetric
			metric.FileName = path
			metric.FileLength = findFileLength(contents)
			metric.tooLongMethods = findTooLongMethod(contents, f, fset)
			metric.nestingDepth = findMaxNestingDepth(contents, f, fset, path)
			metric.totalComments, metric.badComments, metric.duplicateComments = findComments(contents, f, fset)

			metrics = append(metrics, metric)
		}
		return err
	})

	return metrics
}
