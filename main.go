// Package main is the entry point for the powertracker application.
// powertracker monitors and tracks power consumption data from compatible devices.
//
// This is a personal fork of github.com/poolski/powertracker used for
// experimenting with home energy monitoring.
package main

import "github.com/poolski/powertracker/cmd"

func main() {
	cmd.Execute()
}
