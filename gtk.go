// +build nocli

package main

import "../gui"

func runClient() {
	gui.NewGTK().Loop()
}
