package main

import (
	"log"

	"github.com/simonhayward/gonix"
)

func main() {
	if err := gonix.Run(); err != nil {
		log.Fatal(err)
	}
}
