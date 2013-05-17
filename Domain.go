package jserv

import (
	"time"
	"os"
	"io/ioutil"
)

type Domain struct {
	hits int
	directory string
	fileCache map[string][]byte
	cacheModTimes map[string]time.Time
	default404Reply []byte
	defaultPage string
	DoCache bool
}

func NewDomain(directory string) *Domain {
	dom := new(Domain)
	dom.directory = directory
	dom.fileCache = make(map[string][]byte)
	dom.cacheModTimes = make(map[string]time.Time)
	dom.default404Reply = []byte("404 File Not Found")
	dom.defaultPage = "index.html"
	return dom
}

func (d *Domain) SetCacheOptions(opt int) {
	
}

func (d *Domain) ReloadCache() {
	for path,_ := range d.fileCache {
		fi, err := os.Open(path)
		if err != nil {
			delete(d.fileCache, path)
			delete(d.cacheModTimes, path)
		}
		d.fileCache[path],_ = ioutil.ReadAll(fi)
		d.cacheModTimes[path] = time.Now()
	}
}
