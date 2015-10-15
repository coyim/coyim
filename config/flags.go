package config

import "flag"

var (
	ConfigFile    *string = flag.String("config-file", "", "Location of the config file")
	CreateAccount *bool   = flag.Bool("create", false, "If true, attempt to create account")
	DebugFlag     *bool   = flag.Bool("debug", false, "Enable debug logging")
)
