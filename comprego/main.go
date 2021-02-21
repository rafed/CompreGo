package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var LONG_FILE_THRESHOLD int
var LONG_METHOD_THRESHOLD int
var NESTING_DEPTH_THRESHOLD int

func main() {

	lf := flag.Int("lf", 750, "Long File Threshold")
	lm := flag.Int("lm", 75, "Long Method Threshold")
	nd := flag.Int("nd", 5, "Nesting Depth Threshold")

	d := flag.String("d", "", "maintainibility metrics of a project")
	e := flag.String("e", "", "evolution of maintainibility metrics for each version")
	// mec := flag.String("mec", "", "Method and Error correlation")

	TLF := flag.Bool("TLF", false, "Show too long files")
	TLM := flag.Bool("TLM", false, "Show too long methods")
	ND := flag.Bool("ND", false, "Show methods with deep nesting methods")
	LCC := flag.Bool("LCC", false, "Show comments with low cohesion")
	DC := flag.Bool("DC", false, "Show duplicate comments")
	ALL := flag.Bool("ALL", false, "Show verbose results")

	if *ALL {
		*TLF = true
		*TLM = true
		*ND = true
		*LCC = true
		*DC = true
	}

	flag.Parse()

	LONG_FILE_THRESHOLD = *lf
	LONG_METHOD_THRESHOLD = *lm
	NESTING_DEPTH_THRESHOLD = *nd

	argsProvided := 0

	if *d != "" {
		if !isDir(*d) {
			fmt.Fprintf(os.Stderr, "Provide a directory, Usage:\n")
			flag.PrintDefaults()
			os.Exit(1)
		}

		argsProvided++
	}
	if *e != "" {
		if !isDir(*e) {
			fmt.Fprintf(os.Stderr, "Provide a directory, Usage:\n")
			flag.PrintDefaults()
			os.Exit(1)
		}

		argsProvided++
	}

	// if argsProvided == 0 {
	// 	flag.PrintDefaults()
	// 	os.Exit(1)
	// }

	// FLAG STUFF DONE
	start := time.Now()

	if *d != "" {
		fileMetrics := findFileMetrics(*d)
		projectMetric := findProjectMetrics(fileMetrics)
		projectMetric.view(false)
		viewMetricValues(fileMetrics, *TLF, *TLM, *ND, *LCC, *DC)

	} else if *e != "" {
		dirs, err := ioutil.ReadDir(*e)
		if err != nil {
			log.Fatal(err)
		}

		var versionWiseMetrics []ProjectMetric

		for _, d := range dirs {
			if d.IsDir() {
				println("analyzing ", d.Name())
				versionPath := filepath.Join(*e, d.Name())
				fileMetrics := findFileMetrics(versionPath)

				finalMetric := findProjectMetrics(fileMetrics)
				finalMetric.VersionName = d.Name()

				versionWiseMetrics = append(versionWiseMetrics, finalMetric)
			}
		}

		viewEvolutionMetrics(versionWiseMetrics, false)
	}
	// else if *mec != "" {
	// 	dirs, err := ioutil.ReadDir(*mec)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	for _, d := range dirs {
	// 		if d.IsDir() {
	// 			path := filepath.Join(*mec, d.Name())
	// 			forEachMethod(path)
	// 		}
	// 	}
	// }

	fmt.Fprintf(os.Stderr, "\nExecution time: %s\n", time.Since(start))
}
