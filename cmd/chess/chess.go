package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("expected 'perft' or 'play' command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "perft":
		perftCommand(os.Args[2:])
	case "play":
		playCommand(os.Args[2:])
	default:
		fmt.Println("expected 'perft' or 'play' command")
		os.Exit(1)
	}
}
