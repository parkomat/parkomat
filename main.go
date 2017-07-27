package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/parkomat/parkomat/config"
	"github.com/parkomat/parkomat/dns"
	"github.com/parkomat/parkomat/web"
	"github.com/parkomat/parkomat/webdav"
	"os"
	"strconv"
	"sync"
)

func main() {
	log.WithFields(log.Fields{
		"service": "main",
	}).Info("Parkomat (parkomat.io)")

	configFile := flag.String("config_file", "parkomat.toml", "Configuration File")
	dnsOnly := flag.Bool("dns_only", false, "Run only DNS server")
	flag.Parse()

	var c *config.Config
	var err error

	// If you specify environment variable, args will be overwritten
	envConfigFile := os.Getenv("PARKOMAT_CONFIG_FILE")
	if envConfigFile != "" {
		configFile = &envConfigFile
	}

	envDnsOnly := os.Getenv("PARKOMAT_DNS_ONLY")
	if envDnsOnly != "" {
		if s, err := strconv.ParseBool(envDnsOnly); err == nil {
			dnsOnly = &s
		}
	}

	c, err = config.NewConfigFromFile(*configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"service": "main",
			"path":    *configFile,
			"error":   err,
		}).Error("Can't read config file")
		return
	}

	var wg sync.WaitGroup

	// TODO: implement graceful shutdown
	wg.Add(1)
	go func() {
		d := dns.NewDNS(c)
		err = d.Serve("udp")
		if err != nil {
			log.WithFields(log.Fields{
				"service": "main",
				"error":   err,
			}).Error("DNS UDP Error")
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		d := dns.NewDNS(c)
		err = d.Serve("tcp")
		if err != nil {
			log.WithFields(log.Fields{
				"service": "main",
				"error":   err,
			}).Error("DNS TCP Error")
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
				log.WithFields(log.Fields{
					"service": "main",
					"error":   err,
				}).Error("Web Error")
			}
			wg.Done()
		}()

		go func() {
			err = s.ListenAndServeTLSSNI()
			if err != nil {
				log.WithFields(log.Fields{
					"service": "main",
					"error":   err,
				}).Error("Web SSL Error")
			}
			wg.Done()
		}()
	}

	wg.Wait()
	log.WithFields(log.Fields{
		"service": "main",
	}).Info("Exit")
}
