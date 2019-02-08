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
	column  = extract.Flag("column", "Column containing species names (integer starting from 0).").Required().Short('c').Int()

	search  = kingpin.Command("search", "Searches for taxonomy matches to extracted names.")
	firefox = search.Flag("firefox", "Use Firefox browser (uses Chrome by default).").Default("false").Bool()

	infile  = kingpin.Flag("infile", "Path to input file.").Required().String()
	outfile = kingpin.Flag("outfile", "Path to output csv file.").Required().String()
)

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
		fmt.Println("\tExtracting seacrch terms...")
		extractSearchTerms()
	case search.FullCommand():
		fmt.Println("\tSearching for taxonomy matches...")
		searchTaxonomies()
	}
	fmt.Printf("\tFinished. Run time: %s\n\n", time.Since(start))
}
