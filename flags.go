package main

import (
	"flag"
	"github.com/mitchellh/go-homedir"
	"path"
)

var (
	HOME_DIR, _ = homedir.Dir()
	HTTPS_ADDR  = flag.String("https", ":443", "the https address to listen on")
	STORAGE     = flag.String("storage", path.Join(HOME_DIR, "httpsify/certs"), "the ssl certs storage directory")
	HOSTS_FILE  = flag.String("hosts", path.Join(HOME_DIR, "httpsify/hosts.json"), "the sites configurations filename")
)

func InitFlags() {
	flag.Parse()
}
