// Heler functions and structs

package kestrelutils

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"os"
	"strings"
	"unicode"
)

func CheckFile(infile string) {
	// Makes sure imut file exists
	if iotools.Exists(infile) == false {
		fmt.Printf("\n\t[Error] Input file %s not found. Exiting.\n\n", infile)
		os.Exit(1)
	}
}

func PercentEncode(term string) string {
	// Percent encodes apostrophes and spaces
	term = strings.Replace(term, " ", "%20", -1)
	return strings.Replace(term, "'", "%27", -1)
}

func PercentDecode(term string) string {
	// Removes percent encoding from web search
	term = strings.Replace(term, "%20", " ", -1)
	return strings.Replace(term, "%27", "'", -1)
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
