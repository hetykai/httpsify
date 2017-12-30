package main

import (
	"flag"
	"os"
	"path"
)

import "github.com/mitchellh/go-homedir"

var (
	VERSION     = "httpsify/v3.0.0"
	HOME_DIR, _ = homedir.Dir()
	HTTPS_ADDR  = flag.String("https", ":443", "the https address to listen on")
	STORAGE     = flag.String("storage", path.Join(HOME_DIR, ".httpsify/certs"), "the ssl certs storage directory")
	HOSTS_FILE  = flag.String("hosts", path.Join(HOME_DIR, ".httpsify/hosts.json"), "the sites configurations filename")
	HSTS        = flag.String("hsts", "max-age=86400; includeSubDomains", "the hsts header value, empty value means disable")
	EXPOSE_INFO = flag.Bool("expose-info", true, "whether to expose the httpsify info header or not")
)

func InitFlags() {
	flag.Parse()
	os.MkdirAll(*STORAGE, 0755)
}
