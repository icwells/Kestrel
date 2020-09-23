// Heler functions and structs

package kestrelutils

import (
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"os"
	"path"
	"strings"
	"unicode"
)

var (
	APOSTROPHE = "%27"
	SPACE = "%20"
)

func GetLocation() string {
	// Returns path to git repo
	return path.Join(iotools.GetGOPATH(), "src/github.com/icwells/kestrel")
}

func Getutils() string {
	// Returns path to utils directory
	return path.Join(GetLocation(), "utils")
}

func GetAbsPath(f string) string {
	// Prepends GOPATH to file name if needed
	if !strings.Contains(f, string(os.PathSeparator)) {
		f = path.Join(Getutils(), f)
	}
	if iotools.Exists(f) == false {
		fmt.Printf("\n\t[Error] Cannot find %s file. Exiting.\n", f)
		os.Exit(1)
	}
	return f
}

type Configuration struct {
	Host     string
	Database string
	User     string
	Testdb   string
	Tables   string
	Test     bool
}

func SetConfiguration(user string, test bool) Configuration {
	// Gets setting from config.txt
	var c Configuration
	c.Test = test
	c.User = user
	f := iotools.OpenFile(GetAbsPath("config.txt"))
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		s := strings.Split(string(scanner.Text()), "=")
		for idx, i := range s {
			s[idx] = strings.TrimSpace(i)
		}
		switch s[0] {
		case "host":
			c.Host = s[1]
		case "database":
			c.Database = s[1]
		case "test_database":
			c.Testdb = s[1]
		case "table_columns":
			c.Tables = GetAbsPath(s[1])
		}
	}
	return c
}

func ConnectToDatabase(user, pw string, test bool) *dbIO.DBIO {
	// Manages call to Connect and GetTableColumns
	c := SetConfiguration(user, test)
	d := c.Database
	if c.Test == true {
		d = c.Testdb
	}
	db, err := dbIO.Connect(c.Host, d, c.User, pw)
	if err != nil {
		fmt.Println(err)
		os.Exit(1000)
	}
	db.GetTableColumns()
	return db
}

func CheckFile(infile string) {
	// Makes sure imut file exists
	if iotools.Exists(infile) == false {
		fmt.Printf("\n\t[Error] Input file %s not found. Exiting.\n\n", infile)
		os.Exit(1)
	}
}

func PercentEncode(term string) string {
	// Percent encodes apostrophes and spaces
	term = strings.Replace(term, " ", SPACE, -1)
	return strings.Replace(term, "'", APOSTROPHE, -1)
}

func PercentDecode(term string) string {
	// Removes percent encoding from web search
	term = strings.Replace(term, SPACE, " ", -1)
	return strings.Replace(term, APOSTROPHE, "'", -1)
}

func RemoveKey(url string) string {
	// Returns urls with api key removed
	var idx int
	if strings.Contains(url, "&") == true {
		idx = strings.LastIndex(url, "&")
	} else if strings.Contains(url, "?") == true {
		idx = strings.LastIndex(url, "&")
	}
	if idx > 0 {
		url = url[:idx]
	}
	return url
}

func RemoveNonBreakingSpaces(s string) string {
	// Converts non-breaking spaces to standard unicode spaces
	sp := ' '
	ret := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return sp
		}
		return r
	}, s)
	return ret
}
