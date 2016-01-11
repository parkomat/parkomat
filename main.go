package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/parkomat/parkomat/config"
	"github.com/parkomat/parkomat/dns"
	"github.com/parkomat/parkomat/web"
	"github.com/parkomat/parkomat/webdav"
	"sync"
)

func main() {
	configFile := flag.String("config_file", "parkomat.toml", "Configuration File")
	dnsOnly := flag.Bool("dns_only", false, "Run only DNS server")
	flag.Parse()

	var c *config.Config
	var err error

	c, err = config.NewConfigFromFile(*configFile)
	if err != nil {
		glog.Error("[main] Can't read config file from: ", *configFile, ". Error: ", err)
		return
	}

	var wg sync.WaitGroup

	// TODO: implement graceful shutdown
	wg.Add(1)
	go func() {
		d := dns.NewDNS(c)
		err = d.Serve()
		if err != nil {
			glog.Error("[main] DNS error: ", err)
		}
		wg.Done()
	}()

	if *dnsOnly != true {
		s := web.NewServer(c)
		dav := webdav.NewWebDav(c)

		s.Init()

		err = dav.Init()
		if err == nil {
			s.AddHandlerFunc(c.WebDav.Mount, dav.HandlerFunc)
		}

		wg.Add(2)
		go func() {
			err = s.Serve()
			if err != nil {
				glog.Error("[main] Web Error: ", err)
			}
			wg.Done()
		}()

		go func() {
			err = s.ListenAndServeTLSSNI()
			if err != nil {
				glog.Error("[mail] Web SSL Error: ", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	glog.Info("[main] Bye bye...")
}
