// Kestrel taxonomy web scraper

package main

import (
	"fmt"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/icwells/kestrel/src/searchtaxa"
	"github.com/icwells/kestrel/src/terms"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	app = kingpin.New("Kestrel", "Kestrel will search online databases for taxonomy information.")
	infile  = kingpin.Flag("infile", "Path to input file.").Required().Short('i').String()
	outfile = kingpin.Flag("outfile", "Path to output csv file.").Required().Short('o').String()

	ver = kingpin.Command("version", "Prints version info and exits.")

	search = kingpin.Command("search", "Searches for taxonomy matches to input names.")
	col    = search.Flag("column", "Column containing species names (integer starting from 0).").Required().Short('c').Int()

	/*check    = kingpin.Command("check", "Identifies search results with matching search terms and scientific names to streamline manual curration. Give output file stem with -o.")
	taxafile = check.Flag("taxa", "Path to currated taxonomy file.").Default("nil").Short('t').String()

	merge   = kingpin.Command("merge", "Merges search results with source file.")
	prepend = merge.Flag("prepend", "Prepend taxonomies to existing rows (appends by default).").Default("false").Bool()
	resfile = merge.Flag("result", "Path to Kestrel search result file.").Required().Short('r').String()*/
)

func version() {
	fmt.Print("\n\tKestrel is a program for resolving common names and synonyms with scientific names and extracting taxonomies.\n")
	fmt.Print("\n\tCopyright 2020 by Shawn Rupp.\n")
	fmt.Print("\tThis program comes with ABSOLUTELY NO WARRANTY.\n\tThis is free software, and you are welcome to redistribute it under certain conditions.\n")
	os.Exit(0)
}

func taxonomySearch(start *time.Time) {
	// Wraps calls for taxonomy searches
	searchterms := terms.ExtractSearchTerms(*infile, *outfile *col)
	fmt.Printf("\tCurrent run time: %v\n", time.Since(start))
	searchtaxa.SearchTaxonomies(*outfile, searchterms)
}

func main() {
	start := time.Now()
	switch kingpin.Parse() {
	case ver.FullCommand():
		version()
	case search.FullCommand():
		taxonomySearch(start)
	case check.FullCommand():
		fmt.Println("\n\tChecking taxonomy results...")
		kestrelutils.CheckResults()
	case merge.FullCommand():
		fmt.Println("\n\tMerging search results with source file...")
		kestrelutils.MergeResults()
	}
	fmt.Printf("\tFinished. Run time: %v\n\n", time.Since(start))
}
