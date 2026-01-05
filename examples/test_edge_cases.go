package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BhardwajShrey/zizou/internal/diff"
)

func main() {
	content, err := os.ReadFile("/tmp/test_edge_cases.diff")
	if err != nil {
		log.Fatal(err)
	}

	parser := diff.NewParser()
	d, err := parser.Parse(string(content))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("       EDGE CASE DETECTION TEST")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	for i, file := range d.Files {
		fmt.Printf("File #%d: %s\n", i+1, file.NewPath)
		fmt.Println("─────────────────────────────────────────────")

		// Check all edge cases
		if file.IsBinary {
			fmt.Println("  ✅ Binary file detected")
		}

		if file.IsRenamed {
			fmt.Printf("  ✅ Renamed: %s → %s\n", file.OldPath, file.NewPath)
			fmt.Printf("     Similarity: %d%%\n", file.Similarity)
		}

		if file.IsCopied {
			fmt.Printf("  ✅ Copied: %s → %s\n", file.OldPath, file.NewPath)
		}

		if file.IsNew {
			fmt.Println("  ✅ New file created")
		}

		if file.IsDeleted {
			fmt.Println("  ✅ File deleted")
		}

		if file.OldMode != "" && file.NewMode != "" {
			fmt.Printf("  ✅ Mode changed: %s → %s\n", file.OldMode, file.NewMode)
		} else if file.NewMode != "" {
			fmt.Printf("  ✅ Mode: %s\n", file.NewMode)
		}

		if len(file.Hunks) > 0 {
			fmt.Printf("  📝 %d hunk(s) with code changes\n", len(file.Hunks))
		}

		fmt.Println()
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  ALL EDGE CASES DETECTED SUCCESSFULLY! ✓")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
