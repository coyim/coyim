package config

import "flag"

// These flags represent all the available command line flags
var (
	ConfigFile    = flag.String("config-file", "", "Location of the config file")
	CreateAccount = flag.Bool("create", false, "If true, attempt to create account")
	DebugFlag     = flag.Bool("debug", false, "Enable debug logging")
)
