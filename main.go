package main

import (
	"fmt"
	"sd-ingest/cmd"
)

var (
	commit  string
	builtAt string
	version = "0.1"
)

func main() {

	fmt.Printf("\n\nSD Ingest Version %s\n", version)
	fmt.Println("by Jason Underhill\n")
	fmt.Println("https://github.com/junderhill/sd-ingest\n")
	fmt.Printf("Built from %s at %s\n\n", commit, builtAt)

	cmd.Execute()
}
