// +build !nocli

package main

import "./cli"

func runClient() {
	cli.NewCLI().Loop()
}
