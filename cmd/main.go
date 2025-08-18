package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/omid3699/god_says/internal"
)

func main() {
	var (
		amount = flag.Int("amount", internal.DefaultAmount, fmt.Sprintf("Number of words to generate (%d - %d)", internal.MinAmount, internal.MaxAmount))
		help   = flag.Bool("help", false, "Show the help message")
	)
	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nGo port of Terry Davis' \"god says\" program from TempleOS\n")
		fmt.Fprintf(os.Stderr, "Generates random words from the Happy.TXT wordlist.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                    # Generate %d words (default)\n", os.Args[0], internal.DefaultAmount)
		fmt.Fprintf(os.Stderr, "  %s -amount 10         # Generate 10 words\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -amount 100        # Generate 100 words\n", os.Args[0])
		os.Exit(0)
	}

	if *amount < internal.MinAmount || *amount > internal.MaxAmount {
		fmt.Fprintf(os.Stderr, "Error: amount must be between %d and %d\n", internal.MinAmount, internal.MaxAmount)
		os.Exit(1)
	}

	god, err := internal.NewGod(*amount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialize God Says %v \n", err)
	}

	message := god.Speak()
	fmt.Println(message)
}
