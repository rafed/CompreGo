package main

import (
	"fmt"
)

// i am ripo video
func findProjectMetrics(metrics []FileMetric) ProjectMetric {
	totalSloc := 0
	tooLongFiles := 0
	tooLongMethods := 0
	nestingDepthSloc := 0
	totalComments := 0

	badComments := 0
	commentDuplicates := 0

	for _, i := range metrics {
		totalSloc += i.FileLength

		if i.FileLength > LONG_FILE_THRESHOLD {
			// println(i.FileName)
			tooLongFiles += i.FileLength
		}

		for _, lm := range i.tooLongMethods {
			if lm.TooLongMethodLength > LONG_METHOD_THRESHOLD {
				tooLongMethods += lm.TooLongMethodLength
			}
		}

		for _, nd := range i.nestingDepth {
			nestingDepthSloc += nd.MaxNestingDepth
		}

		totalComments += i.totalComments
		badComments += len(i.badComments)
		commentDuplicates += i.duplicateComments.count
	}

	var m ProjectMetric

	if totalSloc == 0 {
		totalSloc++
	}

	m.TooLongFiles = float64(tooLongFiles) / float64(totalSloc) * 100
	m.TooLongMethods = float64(tooLongMethods) / float64(totalSloc) * 100
	m.NestingDepth = float64(nestingDepthSloc) / float64(totalSloc) * 100
	m.CommentCoherence = findCommentCoherence(badComments, totalComments) * 100
	m.CommentDuplicates = findDuplicateComments(commentDuplicates, totalComments) * 100

	return m
}

// ProjectMetric matha amar
type ProjectMetric struct {
	VersionName       string
	TooLongMethods    float64
	TooLongFiles      float64
	NestingDepth      float64
	CommentCoherence  float64
	CommentDuplicates float64
}

// I am ripon videos
func (p ProjectMetric) view(csv bool) {
	if !csv {
		fmt.Println("--- Project Maintainability Metrics ---")
		fmt.Printf("%-25s: %05.2f\n", "Too Long Files", p.TooLongFiles)
		fmt.Printf("%-25s: %05.2f\n", "Too Long Methods", p.TooLongMethods)
		fmt.Printf("%-25s: %05.2f\n", "Nesting Depth", p.NestingDepth)
		fmt.Printf("%-25s: %05.2f\n", "Lack of cohesive comments", p.CommentCoherence)
		fmt.Printf("%-25s: %05.2f\n", "Duplicate comments", p.CommentDuplicates)
	} else {
		// fmt.Println("lf,lm,nd,lcc,cd")
		fmt.Printf("%.2f,%.2f,%.2f,%.2f,%.2f\n", p.TooLongFiles, p.TooLongMethods, p.NestingDepth, p.CommentCoherence, p.CommentDuplicates)
	}

}

// i am ripon video
func viewEvolutionMetrics(metrics []ProjectMetric, csv bool) {
	fmt.Println()
	if csv {
		fmt.Printf("version,long_file,long_method,complexity,comment_coherence,comment_duplicates\n")
		for _, i := range metrics {
			fmt.Printf("%s,%.2f,%.2f,%.2f,%.2f,%.2f\n", i.VersionName, i.TooLongFiles, i.TooLongMethods, i.NestingDepth, i.CommentCoherence, i.CommentDuplicates)
		}
	} else {
		fmt.Printf("%-10s | %-5s | %-5s | %-5s | %-5s | %-5s\n", "Version", "TLF", "TLM", "ND", "LCC", "DC")
		for _, i := range metrics {
			fmt.Printf("%-10s | %05.2f | %05.2f | %05.2f | %05.2f | %05.2f\n", i.VersionName, i.TooLongFiles, i.TooLongMethods, i.NestingDepth, i.CommentCoherence, i.CommentDuplicates)
		}
	}
}

func viewMetricValues(metrics []FileMetric, TLF bool, TLM bool, ND bool, LCC bool, DC bool) {
	if TLF {
		fmt.Println("\nToo Long Files:")
		for _, i := range metrics {
			if i.FileLength > LONG_FILE_THRESHOLD {
				fmt.Println("  ", i.FileName)
			}
		}
	}
	if TLM {
		fmt.Println("\nToo Long Methods:")
		for _, i := range metrics {
			for _, j := range i.tooLongMethods {
				fmt.Printf("  %s:%s()  len:%d\n", i.FileName, j.FunctionName, j.TooLongMethodLength)
			}
		}
	}
	if ND {
		fmt.Println("\nDeep Nesting Depth:")
		for _, i := range metrics {
			for _, j := range i.nestingDepth {
				fmt.Printf("  %s:%s() depth:%d\n", i.FileName, j.FunctionName, j.MaxNestingDepth)
			}
		}
	}
	if LCC {
		fmt.Println("\nLow Comment Cohesion:")
		for _, i := range metrics {
			for _, j := range i.badComments {
				fmt.Printf("  %s:%s()\n", i.FileName, j.FunctionName)
				fmt.Printf("  Comment: %s\n", j.comment)
			}
		}
	}
	if DC {
		fmt.Println("\nDuplicate Comments:")
		for _, i := range metrics {
			// if len()
			fmt.Printf("  %s\n", i.FileName)
			for c, j := range i.duplicateComments.duplicates {
				fmt.Printf("    [%d] %s", c, j.d1)
				fmt.Printf("    [%d] %s", c, j.d2)
			}
		}
	}
}
