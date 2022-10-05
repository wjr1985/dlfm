package main

import (
	"flag"
	"fmt"
	"log"
)

func end(err error) {
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Press Enter to close dlfm...")
	fmt.Scanln()
}

func PrintUsage() {
	fmt.Println("dlfm - utile that shows what are you scrobbling now in your discord status")
	fmt.Println()

	flag.PrintDefaults()
}
