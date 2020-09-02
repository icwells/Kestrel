// Kestrel taxonomy web scraper

package main

import (
	"fmt"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/icwells/kestrel/src/searchtaxa"
	"github.com/icwells/kestrel/src/taxonomy"
	"github.com/icwells/kestrel/src/terms"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	app     = kingpin.New("Kestrel", "Kestrel will search online databases for taxonomy information.")
	infile  = kingpin.Flag("infile", "Path to input file.").Required().Short('i').String()
	outfile = kingpin.Flag("outfile", "Path to output csv file.").Default("").Short('o').String()

	ver = kingpin.Command("version", "Prints version info and exits.")

	format = kingpin.Command("format", "Formats new corpus for searching. New corpus is specified with the '-i' option. Output is written to the utils folder.")

	search   = kingpin.Command("search", "Searches for taxonomy matches to input names.")
	col      = search.Flag("column", "Column containing species names (integer starting from 0; use -1 for a single column file).").Default("-1").Short('c').Int()
	nocorpus = search.Flag("nocorpus", "Perform web search without searching stored corpus.").Default("false").Bool()
	proc     = search.Flag("proc", "The maximum number of concurrent processes (more will use more RAM, but will finish more quickly).").Default("200").Short('p').Int()

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
	case format.FullCommand():
		fmt.Println("\n\tFormatting new corpus...")
		taxonomy.FormatCorpus(*infile)
	case search.FullCommand():
		fmt.Println("\n\tExtracting search terms...")
		searchterms := terms.ExtractSearchTerms(*infile, *outfile, *col)
		fmt.Printf("\tCurrent run time: %v\n", time.Since(start))
		fmt.Println("\n\tSearching for taxonomy matches...")
		searchtaxa.SearchTaxonomies(*outfile, searchterms, *proc, *nocorpus)
	case merge.FullCommand():
		fmt.Println("\n\tMerging search results with source file...")
		kestrelutils.MergeResults(*infile, *resfile, *outfile, *col, *prepend)
	}
	fmt.Printf("\tFinished. Run time: %v\n\n", time.Since(start))
}
