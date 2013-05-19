package main

import (
	"github.com/whyrusleeping/jserv"
	"net/http"
	"os"
	"time"
)

func main() {
	s := jserv.NewServer()
	//each domain object has its root directory set 
	jeromy := jserv.NewDomain("/home/whyrusleeping/jero.my/")
	info := jserv.NewDomain("/home/whyrusleeping/infopage/")

	//Start a reload daemon for jero.my site
	go func() {
		for {
			time.Sleep(time.Minute * 5)
			jeromy.ReloadCache()
		}
	}()

	//set the domain string to map the domain object to
	s.AddDomain("jero.my", jeromy)
	s.AddDomain("info.jero.my", info)

	panic(http.ListenAndServe(":80", s))
}
