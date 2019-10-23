package main

import "fmt"

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

	m.TooLongFiles = float32(too_long_files) / float32(total_sloc)
	m.TooLongMethods = float32(too_long_methods) / float32(total_sloc)
	m.NestingDepth = float32(nesting_depth_sloc) / float32(total_sloc)
	m.CommentIncompleteness = 9999

	return m
}

type ProjectMetric struct {
	VersionName           string
	TooLongMethods        float32
	TooLongFiles          float32
	NestingDepth          float32
	CommentIncompleteness float32
}

func (p ProjectMetric) view() {
	fmt.Println("--- Project Maintainibility Metrics ---")
	fmt.Println("Too Long Files: ", p.TooLongFiles)
	fmt.Println("Too Long Methods: ", p.TooLongMethods)
	fmt.Println("Nesting Depth: ", p.NestingDepth)
	fmt.Println("Comment Incompleteness ", p.CommentIncompleteness)
}

func viewEvolutionMetrics(metrics []ProjectMetric) {
	fmt.Printf("%10s: %7s | %9s | %10s | %7s \n", "Version", "Long File", "Long Method", "Nesting Depth", "Comment")
	for _, i := range metrics {
		fmt.Printf("%10s: %7f | %9f | %10f | %7f\n", i.VersionName, i.TooLongFiles, i.TooLongMethods, i.NestingDepth, i.CommentIncompleteness)
	}
}
