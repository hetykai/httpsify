package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"golang.org/x/crypto/acme/autocert"
)

var (
	// CMD options
	listen      = flag.String("listen", ":443", "the local listen address")
	domains     = flag.String("domains", "", "a comma separated strings of domain[->[ip]:port]")
	backend     = flag.String("backend", ":80", "the default backend to be used")
	sslCacheDir = flag.String("ssl-cache-dir", "./httpsify-ssl-cache", "the cache directory to cache generated ssl certs")

	// internal vars
	domain_backend = map[string]string{}
	whitelisted    = []string{}
)

func main() {
	flag.Parse()

	if *domains == "" {
		flag.Usage()
		fmt.Println(`Example(1): httpsify -domains "example.org,api.example.org->localhost:366, api2.example.org->:367"`)
		fmt.Println(`Example(2): httpsify -domains "www.site.com,apiv1.site.com->:8080,apiv2.site.com->:8081"`)
		return
	}

	for _, zone := range strings.Split(*domains, ",") {
		parts := strings.SplitN(zone, "->", 2)
		if len(parts) < 2 {
			parts = append(parts, *backend)
		}
		parts[1] = fixUrl(parts[1])
		domain_backend[parts[0]] = parts[1]
		whitelisted = append(whitelisted, parts[0])
	}

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(whitelisted...),
		Cache:      autocert.DirCache(*sslCacheDir),
	}

	h := handler()

	s := &http.Server{
		Addr:      *listen,
		Handler:   h,
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	}

	log.Fatal(s.ListenAndServeTLS("", ""))
}

// fix the specified url
// this function will make sure that "http://" already exists,
// also it will make sure that it has a hostname .
func fixUrl(u string) string {
	u = strings.TrimPrefix(strings.TrimSpace(u), "https://")
	if strings.Index(u, ":") == 0 {
		u = "localhost" + u
	}
	if !strings.HasPrefix(u, "ws://") && !strings.HasPrefix(u, "http://") {
		u = "http://" + u
	}
	u = strings.TrimRight(u, "/")
	return u
}

// the proxy handler
func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = strings.SplitN(r.Host, ":", 2)[0]
		if _, found := domain_backend[r.Host]; !found {
			http.Error(w, r.Host+": not found", http.StatusNotImplemented)
			return
		}
		r.Header["X-Forwarded-Proto"] = []string{"https"}
		r.Header["X-Forwarded-For"] = append(r.Header["X-Forwarded-For"], strings.SplitN(r.RemoteAddr, ":", 2)[0])
		u, _ := url.Parse(domain_backend[r.Host] + "/" + strings.TrimLeft(r.URL.RequestURI(), "/"))
		if strings.ToLower(r.Header.Get("Upgrade")) == "websocket" {
			NewWebsocketReverseProxy(u).ServeHTTP(w, r)
			return
		} else {
			proxy := httputil.NewSingleHostReverseProxy(u)
			defaultDirector := proxy.Director
			proxy.Director = func(req *http.Request) {
				defaultDirector(req)
				req.Host = r.Host
				req.URL = u
			}
			proxy.ServeHTTP(w, r)
			return
		}
	})
}

// the websocket proxy handler
func NewWebsocketReverseProxy(u *url.URL) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		backConn, err := net.Dial("tcp", u.Host)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer backConn.Close()
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}
		clientConn, _, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer clientConn.Close()
		message := r.Method + " " + r.URL.RequestURI() + " " + r.Proto + "\n"
		message += "Host: " + r.Host + "\n"
		for k, vals := range r.Header {
			for _, v := range vals {
				message += k + ": " + v + "\n"
			}
		}
		message += "\n"
		go io.Copy(backConn, io.MultiReader(strings.NewReader(message), r.Body, clientConn))
		io.Copy(clientConn, backConn)
	})
}
