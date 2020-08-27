package gui

import "runtime"

type supportedSymbols struct {
	bullet string
}

var symbols supportedSymbols

func initUnicodeSymbols() {
	if runtime.GOOS == "windows" {
		symbols = windowsSymbols()
	} else {
		symbols = defaultSymbols()
	}
}

func windowsSymbols() supportedSymbols {
	return supportedSymbols{
		bullet: "*",
	}
}

func defaultSymbols() supportedSymbols {
	return supportedSymbols{
		bullet: "•",
	}
}
