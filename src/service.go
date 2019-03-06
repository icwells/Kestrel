// Defines service struct for selenium

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/tebeka/selenium"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type service struct {
	service *selenium.Service
	err     error
	log     *os.File
	port    int
	browser string
}

func (s *service) startService() {
	// Initialzes new selenium browser
	s.log = iotools.CreateFile("seleniumLog.txt")
	gopath := iotools.GetGOPATH()
	dir := path.Join(gopath, "src/github.com/tebeka/selenium/vendor")
	seleniumpath := getSeleniumPath(dir)
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),
		selenium.Output(s.log),
		selenium.ChromeDriver(getDriverPath(path.Join(dir, "chromedriver-*"))),
	}
	s.service, s.err = selenium.NewSeleniumService(seleniumpath, s.port, opts...)
}

func (s *service) stop() {
	// Closes service
	s.service.Stop()
	s.log.Close()
}

func newService() service {
	var s service
	s.port = 8080
	s.browser = "chrome"
	s.startService()
	return s
}

func (s *service) getBrowser() (selenium.WebDriver, error) {
	// Returns browser instance and error
	return selenium.NewRemote(selenium.Capabilities{"browserName": s.browser}, fmt.Sprintf("http://localhost:%d/wd/hub", s.port))
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
