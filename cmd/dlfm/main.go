package main

import (
	hlps "github.com/dikey0ficial/dlfm/v5/internal/helpers"
	"github.com/dikey0ficial/dlfm/v5/internal/lib"

	"log"
	"os"
)

var (
	app        lib.App
	stdl, errl *log.Logger
)

func init() {
	stdl = log.New(os.Stdout, "", log.Ltime)
	errl = log.New(os.Stderr, "[ERROR] ", log.Ltime|log.Lshortfile)

	ParseFlags()

	if flags.Help {
		PrintUsage()
		end(nil)
		os.Exit(0)
	}

	var confPath string

	if flags.ConfigPath != "" {
		if _, err := os.Stat(flags.ConfigPath); err != nil {
			errl.Printf("Specified config doesn't exist")
			end(nil)
			os.Exit(1)
		}

		confPath = flags.ConfigPath
	} else if path, err := hlps.GetConfigPath(); err != nil {
		end(err)
		os.Exit(1)
	} else {
		confPath = path
	}

	// unless it creates local app instead of using global
	var err error

	app, err = lib.NewApp(confPath, stdl, errl)

	if err != nil {
		errl.Printf("Error initializing (%v)\n", err)
		end(nil)
		os.Exit(1)
	}
}

func main() {
	code := 0
	err := app.Run()
	if err != nil {
		errl.Printf("%v\n", err)
	}
	end(nil)
	os.Exit(code)
}
