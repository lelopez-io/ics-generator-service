package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/lelopez-io/ics-generator-service/internal/ics"
)

// Make these variables package-level so they can be accessed in tests
var (
	serverMode *bool
	inputFile  *string
	outputFile *string
	port       *int
)

func init() {
	serverMode = flag.Bool("server", false, "Run in server mode")
	inputFile = flag.String("input", "", "Input JSON file (for local mode)")
	outputFile = flag.String("output", "output/calendar.ics", "Output ICS file (for local mode)")
	port = flag.Int("port", 8080, "Port to run the server on (for server mode)")
}

func main() {
	flag.Parse()

	if *serverMode {
		log.Printf("Starting server on port %d...\n", *port)
		http.HandleFunc("/generate-ics", ics.HandleGenerateICS)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	} else {
		if *inputFile == "" {
			log.Fatal("Input file is required in local mode. Use --input flag.")
		}
		log.Println("Running in local mode...")
		ics.GenerateLocalICS(*inputFile, *outputFile)
	}
}
