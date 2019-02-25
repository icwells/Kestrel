// Calls selenium to perform Google search

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/tebeka/selenium"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (s *searcher) getSearchResults(ch chan int, res, k string) {
	fmt.Println(res)
	os.Exit(0)
}

func (s *searcher) seleniumSearch(browser selenium.WebDriver, k string) string {
	// Gets Google search result page
	var ret string
	elem, err := browser.FindElement(selenium.ByName, "q")
	if err == nil {
		elem.SendKeys(percentDecode(k) + " taxonomy" + selenium.ReturnKey)
		ret, err = browser.PageSource()
		if err != nil {
			// Ensure empty return
			ret = ""
		}
	}
	return ret
}

func getDriverPath(path string) string {
	// Returns path to driver
	var ret string
	p, err := filepath.Glob(path)
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

func startService(firefox bool) (*selenium.Service, error) {
	// Initialzes new selenium browser
	var browser string
	port := 8080
	gopath := iotools.GetGOPATH()
	dir := path.Join(gopath, "src/github.com/tebeka/selenium/vendor")
	seleniumpath := path.Join(dir, "selenium-server-standalone-3.4.jar")
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(), 
		selenium.Output(os.Stderr),
	}
	if firefox == true {
		browser = "Firefox"
		gdpath := getDriverPath(path.Join(dir, "geckodriver-*"))
		opts = append(opts, selenium.GeckoDriver(gdpath))
	} else {
		browser = "Chrome"
		cdpath := getDriverPath(path.Join(dir, "chromedriver-*"))
		fmt.Println(dir)
		opts = append(opts, selenium.ChromeDriver(cdpath))
	}
	fmt.Printf("\tPerfoming Selenium search with %s browser...\n", browser)
	return selenium.NewSeleniumService(seleniumpath, port, opts...)
}

func getBrowser(firefox bool) (*selenium.Service, selenium.WebDriver, error) {
	// Returns selenium service, browser instance, and error
	var wd selenium.WebDriver
	service, err := startService(firefox)
	if err == nil {
		browser := "chrome"
		if firefox == true {
			browser = "firefox"
		}
		caps := selenium.Capabilities{"browserName": browser}
		wd, err = selenium.NewRemote(caps, "http://www.google.com")
	}
	return service, wd, err
}
