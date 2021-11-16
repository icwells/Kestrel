// Kestrel taxonomy web scraper

package main

import (
	"bufio"
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/icwells/kestrel/src/searchtaxa"
	"github.com/icwells/kestrel/src/taxonomy"
	"github.com/icwells/kestrel/src/terms"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"strings"
	"time"
)

var (
	app     = kingpin.New("Kestrel", "Kestrel will search online databases for taxonomy information.")
	infile  = kingpin.Flag("infile", "Path to input file.").Default("").Short('i').String()
	outfile = kingpin.Flag("outfile", "Path to output csv file.").Default("").Short('o').String()
	proc    = kingpin.Flag("proc", "The maximum number of concurrent processes for search or database upload (more will use more RAM, but will finish more quickly).").Default("200").Short('p').Int()
	user    = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()

	ver = kingpin.Command("version", "Prints version info and exits.")

	upload = kingpin.Command("upload", "Formats and uploads taxonomy databases to MySQL database for searching. Databases must first be downloaded into the databases directory using './install.sh dowload'.")

	dump = kingpin.Command("dump", "Saves MySQL tables (if present) to current directory as csv files.")

	search   = kingpin.Command("search", "Searches for taxonomy matches to input names.")
	col      = search.Flag("column", "Column containing species names (integer starting from 0; use -1 for a single column file).").Default("-1").Short('c').Int()
	nocorpus = search.Flag("nocorpus", "Perform web search without searching SQL corpus.").Default("false").Bool()
	password = search.Flag("password", "MySQL password (for testing; will prompt for password by default).").String()

	merge   = kingpin.Command("merge", "Merges search results with source file.")
	prepend = merge.Flag("prepend", "Prepend taxonomies to existing rows (appends by default).").Default("false").Bool()
	resfile = merge.Flag("result", "Path to Kestrel search result file.").Required().Short('r').String()
)

func version() {
	fmt.Print("\n\tKestrel is a program for resolving common names and synonyms with scientific names and extracting taxonomies.\n")
	fmt.Print("\n\tCopyright 2021 by Shawn Rupp.\n")
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
		db = dbIO.ReplaceDatabase(c.Host, c.Database, *user, "")
		db.NewTables(c.Tables)
		// Remove types from columns map
		db.GetTableColumns()
	} else {
		fmt.Println("\tExiting.")
		os.Exit(0)
	}
	return db
}

func dumpTables(db *dbIO.DBIO, logger *log.Logger) {
	// Saves taxonomy and common name tables to current directory
	for k, v := range db.Columns {
		logger.Printf("Saving %s...\n", k)
		iotools.WriteToCSV(fmt.Sprintf("%s.csv", k), v, db.GetTable(k))
	}
}

func main() {
	var start time.Time
	var db *dbIO.DBIO
	logger := kestrelutils.GetLogger()
	switch kingpin.Parse() {
	case ver.FullCommand():
		version()
	case upload.FullCommand():
		db = newDatabase()
		start = db.Starttime
		logger.Println("Uploading taxonomies to MySQL database...")
		taxonomy.UploadDatabases(db, *proc, logger)
	case dump.FullCommand():
		db = kestrelutils.ConnectToDatabase(*user, *password, false)
		start = db.Starttime
		logger.Println("Saving taxonomy tables to current directory...")
		dumpTables(db, logger)
	case search.FullCommand():
		db = kestrelutils.ConnectToDatabase(*user, *password, false)
		logger.Println("Extracting search terms...")
		start = db.Starttime
		searchterms := terms.ExtractSearchTerms(*infile, *outfile, *col, logger)
		logger.Printf("Current run time: %v\n", time.Since(start))
		logger.Println("Searching for taxonomy matches...")
		searchtaxa.SearchTaxonomies(db, *outfile, searchterms, *proc, *nocorpus, logger)
	case merge.FullCommand():
		start = time.Now()
		logger.Println("Merging search results with source file...")
		kestrelutils.MergeResults(*infile, *resfile, *outfile, *col, *prepend, logger)
	}
	logger.Printf("Finished. Run time: %v\n\n", time.Since(start))
}
