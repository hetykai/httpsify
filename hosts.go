package main

import (
	"encoding/json"
	"os"
	"sync"
)

var (
	HOSTS  = make(map[string][]string)
	LOCKER = new(sync.Mutex)
)

func InitSites() error {
	LOCKER.Lock()
	defer LOCKER.Unlock()
	f, e := os.Open(*HOSTS_FILE)
	if e != nil {
		return e
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&HOSTS)
}
