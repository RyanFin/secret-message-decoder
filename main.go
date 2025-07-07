package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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

	// extract the grid structure from the Google Doc and store in a cell table structure
	cells := parseCharacterGridFromDoc(doc)

	fmt.Println("cells: ", cells)

	// Exit if no valid data was parsed.
	if len(cells) == 0 {
		// exit the program
		log.Fatal("No valid table data found.")
	}

	renderGridFromCells(cells)

}

// parseTableFromDoc parses the first HTML table in the provided goquery.Document,
// extracting rows that contain three columns: x-coordinate, character, and y-coordinate.
//
// @param doc *goquery.Document - The parsed HTML document containing the table.
// @returns []Cell - A slice of Cell structs representing characters positioned by their x and y coordinates.
//
// The function expects the table rows (except the header) to have exactly three columns:
// - The first column is the x-coordinate (int).
// - The second column is a character (string).
// - The third column is the y-coordinate (int).
//
// Rows with invalid or missing data are skipped.
// The returned slice can be used to reconstruct a character grid based on these coordinates.
func parseCharacterGridFromDoc(doc *goquery.Document) []Cell {
	var cells []Cell

	// Find the first table in the document and iterate over its rows.
	doc.Find("table").First().Find("tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return // Skip the header row.
		}

		tds := s.Find("td")
		if tds.Length() < 3 {
			return // Skip if there aren't enough columns.
		}

		// Extract x-coordinate, character, and y-coordinate values.
		xStr := strings.TrimSpace(tds.Eq(0).Text())
		char := strings.TrimSpace(tds.Eq(1).Text())
		yStr := strings.TrimSpace(tds.Eq(2).Text())

		// Convert x and y to integers.
		x, err1 := strconv.Atoi(xStr)
		y, err2 := strconv.Atoi(yStr)

		// Skip rows with invalid integer conversion.
		if err1 != nil || err2 != nil {
			log.Printf("Skipping invalid row: %v %v\n", xStr, yStr)
			return
		}

		// Add the parsed cell to the list.
		cells = append(cells, Cell{X: x, Y: y, C: char})
	})

	return cells
}

// renderGridFromCells takes a slice of Cell structs representing characters with X,Y coordinates,
// builds a 2D grid of characters, and prints it to the console with the Y-axis flipped (bottom-to-top).
//
// @param cells []Cell - slice of Cell structs containing X, Y coordinates and character C to be placed.
// @returns none (prints output directly to stdout).
func renderGridFromCells(cells []Cell) {
	// Determine max X and Y coordinates to define grid size
	maxX, maxY := 0, 0
	for _, c := range cells {
		if c.X > maxX {
			maxX = c.X
		}
		if c.Y > maxY {
			maxY = c.Y
		}
	}

	// Initialize 2D grid filled with spaces
	grid := make([][]string, maxY+1)
	for i := range grid {
		grid[i] = make([]string, maxX+1)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	// Place each character at its (X,Y) position in the grid
	for _, c := range cells {
		grid[c.Y][c.X] = c.C
	}

	// Print the grid from bottom (maxY) to top (0), trimming trailing spaces per line
	for i := maxY; i >= 0; i-- {
		line := strings.Join(grid[i], "")
		fmt.Println(strings.TrimRight(line, " "))
	}
}
