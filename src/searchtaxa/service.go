// Defines service struct for selenium

package searchtaxa

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"os"
	"os/exec"
	"path"
)

type service struct {
	browser string
	driver  string
	err     error
	ip      string
	log     *os.File
	port    int
	service *selenium.Service
}

func newService() *service {
	// Initializes new struct
	s := new(service)
	s.browser = "chrome"
	s.driver = "chromedriver"
	s.ip = "http://127.0.0.1"
	s.port = 8090
	s.startService()
	return s
}

func (s *service) getBrowser() (selenium.WebDriver, error) {
	// Returns browser instance and error
	caps := selenium.Capabilities{"browserName": s.browser,
		"pageLoadStrategy": "normal",
	}
	opt := new(chrome.Capabilities)
	opt.W3C = false
	caps.AddChrome(*opt)
	return selenium.NewRemote(caps, fmt.Sprintf("%s:%d/wd/hub", s.ip, s.port))
}

func (s *service) startService() {
	// Initialzes new selenium service
	var err error
	dir := path.Join(iotools.GetGOPATH(), "src/github.com/tebeka/selenium/vendor")
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),
		selenium.Output(s.log),
		selenium.ChromeDriver(path.Join(dir, "chromedriver")),
	}
	s.service, err = selenium.NewSeleniumService(path.Join(dir, "selenium-server.jar"), s.port, opts...)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *service) stop() {
	// Closes service
	s.service.Stop()
	// Flush log before closing
	s.log.Sync()
	s.log.Close()
}

func (s *service) KillChromeDrivers() {
	// Kills all remaining chromedriver processes
	for _, i := range []string{s.driver, s.browser} {
		cmd := exec.Command("killall", "-q", i)
		cmd.Run()
	}
}
