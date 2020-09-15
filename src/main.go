// Kestrel taxonomy web scraper

package main

import (
	"bufio"
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/icwells/kestrel/src/searchtaxa"
	"github.com/icwells/kestrel/src/taxonomy"
	"github.com/icwells/kestrel/src/terms"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
	"time"
)

var (
	app     = kingpin.New("Kestrel", "Kestrel will search online databases for taxonomy information.")
	infile  = kingpin.Flag("infile", "Path to input file.").Default("").Short('i').String()
	outfile = kingpin.Flag("outfile", "Path to output csv file.").Default("").Short('o').String()
	user    = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Default("root").String()

	ver = kingpin.Command("version", "Prints version info and exits.")

	upload = kingpin.Command("upload", "Formats and uploads taxonomy databases to MySQL database for searching. Databases must first be downloaded into the databases directory using './install.sh dowload'.")

	search   = kingpin.Command("search", "Searches for taxonomy matches to input names.")
	col      = search.Flag("column", "Column containing species names (integer starting from 0; use -1 for a single column file).").Default("-1").Short('c').Int()
	nocorpus = search.Flag("nocorpus", "Perform web search without searching SQL corpus.").Default("false").Bool()
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

func newDatabase() *dbIO.DBIO {
	// Creates new database and tables
	var db *dbIO.DBIO
	c := kestrelutils.SetConfiguration(*user, false)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n\tAre you sure you want to initialize a new database? This will erase existing data. (y|n) ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(strings.ToLower(text))
	if text == "y" || text == "yes" {
		db = dbIO.CreateDatabase(c.Host, c.Database, *user)
		db.NewTables(c.Tables)
	} else {
		fmt.Println("\tExiting.")
		os.Exit(0)
	}
	return db
}

func main() {
	var start time.Time
	var db *dbIO.DBIO
	switch kingpin.Parse() {
	case ver.FullCommand():
		version()
	case upload.FullCommand():
		db = newDatabase()
		start = db.Starttime
		fmt.Println("\n\tUploading taxonomies to MySQL database...")
		taxonomy.UploadDatabases(db)
	case search.FullCommand():
		db = kestrelutils.ConnectToDatabase(*user, false)
		fmt.Println("\n\tExtracting search terms...")
		start = db.Starttime
		searchterms := terms.ExtractSearchTerms(*infile, *outfile, *col)
		fmt.Printf("\tCurrent run time: %v\n", time.Since(start))
		fmt.Println("\n\tSearching for taxonomy matches...")
		searchtaxa.SearchTaxonomies(db, *outfile, searchterms, *proc, *nocorpus)
	case merge.FullCommand():
		start = time.Now()
		fmt.Println("\n\tMerging search results with source file...")
		kestrelutils.MergeResults(*infile, *resfile, *outfile, *col, *prepend)
	}
	fmt.Printf("\tFinished. Run time: %v\n\n", time.Since(start))
}
