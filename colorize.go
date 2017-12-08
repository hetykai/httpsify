package main

import "fmt"
import "github.com/fatih/color"

// colorize the log output
func colorize(c color.Attribute, args ...interface{}) {
	color.Set(c)
	fmt.Println(args...)
	color.Unset()
}
