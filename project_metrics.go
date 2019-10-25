package main

import (
	"fmt"
)

func findProjectMetrics(metrics []FileMetric) ProjectMetric {
	total_sloc := 0
	too_long_files := 0
	too_long_methods := 0
	nesting_depth_sloc := 0

	for _, i := range metrics {
		total_sloc += i.FileLength

		if i.FileLength > LONG_FILE_THRESHOLD {
			too_long_files += i.FileLength
		}

		if i.TooLongMethod > LONG_METHOD_THRESHOLD {
			too_long_methods += i.TooLongMethod
		}

		nesting_depth_sloc += i.MaxNestingDepth
	}

	var m ProjectMetric

	if total_sloc == 0 {
		total_sloc++
	}

	m.TooLongFiles = float64(too_long_files) / float64(total_sloc)
	m.TooLongMethods = float64(too_long_methods) / float64(total_sloc)
	m.NestingDepth = float64(nesting_depth_sloc) / float64(total_sloc)
	m.CommentIncompleteness = 99.99 //////

	return m
}

type ProjectMetric struct {
	VersionName           string
	TooLongMethods        float64
	TooLongFiles          float64
	NestingDepth          float64
	CommentIncompleteness float64
}

func (p ProjectMetric) view() {
	fmt.Println("--- Project Maintainibility Metrics ---")
	fmt.Printf("Too Long Files: %.2f\n", p.TooLongFiles*100)
	fmt.Printf("Too Long Methods: %.2f\n", p.TooLongMethods*100)
	fmt.Printf("Nesting Depth: %.2f\n", p.NestingDepth*100)
	fmt.Printf("Comment Incompleteness: %.2f\n", p.CommentIncompleteness*100)
}

func viewEvolutionMetrics(metrics []ProjectMetric, csv bool) {
	if csv {
		fmt.Printf("version,long_file,long_method,complexity,comment")
		for _, i := range metrics {
			fmt.Printf("%s,%.2f,%.2f,%.2f,%.2f\n", i.VersionName, i.TooLongFiles, i.TooLongMethods, i.NestingDepth, i.CommentIncompleteness)
		}
	} else {
		fmt.Printf("%10s: %7s | %9s | %10s | %7s \n", "Version", "Long File", "Long Method", "Complexity", "Comment")
		for _, i := range metrics {
			fmt.Printf("%10s: %7f | %9f | %10f | %7f\n", i.VersionName, i.TooLongFiles, i.TooLongMethods, i.NestingDepth, i.CommentIncompleteness)
		}
	}
}
