package config

import "flag"

// These flags represent all the available command line flags
var (
	ConfigFile    = flag.String("config-file", "", "Location of the config file")
	CreateAccount = flag.Bool("create", false, "If true, attempt to create account")
	DebugFlag     = flag.Bool("debug", false, "Enable debug logging")
	AccountFlag   = flag.String("account", "", "The account the CLI should connect to, if more than one is configured")
	MultiFlag     = flag.Bool("multi", false, "If true, will not try to unify the windows, but create separate instances")
)
