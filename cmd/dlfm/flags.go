package main

import (
	"flag"
)

type Flags struct {
	ConfigPath string
	Help       bool
}

var flags Flags

func ParseFlags() {
	flag.StringVar(&flags.ConfigPath, "conf", "", "uses specified `config`")
	flag.StringVar(&flags.ConfigPath, "c", "", "alias for --conf")

	flag.BoolVar(&flags.Help, "help", false, "shows help and exits")
	flag.BoolVar(&flags.Help, "h", false, "alias for --help")

	flag.Parse()
}
