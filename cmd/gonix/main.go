package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/simonhayward/gonix"
)

// Injected by Nix
var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

func main() {
	versionFlag := flag.Bool("version", false, "print the version info and exit")

	flag.Parse()

	if *versionFlag {
		ts, _ := strconv.ParseInt(BuildTime, 10, 64)
		prettyDate := time.Unix(ts, 0).UTC().Format("2006-01-02 15:04:05 UTC")

		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Commit:     %s\n", Commit)
		fmt.Printf("Built At:   %s\n", prettyDate)
		os.Exit(0)
	}

	if err := gonix.Run(); err != nil {
		log.Fatal(err)
	}
}
