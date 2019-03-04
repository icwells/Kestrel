// Kestrel taxonomy web scraper

package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	app = kingpin.New("Kestrel", "Kestrel will search online databases for taxonomy information.")
	ver = kingpin.Command("version", "Prints version info and exits.")

	extract = kingpin.Command("extract", "Extracts and filters input names.")
	ecol    = extract.Flag("column", "Column containing species names (integer starting from 0).").Required().Short('c').Int()

	search = kingpin.Command("search", "Searches for taxonomy matches to extracted names.")

	check    = kingpin.Command("check", "Identifies search results with matching search terms and scientific names to streamline manual curration. Give output file stem with -o.")
	taxafile = check.Flag("taxa", "Path to currated taxonomy file.").Default("nil").Short('t').String()

	merge   = kingpin.Command("merge", "Merges search results with source file.")
	prepend = merge.Flag("prepend", "Prepend taxonomies to existing rows (appends by default).").Default("false").Bool()
	mcol    = merge.Flag("names", "Column containing species names (integer starting from 0).").Required().Short('n').Int()
	resfile = merge.Flag("result", "Path to Kestrel search result file.").Required().Short('r').String()

	infile  = kingpin.Flag("infile", "Path to input file.").Required().Short('i').String()
	outfile = kingpin.Flag("outfile", "Path to output csv file.").Required().Short('o').String()
)

func version() {
	fmt.Print("\n\tKestrel v1.0 (~) is a program for resolving common names and synonyms with scientific names and extracting taxonomies.\n")
	fmt.Print("\n\tCopyright 2019 by Shawn Rupp.\n")
	fmt.Print("\tThis program comes with ABSOLUTELY NO WARRANTY.\n\tThis is free software, and you are welcome to redistribute it under certain conditions.\n")
	os.Exit(0)
}

func main() {
	start := time.Now()
	switch kingpin.Parse() {
	case ver.FullCommand():
		version()
	case extract.FullCommand():
		fmt.Println("\n\tExtracting seacrch terms...")
		extractSearchTerms()
	case search.FullCommand():
		fmt.Println("\n\tSearching for taxonomy matches...")
		searchTaxonomies(start)
	case check.FullCommand():
		fmt.Println("\n\tChecking taxonomy results...")
		checkResults()
	case merge.FullCommand():
		fmt.Println("\n\tMerging search results with source file...")
		mergeResults()
	}
	fmt.Printf("\tFinished. Run time: %v\n\n", time.Since(start))
}
