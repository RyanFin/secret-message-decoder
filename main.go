package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

// Cell represents a character to be drawn at a specific (X, Y) coordinate.
type Cell struct {
	X int    // X-coordinate (column)
	Y int    // Y-coordinate (row)
	C string // Character to draw
}

func main() {
	// Ensure a URL is provided as a command-line argument.
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <public_google_doc_url>")
	}
	url := os.Args[1]

	// Fetch the content of the Google Doc.
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch document: %v", err)
	}
	defer resp.Body.Close()

	// Check for successful HTTP response.
	if resp.StatusCode != 200 {
		log.Fatalf("Failed to fetch document: HTTP %d", resp.StatusCode)
	}

	// Parse the HTML content of the document.
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse document HTML: %v", err)
	}

	// Debug print to verify the parsed document (can be removed later).
	fmt.Println("document: ", doc)

}
