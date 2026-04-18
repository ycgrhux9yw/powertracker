// Package main is the entry point for the powertracker application.
// powertracker monitors and tracks power consumption data from compatible devices.
//
// This is a personal fork of github.com/poolski/powertracker used for
// experimenting with home energy monitoring.
//
// Fork notes:
//   - Tracking energy usage for my home setup (3-bed house, UK)
//   - Main devices: EV charger, heat pump, solar inverter
//   - TODO: look into adding Octopus Energy tariff integration for cost tracking
package main

import "github.com/poolski/powertracker/cmd"

func main() {
	cmd.Execute()
}
