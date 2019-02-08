// Defines extract functions

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"os"
	"path"
	"strings"
	"unicode"
)

func checkFile(infile string) {
	if iotools.Exists(infile) == false {
		fmt.Print("\n\t[Error] Cannot find input file. Exiting.\n")
		os.Exit(1)
	}
}

func checkPunctuation(term string) bool {
	// Returns false if term contains puntuation
	for _, i := range []rune(term) {
		if i != '.' && unicode.IsPunct(i) == true {
			return false
		}
	}
	return true
}

func checkForNum(term string) bool {
	// Returns false if term contains a number
	n := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for _, i := range n {
		if strings.Contains(term, i) == true {
			return false
		}
	}
	return true
}

func checkName(query string) (string, string) {
	// Checks for formatting errors
	

}

func filter(query string) (string, string, string) {
	// Filters input query
	var term, reason string
	if len(query) >= 3 {
		term, reason = checkName(query)
		if len(reason) == 0 {
			if checkForNum(term) == false {
				reason = "numberContent"
			} else if checkPunctuation(term) == false {

			}
	}
	} else {
		reason = "tooShort"
	}
	return query, term, reason
}

func filterTerms(infile string, c int) ([][]string, [][]string) {
	// Reads terms from given column and checks formatting
	first := true
	var d string
	var pass, fail [][]string
	f := iotools.OpenFile(infile)
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		if first == false {
			s := strings.Split(line, d)
			if len(s) > c {
				query, term, reason := filter(line[c])
				// Append terms with no fail reason to pass; else append to fail
				if len(reason) == 0 {
					pass = append(pass, []string{query, term})
				} else {
					fail = append(fail, []string{query, term, reason})
				}
			}
		} else {
			d = iotool.GetDelim(line)
			first = false
		}
	}
	return pass, fail
}

func extractSearchTerms() {
	// Extracts and formats input terms
	checkFile(*infile)
	dir, _ := path.Split(*outfile)
	misses := path.Join(p, "KestrelRejected.csv")
	pass, fail := filterTerms(*infile, *column)
	fmt.Fprintln("\tSuccessfully formatted %d entries.\n\t%d entries failed formatting.", len(pass), len(fail))
	iotools.WriteCSV(*outfile, "Query,SearchTerm\n", pass)
	iotools.WriteCSV(misses, "Query,SearchTerm,Reason\n", fail)
}
