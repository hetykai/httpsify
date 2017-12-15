package main

import (
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

import (
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/roundrobin"
	"golang.org/x/crypto/acme/autocert"
)

// Initialize the autocert manager and configure it,
// also create an instance of the http.Server and link the autocert manager to it.
func InitServer() error {
	m := autocert.Manager{
		Cache:  autocert.DirCache(*STORAGE),
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			if _, ok := HOSTS[host]; ok {
				return nil
			}
			return errors.New("Unkown host(" + host + ")")
		},
	}
	s := &http.Server{
		Addr:      *HTTPS_ADDR,
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   ServeHTTP(),
	}
	log.SetOutput(ioutil.Discard)
	return s.ListenAndServeTLS("", "")
}

// The main server handler
func ServeHTTP() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if upstreams, ok := HOSTS[req.Host]; ok {
			forwarder, _ := forward.New(forward.PassHostHeader(true))
			loadbalancer, _ := roundrobin.New(forwarder)
			for _, upstream := range upstreams {
				if url, err := url.Parse(upstream); err == nil {
					loadbalancer.UpsertServer(url)
				} else {
					colorize(color.FgRed, "⇛", err.Error())
				}
			}
			res.Header().Set("X-HTTPSIFY-Version", VERSION)
			loadbalancer.ServeHTTP(res, req)
			return
		}
		http.Error(res, "The request service couldn't be found here", http.StatusNotImplemented)
	})
}
