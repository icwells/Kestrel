// Kestrel taxonomy web scraper

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	app = kingpin.New("Kestrel", "Kestrel will search online databases for taxonomy information.")
	ver = kingpin.Command("version", "Prints version info and exits.")

	extract = kingpin.Command("extract", "Extracts and filters input names.")

	search  = kingpin.Command("search", "Searches for taxonomy matches to extracted names.")
	firefox = search.Flag("firefox", "Use Firefox browser (uses Chrome by default).").Default("false").Bool()

	check = kingpin.Command("check", "Identifies search results with matching search terms and scientific names to streamline manual curration. Give output file stem with -o.")

	merge   = kingpin.Command("merge", "Merges search results with source file.")
	resfile = merge.Flag("result", "Path to Kestrel search result file.").Required().Short('r').String()

	infile  = kingpin.Flag("infile", "Path to input file.").Required().Short('i').String()
	outfile = kingpin.Flag("outfile", "Path to output csv file.").Required().Short('o').String()
	column  = kingpin.Flag("column", "Column containing species names (integer starting from 0; required for extract and merge).").Default("-1").Short('c').Int()
)

func checkColumn(c int) {
	// Quits if column wasn't given
	if c < 0 {
		fmt.Print("\n\t[Error] Column number not given. Exiting.\n\n")
		os.Exit(2)
	}
}

func checkFile(infile string) {
	// Makes sure imut file exists
	if iotools.Exists(infile) == false {
		fmt.Printf("\n\t[Error] Input file %s not found. Exiting.\n\n", infile)
		os.Exit(1)
	}
}

func version() {
	fmt.Print("\n\tKestrel v1.0 (~) is a program for resolving common names and synonyms with scientific names and extracting taxonomies.\n")
	fmt.Print("\n\tCopyright 2019 by Shawn Rupp.\n")
	fmt.Print("\tThis program comes with ABSOLUTELY NO WARRANTY.\n\tThis is free software, and you are welcome to redistribute it under certain conditions.\n")
	os.Exit(0)
}

func main() {
	var start time.Time
	switch kingpin.Parse() {
	case ver.FullCommand():
		version()
	case extract.FullCommand():
		fmt.Println("\n\tExtracting seacrch terms...")
		extractSearchTerms()
	case search.FullCommand():
		fmt.Println("\n\tSearching for taxonomy matches...")
		searchTaxonomies()
	case check.FullCommand():
		fmt.Println("\n\tChecking taxonomy results...")
		checkResults()
	case merge.FullCommand():
		fmt.Println("\n\tMerging search results with source file...")
		mergeResults()
	}
	fmt.Printf("\tFinished. Run time: %s\n\n", time.Since(start))
}
