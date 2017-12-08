package main

import "github.com/fatih/color"

func main() {
	colorize(color.FgGreen, "⇛ Parsing the specified flags ...")
	InitFlags()

	colorize(color.FgGreen, "⇛ Loading the provided json hosts file from ("+*HOSTS_FILE+") ...")
	if err := InitHostsList(); err != nil {
		colorize(color.FgRed, "⇛", err.Error())
		return
	}

	colorize(color.FgGreen, "⇛ Watching the hosts file for any change to hot reload it ...")
	go func() {
		WatchChanges()
	}()

	colorize(color.FgGreen, "⇛ Running the HTTPS (HTTP/2) server on address ("+*HTTPS_ADDR+") ...")
	colorize(color.FgRed, InitServer())
}
