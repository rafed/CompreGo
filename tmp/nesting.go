package main

/*
 * Copyright (C) 2010 The Android Open Source Project
 */

func main() {
	P := 0
	if P == 1 {
		P = 3
	} else {
		if P == 2 {
			P = 4
		} else {
			if P == 3 {
				if P >= 3 {
					if true {
						P = 5
					}
				}
			}
		}
	}
}

type DefaultJSONParser struct {
	segments []int
}

func extract(parser DefaultJSONParser) {

	for i := 0; i < len(parser.segments); i++ {

		if i != 0 {
			continue
		}

		// var eval bool
		eval := true
		println(eval)

		last := true
		if !last {
			// nextSegment := segments[i+1]
			if i > 0 {
				eval = true
			} else if i >= 0 && i < 3 {
				eval = true
			} else if i > 5 {
				eval = true
			} else {
				eval = false
			}
		} else {
			eval = true
		}
	}

	// if eval == true println("hello")
}

/**
 * dumped
 */
func executePendingTransactions() bool {
	b := false
	return b
}

/**
 * dumped
 */
func executeTransactions() bool {
	b := false
	return b
}
