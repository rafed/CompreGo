package main

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FileMetric contains metrics of a file
type FileMetric struct {
	FileName        string
	TooLongMethod   int
	MaxNestingDepth int
	FileLength      int
	TotalComments   int
	BadComments     int
}

func findFileMetrics(filename string) []FileMetric {
	var metrics []FileMetric

	filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			fset := token.NewFileSet()
			f, _ := parser.ParseFile(fset, path, nil, parser.ParseComments)
			file, _ := ioutil.ReadFile(path)
			contents := string(file)
			contents += " "

			if len(contents) == 0 {
				return err
			}

			var metric FileMetric
			metric.FileName = path
			metric.FileLength = findFileLength(contents)
			metric.TooLongMethod = findTooLongMethod(contents, f, fset, LONG_METHOD_THRESHOLD)
			metric.MaxNestingDepth = findMaxNestingDepth(contents, f, fset)
			metric.BadComments, metric.TotalComments = findComments(contents, f, fset)

			metrics = append(metrics, metric)
		}
		return err
	})

	return metrics
}
