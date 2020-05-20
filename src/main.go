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
	app     = kingpin.New("Kestrel", "Kestrel will search online databases for taxonomy information.")
	infile  = kingpin.Flag("infile", "Path to input file.").Required().Short('i').String()
	outfile = kingpin.Flag("outfile", "Path to output csv file.").Required().Short('o').String()
	col     = search.Flag("column", "Column containing species names (integer starting from 0).").Default(0).Short('c').Int()

	ver    = kingpin.Command("version", "Prints version info and exits.")
	search = kingpin.Command("search", "Searches for taxonomy matches to input names.")

	/*check    = kingpin.Command("check", "Identifies search results with matching search terms and scientific names to streamline manual curration. Give output file stem with -o.")
	taxafile = check.Flag("taxa", "Path to currated taxonomy file.").Default("nil").Short('t').String()*/

	merge   = kingpin.Command("merge", "Merges search results with source file.")
	prepend = merge.Flag("prepend", "Prepend taxonomies to existing rows (appends by default).").Default("false").Bool()
	resfile = merge.Flag("result", "Path to Kestrel search result file.").Required().Short('r').String()
)

func version() {
	fmt.Print("\n\tKestrel is a program for resolving common names and synonyms with scientific names and extracting taxonomies.\n")
	fmt.Print("\n\tCopyright 2020 by Shawn Rupp.\n")
	fmt.Print("\tThis program comes with ABSOLUTELY NO WARRANTY.\n\tThis is free software, and you are welcome to redistribute it under certain conditions.\n")
	os.Exit(0)
}

func main() {
	start := time.Now()
	switch kingpin.Parse() {
	case ver.FullCommand():
		version()
	case search.FullCommand():
		fmt.Println("\n\tExtracting search terms...")
		searchterms := terms.ExtractSearchTerms(*infile, *outfile*col)
		fmt.Printf("\tCurrent run time: %v\n", time.Since(start))
		fmt.Println("\n\tSearching for taxonomy matches...")
		searchtaxa.SearchTaxonomies(*outfile, searchterms)
	case merge.FullCommand():
		fmt.Println("\n\tMerging search results with source file...")
		kestrelutils.MergeResults(*infile, *resfile, *outfile, *col, *prepend)
	}
	fmt.Printf("\tFinished. Run time: %v\n\n", time.Since(start))
}
