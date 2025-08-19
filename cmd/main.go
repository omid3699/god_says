package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/omid3699/god_says/cmd/server"
	"github.com/omid3699/god_says/internal"
)

func main() {
	var (
		amount = flag.Int("amount", internal.DefaultAmount, fmt.Sprintf("Number of words to generate (%d - %d)", internal.MinAmount, internal.MaxAmount))
		help   = flag.Bool("help", false, "Show the help message")
		http   = flag.Bool("http", false, "Start an HTTP server")
		host   = flag.String("host", "127.0.0.1", "The HTTP server host default is 127.0.0.1")
		port   = flag.Int("port", 3333, "The listening port of HTTP server")
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
		fmt.Fprint(os.Stderr, "\nGod Says HTTP server \n")
		fmt.Fprintf(os.Stderr, "  %s -http                    # Start HTTP server with default host and port 127.0.0.1:3333 \n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -http -host 0.0.0.0      # Start HTTP with 0.0.0.0 as host \n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -http -port 8080         # Start HTTP server listening on port 8080 \n", os.Args[0])
		os.Exit(0)
	}

	if !*http {
		// Run in CLI mode

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
	} else {
		// Run in in HTTP server mode
		log.Printf("Starting God Says HTTP server host: %s port: %d", *host, *port)
		err := server.RunServer(*host, *port)
		if err != nil {
			log.Fatalf("Error in running God Says HTTP server: %s", err)
		}
	}
}
