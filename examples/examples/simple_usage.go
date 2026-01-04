package main

import (
	"fmt"
	"log"

	"github.com/BhardwajShrey/zizou/internal/diff"
)

func main() {
	// Simple example showing the most common use cases

	diffContent := `diff --git a/app/main.go b/app/main.go
index 30dcdee..2cb2e68 100644
--- a/app/main.go
+++ b/app/main.go
@@ -23,7 +23,8 @@ func main() {

     image := os.Args[2]

-    EnterNewJail(os.Args[3], image)
+    jailPath := EnterNewJail(os.Args[3], image)
+    defer os.Remove(jailPath)

     cmd := exec.Command(command, args...)`

	// Parse
	parser := diff.NewParser()
	d, err := parser.Parse(diffContent)
	if err != nil {
		log.Fatal(err)
	}

	// Example 1: Get basic stats
	fmt.Println("=== Statistics ===")
	stats := d.Stats()
	fmt.Printf("Files changed: %d\n", stats.Files)
	fmt.Printf("Lines added: %d\n", stats.LinesAdded)
	fmt.Printf("Lines removed: %d\n", stats.LinesRemoved)
	fmt.Println()

	// Example 2: List modified files
	fmt.Println("=== Modified Files ===")
	for _, file := range d.GetModifiedFiles() {
		fmt.Printf("  • %s\n", file)
	}
	fmt.Println()

	// Example 3: Show only what was added
	fmt.Println("=== New Code ===")
	for _, line := range d.GetAddedLines() {
		fmt.Printf("%s:%d → %s\n", line.File, line.Line, line.Content)
	}
	fmt.Println()

	// Example 4: Show only what was removed
	fmt.Println("=== Removed Code ===")
	for _, line := range d.GetRemovedLines() {
		fmt.Printf("%s → %s\n", line.File, line.Content)
	}
	fmt.Println()

	// Example 5: Detailed file analysis
	fmt.Println("=== Detailed Analysis ===")
	for _, file := range d.Files {
		fmt.Printf("\nFile: %s\n", file.NewPath)

		for hunkIdx, hunk := range file.Hunks {
			if hunk.HasChanges() {
				fmt.Printf("  Change section #%d (lines %d-%d):\n",
					hunkIdx+1, hunk.NewStart, hunk.NewStart+hunk.NewCount-1)

				for _, line := range hunk.Lines {
					if line.Type != diff.LineContext {
						prefix := "+"
						if line.Type == diff.LineRemoved {
							prefix = "-"
						}
						fmt.Printf("    %s %s\n", prefix, line.Content)
					}
				}
			}
		}
	}
}
