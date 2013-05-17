package jserv

import (
	"net/http"
	"fmt"
	"bytes"
	"strings"
	"os"
	"time"
	"io/ioutil"
)

func ValidFilePath(url string) bool {
	if strings.Contains(url, "..") {
		return false
	}
	_, err := os.Stat(url)
	if err != nil {
		return false
	}
	return true
}

type Server struct {
	domainMap map[string]*Domain
}

func NewServer() *Server {
	s := new(Server)
	s.domainMap = make(map[string]*Domain)
	return s
}

func (s *Server) AddDomain(domStr string, domain *Domain) {
	s.domainMap[domStr] = domain
}

func (s *Server) GetDomain(domStr string) *Domain {
	d, ok := s.domainMap[domStr]
	if ok {
		return d
	}
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dom, ok := s.domainMap[r.Host]
	if !ok {
		w.Write([]byte(fmt.Sprintf("Unrecognized domain source: '%s'", r.Host)))
		return
	}
	dom.hits++
	var file string
	if r.URL.Path == "/" {
		file = dom.directory + dom.defaultPage
	} else {
		file = dom.directory + r.URL.Path[1:]
	}
	toSend, ok := dom.fileCache[file]
	if ok {
		read := bytes.NewReader(toSend)
		http.ServeContent(w,r, file, dom.cacheModTimes[file], read)
		return
	}

	//File not cached, make sure its real and send/cache it
	if !ValidFilePath(file) {
		read := bytes.NewReader(dom.default404Reply)
		http.ServeContent(w,r, "NOTFOUND.html", time.Now(), read)
		return
	}

	f, err := os.Open(file)
	if err != nil {
		w.Write([]byte("well this is awkward..."))
	}
	content,_ := ioutil.ReadAll(f)
	read := bytes.NewReader(content)
	http.ServeContent(w,r,file,time.Now(),read)
	dom.fileCache[file] = content
	dom.cacheModTimes[file] = time.Now()
}
