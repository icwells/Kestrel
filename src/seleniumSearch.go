// Calls selenium to perform Google search

package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/icwells/go-tools/iotools"
	"github.com/tebeka/selenium"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func (s *searcher) parseURLs(urls map[string]string) map[string]taxonomy {
	// Attempts to find taxonomy from given urls
	taxa := make(map[string]taxonomy)
	for k, v := range urls {
		var source string
		t := newTaxonomy()
		if strings.Contains(v, "#") == true {
			// Remove subheader link
			v = v[:strings.Index(v, "#")]
		}
		switch k {
		/*case s.urls.wiki:
		t.scrapeWiki(v)
		source = "WIKI"*/
		case s.urls.itis:
			t.scrapeItis(v)
			source = "ITIS"
		}
		if t.found == true {
			taxa[source] = t
		}
	}
	return taxa
}

func (s *searcher) getURLs(res string) map[string]string {
	// Returns slice os urls to scrape
	ret := make(map[string]string)
	page, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err == nil {
		page.Find("a").Each(func(i int, r *goquery.Selection) {
			// Examine all attach tags for target links
			url, ex := r.Attr("href")
			if ex == true && strings.Count(url, ":") <= 1 && strings.Contains(url, "(") == false {
				// Skip urls from webcaches and disambiguation pages
				for _, i := range []string{s.urls.wiki, s.urls.itis} {
					if strings.Contains(url, i) == true {
						if _, exists := ret[i]; exists == false {
							ret[i] = url
						}
						break
					}
				}
			}
		})
	}
	return ret
}

func (s *searcher) getSearchResults(wg *sync.WaitGroup, mut *sync.RWMutex, res, k string) {
	// Parses urls from google search results
	defer wg.Done()
	found := false
	urls := s.getURLs(res)
	taxa := s.parseURLs(urls)
	if len(taxa) >= 1 {
		// Only attempt getMatch once
		found = s.getMatch(s.terms[k].term, taxa)
	}
	mut.Lock()
	if found == true {
		s.writeMatches(k)
	} else {
		// Write missed queries to file
		s.writeMisses(k)
	}
	mut.Unlock()
}

func (s *searcher) seleniumSearch(browser selenium.WebDriver, k string) string {
	// Gets Google search result page
	var ret string
	er := browser.Get("http://www.google.com")
	if er == nil {
		elem, err := browser.FindElement(selenium.ByName, "q")
		if err == nil {
			elem.SendKeys(percentDecode(k) + " taxonomy" + selenium.ReturnKey)
			ret, err = browser.PageSource()
			if err != nil {
				// Ensure empty return
				ret = ""
			}
		}
	}
	return ret
}

func getDriverPath(dir string) string {
	// Returns path to driver
	var ret string
	p, err := filepath.Glob(dir)
	if err == nil {
		for _, i := range p {
			if strings.Contains(i, ".zip") == false && strings.Contains(i, ".tar") == false {
				if iotools.Exists(i) == true {
					ret = i
					break
				}
			}
		}
	}
	return ret
}

func getSeleniumPath(dir string) string {
	// Returns path to selenium jar
	var ret string
	p, err := filepath.Glob(path.Join(dir, "selenium-server-standalone-*"))
	if err == nil {
		if len(p) > 1 {
			// Get highest version number
			ver := 0.0
			for _, i := range p {
				n := i[strings.LastIndex(i, "-")+1 : strings.LastIndex(i, ".")]
				if strings.Count(n, ".") > 1 {
					n = n[:strings.LastIndex(n, ".")]
				}
				v, er := strconv.ParseFloat(n, 64)
				if er == nil && v > ver {
					ver = v
					ret = i
				}
			}
		} else if len(p) == 1 {
			ret = p[0]
		}
		if iotools.Exists(ret) == false {
			ret = ""
		}
	}
	return ret
}

func startService(port int, browser string) (*selenium.Service, error) {
	// Initialzes new selenium browser
	gopath := iotools.GetGOPATH()
	dir := path.Join(gopath, "src/github.com/tebeka/selenium/vendor")
	seleniumpath := getSeleniumPath(dir)
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),
		selenium.Output(os.Stderr),
	}
	cdpath := getDriverPath(path.Join(dir, "chromedriver-*"))
	opts = append(opts, selenium.ChromeDriver(cdpath))
	fmt.Printf("\tPerfoming Selenium search with %s browser...\n\n", browser)
	return selenium.NewSeleniumService(seleniumpath, port, opts...)
}

func getBrowser() (*selenium.Service, selenium.WebDriver, error) {
	// Returns selenium service, browser instance, and error
	var wd selenium.WebDriver
	port := 8080
	browser := "chrome"
	service, err := startService(port, browser)
	if err == nil {
		caps := selenium.Capabilities{"browserName": browser}
		wd, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	}
	// Println blank line after selenium stdout
	fmt.Println()
	return service, wd, err
}
