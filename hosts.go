package main

import (
	"encoding/json"
	"os"
	"sync"
)

import (
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

var (
	HOSTS  = make(map[string][]string)
	LOCKER = new(sync.Mutex)
)

// Load the hosts file into memory
func InitHostsList() error {
	LOCKER.Lock()
	defer LOCKER.Unlock()
	f, e := os.OpenFile(*HOSTS_FILE, os.O_RDONLY|os.O_CREATE, 0666)
	if e != nil {
		return e
	}
	defer f.Close()
	finfo, err := f.Stat()
	if err != nil {
		return err
	}
	if finfo.Size() == 0 {
		return nil
	}
	return json.NewDecoder(f).Decode(&HOSTS)
}

// watch the hosts file for any changes and the reload it
func WatchChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()
	watcher.Add(*HOSTS_FILE)
	for {
		select {
		case <-watcher.Events:
			colorize(color.FgYellow, "⇛ There is a change in the hosts file, reloading ...")
			if err := InitHostsList(); err != nil {
				colorize(color.FgRed, "⇛", err.Error())
			} else {
				colorize(color.FgGreen, "⇛ The hosts file has been reloaded successfully!")
			}
		}
	}
}
