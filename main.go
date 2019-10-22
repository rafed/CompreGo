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

func main() {

	lf := flag.Int("lf", 750, "Long File Threshold (deafult 750)")
	lm := flag.Int("lm", 75, "Long Method Threshold (default 75)")

	d := flag.String("d", "", "maintainibility metrics of a project")
	e := flag.String("e", "", "evolution of maintainibility metrics for each version")

	flag.Parse()

	LONG_FILE_THRESHOLD = *lf
	LONG_METHOD_THRESHOLD = *lm

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

	if argsProvided > 1 || argsProvided == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// // FLAG STUFF DONE
	start := time.Now()

	if *d != "" {
		fileMetrics := findFileMetrics(*d)
		projectMetric := findProjectMetrics(fileMetrics)
		projectMetric.view()
		// viewProjectMetrics(structs)
	} else if *e != "" {
		dirs, err := ioutil.ReadDir(*e)
		if err != nil {
			log.Fatal(err)
		}

		var versionWiseMetrics []ProjectMetric

		for _, d := range dirs {
			if d.IsDir() {
				versionPath := filepath.Join(*e, d.Name())
				fileMetrics := findFileMetrics(versionPath)

				finalMetric := findProjectMetrics(fileMetrics)
				finalMetric.VersionName = d.Name()

				versionWiseMetrics = append(versionWiseMetrics, finalMetric)
			}
		}

		viewEvolutionMetrics(versionWiseMetrics)
	}

	fmt.Fprintf(os.Stderr, "Execution time: %s\n", time.Since(start))
}
