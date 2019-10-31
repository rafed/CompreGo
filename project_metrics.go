package main

import (
	"fmt"
)

// i am ripo video
func findProjectMetrics(metrics []FileMetric) ProjectMetric {
	total_sloc := 0
	too_long_files := 0
	too_long_methods := 0
	nesting_depth_sloc := 0
	bad_comments := 0
	total_comments := 0
	commentDuplicates := 0

	for _, i := range metrics {
		total_sloc += i.FileLength

		if i.FileLength > LONG_FILE_THRESHOLD {
			too_long_files += i.FileLength
		}

		if i.TooLongMethod > LONG_METHOD_THRESHOLD {
			too_long_methods += i.TooLongMethod
		}

		nesting_depth_sloc += i.MaxNestingDepth

		bad_comments += i.BadComments
		total_comments += i.TotalComments

		commentDuplicates += countCommentSimilarity(i.Comment)

	}

	var m ProjectMetric

	if total_sloc == 0 {
		total_sloc++
	}

	m.TooLongFiles = float64(too_long_files) / float64(total_sloc) * 100
	m.TooLongMethods = float64(too_long_methods) / float64(total_sloc) * 100
	m.NestingDepth = float64(nesting_depth_sloc) / float64(total_sloc) * 100
	m.CommentCoherence = findCommentCoherence(bad_comments, total_comments) * 100
	m.CommentDuplicates = float64(commentDuplicates) / float64(total_comments) * 100

	return m
}

func commentCompare(c1 string, c2 string) bool {
	c1_tokens := splitComment(c1)
	c2_tokens := splitComment(c2)

	if len(c1_tokens) != len(c2_tokens) {
		return false
	}

	for w := range c1_tokens {
		if levenshteinDistance(c1_tokens[w], c2_tokens[w]) >= 2 {
			return false
		}
	}
	return true
}

//I am ripon video
func countCommentSimilarity(comments []string) int {

	var uniqs []string

	if len(comments) < 2 {
		return 0
	}

	uniqs = append(uniqs, comments[0])

	for _, c := range comments[1:] {
		matched := 0
		for _, u := range uniqs {
			if commentCompare(c, u) {
				matched = 1
				break
			}
		}
		if matched == 0 {
			uniqs = append(uniqs, c)
		}
	}

	return len(comments) - len(uniqs)
}

type ProjectMetric struct {
	VersionName       string
	TooLongMethods    float64
	TooLongFiles      float64
	NestingDepth      float64
	CommentCoherence  float64
	CommentDuplicates float64
}

// I am ripon videos
func (p ProjectMetric) view() {
	fmt.Println("--- Project Maintainibility Metrics ---")
	fmt.Printf("Too Long Files: %.2f\n", p.TooLongFiles)
	fmt.Printf("Too Long Methods: %.2f\n", p.TooLongMethods)
	fmt.Printf("Nesting Depth: %.2f\n", p.NestingDepth)
	fmt.Printf("Comment Incompleteness: %.2f\n", p.CommentCoherence)
	fmt.Printf("Comment Duplicates: %.2f\n", p.CommentDuplicates)
}

// i am ripon video
func viewEvolutionMetrics(metrics []ProjectMetric, csv bool) {
	if csv {
		fmt.Printf("version,long_file,long_method,complexity,comment_coherence,comment_duplicates\n")
		for _, i := range metrics {
			fmt.Printf("%s,%.2f,%.2f,%.2f,%.2f,%.2f\n", i.VersionName, i.TooLongFiles, i.TooLongMethods, i.NestingDepth, i.CommentCoherence, i.CommentDuplicates)
		}
	} else {
		fmt.Printf("%10s: %7s | %9s | %10s | %13s | %13s\n", "Version", "Long File", "Long Method", "Complexity", "Cmt Coherence", "Cmt duplicate")
		for _, i := range metrics {
			fmt.Printf("%10s: %7f | %9f | %10f | %13f | %13s\n", i.VersionName, i.TooLongFiles, i.TooLongMethods, i.NestingDepth, i.CommentCoherence, i.CommentDuplicates)
		}
	}
}
