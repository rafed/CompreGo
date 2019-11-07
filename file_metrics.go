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
	FileName        string
	TooLongMethod   int
	MaxNestingDepth int
	FileLength      int
	TotalComments   int
	BadComments     int
	Comment         []string
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

			var metric FileMetric
			metric.FileName = path
			metric.FileLength = findFileLength(contents)
			metric.TooLongMethod = findTooLongMethod(contents, f, fset, LONG_METHOD_THRESHOLD)
			// println(info.Name())
			metric.MaxNestingDepth = findMaxNestingDepth(contents, f, fset, path)
			metric.BadComments, metric.TotalComments, metric.Comment = findComments(contents, f, fset)

			metrics = append(metrics, metric)
		}
		return err
	})

	return metrics
}
