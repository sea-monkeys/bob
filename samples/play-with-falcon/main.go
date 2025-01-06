package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Command flags
	namePtr := flag.String("name", "World", "The name to greet or say goodbye to")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  sam --version\n")
		fmt.Fprintf(os.Stderr, "  sam --help\n")
		fmt.Fprintf(os.Stderr, "  sam hello\n")
		fmt.Fprintf(os.Stderr, "  sam hello --name <name>\n")
		fmt.Fprintf(os.Stderr, "  sam goodbye\n")
		fmt.Fprintf(os.Stderr, "  sam goodbye --name <name>\n")
		flag.PrintDefaults()
	}

	// Parse command-line arguments
	flag.Parse()

	// Check command and execute
	switch flag.Arg(0) {
	case "hello":
		fmt.Printf("Hello, %s!\n", *namePtr)
	case "goodbye":
		fmt.Printf("Goodbye, %s!\n", *namePtr)
	case "--version":
		fmt.Println("sam version 1.0")
	case "--help":
		flag.Usage()
	default:
		flag.Usage()
	}
}