// Calls selenium to perform Google search

package searchtaxa

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/icwells/kestrel/src/kestrelutils"
	"github.com/icwells/kestrel/src/taxonomy"
	"github.com/tebeka/selenium"
	"log"
	"os"
	"strings"
)

func (s *searcher) parseURLs(urls map[string]string) map[string]*taxonomy.Taxonomy {
	// Attempts to find taxonomy from given urls
	taxa := make(map[string]*taxonomy.Taxonomy)
	for k, v := range urls {
		var source string
		t := taxonomy.NewTaxonomy()
		if strings.Contains(v, "#") == true {
			// Remove subheader link
			v = v[:strings.Index(v, "#")]
		}
		switch k {
		case s.urls.wiki:
			t.ScrapeWiki(v)
			source = "WIKI"
		case s.urls.itis:
			t.ScrapeItis(v)
			source = "ITIS"
		}
		if t.Found == true {
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

func (s *searcher) seleniumSearch(k string) string {
	// Gets Google search result page
	var ret string
	log.SetOutput(s.service.log)
	browser, e := s.service.getBrowser()
	if e == nil {
		defer browser.Quit()
		er := browser.Get("https://duckduckgo.com/")
		if er == nil {
			elem, err := browser.FindElement(selenium.ByName, "q")
			if err == nil {
				elem.SendKeys(kestrelutils.PercentDecode(k) + " taxonomy" + selenium.ReturnKey)
				ret, err = browser.PageSource()
				if err != nil {
					// Ensure empty return
					ret = ""
				}
			}
		}
	}
	log.SetOutput(os.Stdout)
	return ret
}

func (s *searcher) getSearchResults(k string) bool {
	// Parses urls from google search results
	found := false
	res := s.seleniumSearch(k)
	urls := s.getURLs(res)
	taxa := s.parseURLs(urls)
	if len(taxa) >= 1 {
		// Only attempt getMatch once
		found = s.getMatch(s.terms[k].Term, taxa)
	}
	return found
}
