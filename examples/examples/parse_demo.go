package main

import (
	"fmt"
	"log"

	"github.com/BhardwajShrey/zizou/internal/diff"
)

func main() {
	// Your example diff
	diffContent := `diff --git a/app/main.go b/app/main.go
index 30dcdee..2cb2e68 100644
--- a/app/main.go
+++ b/app/main.go
@@ -23,7 +23,8 @@ func main() {

     image := os.Args[2]

-    EnterNewJail(os.Args[3], image)
+    jailPath := EnterNewJail(os.Args[3], image)
+    defer os.Remove(jailPath)

     cmd := exec.Command(command, args...)
     cmd.Stderr = os.Stderr
@@ -37,14 +38,12 @@ func main() {
     os.Exit(cmd.ProcessState.ExitCode())
 }

-func EnterNewJail(filepath string, image string) {
+func EnterNewJail(filepath string, image string) string {
     tempDirPath, err := os.MkdirTemp("", "temp_folder_*")
     if err != nil {
         throwerror.ThrowError(err, "Unable to create temp directory")
     }

-    defer os.Remove(tempDirPath)
-
     err = os.Chmod(tempDirPath, 0777)
     if err != nil {
         throwerror.ThrowError(err, "Error modifying rwx on tempDirPath")
@@ -58,4 +57,6 @@ func EnterNewJail(filepath string, image string) {
     if err != nil {
         throwerror.ThrowError(err, fmt.Sprintf("Error in executing chroot on %s", tempDirPath))
     }
+
+    return tempDirPath
 }`

	// Parse the diff
	parser := diff.NewParser()
	parsedDiff, err := parser.Parse(diffContent)
	if err != nil {
		log.Fatalf("Failed to parse diff: %v", err)
	}

	// Extract information
	fmt.Println("=== DIFF ANALYSIS ===\n")

	// 1. Changed files
	fmt.Printf("Changed Files: %d\n", len(parsedDiff.Files))
	for _, file := range parsedDiff.Files {
		fmt.Printf("  - %s\n", file.NewPath)
	}
	fmt.Println()

	// 2. Detailed analysis per file
	for _, file := range parsedDiff.Files {
		fmt.Printf("File: %s\n", file.NewPath)
		fmt.Printf("Hunks: %d\n\n", len(file.Hunks))

		for hunkIdx, hunk := range file.Hunks {
			fmt.Printf("  Hunk #%d: @@ -%d,%d +%d,%d @@\n",
				hunkIdx+1, hunk.OldStart, hunk.OldCount, hunk.NewStart, hunk.NewCount)

			// Count changes
			added := 0
			removed := 0
			context := 0

			for _, line := range hunk.Lines {
				switch line.Type {
				case diff.LineAdded:
					added++
				case diff.LineRemoved:
					removed++
				case diff.LineContext:
					context++
				}
			}

			fmt.Printf("    Added: %d, Removed: %d, Context: %d\n\n", added, removed, context)

			// Show actual changes
			fmt.Println("    Changes:")
			for _, line := range hunk.Lines {
				prefix := " "
				switch line.Type {
				case diff.LineAdded:
					prefix = "+"
				case diff.LineRemoved:
					prefix = "-"
				}
				fmt.Printf("    %s %s\n", prefix, line.Content)
			}
			fmt.Println()
		}
	}

	// 3. Extract only added lines
	fmt.Println("=== ADDED LINES ONLY ===")
	for _, file := range parsedDiff.Files {
		for _, hunk := range file.Hunks {
			for _, line := range hunk.Lines {
				if line.Type == diff.LineAdded {
					fmt.Printf("%s:%d: %s\n", file.NewPath, line.Number, line.Content)
				}
			}
		}
	}
	fmt.Println()

	// 4. Extract only removed lines
	fmt.Println("=== REMOVED LINES ONLY ===")
	for _, file := range parsedDiff.Files {
		for _, hunk := range file.Hunks {
			for _, line := range hunk.Lines {
				if line.Type == diff.LineRemoved {
					fmt.Printf("%s: %s\n", file.NewPath, line.Content)
				}
			}
		}
	}
	fmt.Println()

	// 5. Show context around changes
	fmt.Println("=== CHANGES WITH CONTEXT ===")
	for _, file := range parsedDiff.Files {
		for _, hunk := range file.Hunks {
			for i, line := range hunk.Lines {
				if line.Type == diff.LineAdded || line.Type == diff.LineRemoved {
					// Show 2 lines before and after
					start := max(0, i-2)
					end := min(len(hunk.Lines), i+3)

					fmt.Printf("Change in %s around line %d:\n", file.NewPath, line.Number)
					for j := start; j < end; j++ {
						marker := " "
						if j == i {
							marker = "→"
						}
						prefix := " "
						switch hunk.Lines[j].Type {
						case diff.LineAdded:
							prefix = "+"
						case diff.LineRemoved:
							prefix = "-"
						}
						fmt.Printf("%s %s %s\n", marker, prefix, hunk.Lines[j].Content)
					}
					fmt.Println()
				}
			}
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
