package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shreybhardwaj/zizou/internal/diff"
)

func main() {
	// Read the test diff file
	content, err := os.ReadFile("test_diff.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Parse it
	parser := diff.NewParser()
	d, err := parser.Parse(string(content))
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	// Display results
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("           ZIZOU DIFF PARSER - TEST RUN")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Statistics
	stats := d.Stats()
	fmt.Printf("📊 STATISTICS:\n")
	fmt.Printf("   Files changed:  %d\n", stats.Files)
	fmt.Printf("   Lines added:    +%d\n", stats.LinesAdded)
	fmt.Printf("   Lines removed:  -%d\n", stats.LinesRemoved)
	fmt.Printf("   Net change:     %+d\n\n", stats.LinesAdded-stats.LinesRemoved)

	// Modified files
	fmt.Printf("📁 MODIFIED FILES:\n")
	for _, file := range d.GetModifiedFiles() {
		fmt.Printf("   • %s\n", file)
	}
	fmt.Println()

	// Detailed analysis
	fmt.Printf("🔍 DETAILED ANALYSIS:\n\n")
	for _, file := range d.Files {
		fmt.Printf("File: %s\n", file.NewPath)
		fmt.Printf("Hunks: %d\n\n", len(file.Hunks))

		for i, hunk := range file.Hunks {
			fmt.Printf("  Hunk #%d: @@ -%d,%d +%d,%d @@\n",
				i+1, hunk.OldStart, hunk.OldCount, hunk.NewStart, hunk.NewCount)

			addedCount := 0
			removedCount := 0
			for _, line := range hunk.Lines {
				if line.Type == diff.LineAdded {
					addedCount++
				} else if line.Type == diff.LineRemoved {
					removedCount++
				}
			}

			fmt.Printf("  Changes: +%d -%d\n\n", addedCount, removedCount)
		}
	}

	// Show added lines
	fmt.Printf("✅ ADDED LINES:\n")
	for _, line := range d.GetAddedLines() {
		fmt.Printf("   %s:%d\n", line.File, line.Line)
		fmt.Printf("   + %s\n\n", line.Content)
	}

	// Show removed lines
	fmt.Printf("❌ REMOVED LINES:\n")
	for _, line := range d.GetRemovedLines() {
		fmt.Printf("   %s\n", line.File)
		fmt.Printf("   - %s\n\n", line.Content)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("                  PARSE SUCCESSFUL! ✓")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
