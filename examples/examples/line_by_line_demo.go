package main

import (
	"fmt"
	"strings"
)

// This demo shows how to read a diff line-by-line
func main() {
	diff := `diff --git a/app/main.go b/app/main.go
index 30dcdee..2cb2e68 100644
--- a/app/main.go
+++ b/app/main.go
@@ -23,7 +23,8 @@ func main() {

     image := os.Args[2]

-    EnterNewJail(os.Args[3], image)
+    jailPath := EnterNewJail(os.Args[3], image)
+    defer os.Remove(jailPath)

     cmd := exec.Command(command, args...)`

	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("     GIT DIFF LINE-BY-LINE EXPLANATION")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println()

	lines := strings.Split(diff, "\n")

	for i, line := range lines {
		lineNum := i + 1
		fmt.Printf("Line %2d: ", lineNum)

		// Show the line with visible spacing
		displayLine := strings.Replace(line, "\t", "→", -1)
		if line == "" {
			fmt.Print("(empty line)")
		} else if strings.HasPrefix(line, " ") && len(line) > 1 {
			fmt.Printf("«%s»", displayLine) // « » show it's a context line with space
		} else {
			fmt.Print(displayLine)
		}
		fmt.Println()

		// Explain what this line means
		fmt.Print("         ")
		explainLine(line, lineNum)
		fmt.Println()
		fmt.Println()
	}

	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("                 SUMMARY")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("Prefixes:")
	fmt.Println("  (none)   = Header/metadata line")
	fmt.Println("  « »      = Context (space prefix, unchanged)")
	fmt.Println("  -        = Removed line (deleted)")
	fmt.Println("  +        = Added line (new code)")
	fmt.Println()
	fmt.Println("This change:")
	fmt.Println("  • Modified: app/main.go")
	fmt.Println("  • Removed:  1 line")
	fmt.Println("  • Added:    2 lines")
	fmt.Println("  • Net:      +1 line")
}

func explainLine(line string, lineNum int) {
	switch {
	case strings.HasPrefix(line, "diff --git"):
		parts := strings.Fields(line)
		if len(parts) >= 4 {
			oldFile := strings.TrimPrefix(parts[2], "a/")
			newFile := strings.TrimPrefix(parts[3], "b/")
			fmt.Printf("📁 FILE HEADER: Comparing '%s' (old) with '%s' (new)", oldFile, newFile)
		}

	case strings.HasPrefix(line, "index "):
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			hashes := strings.Split(parts[1], "..")
			mode := parts[2]
			fmt.Printf("🔖 GIT METADATA: Hash %s → %s, Mode: %s", hashes[0][:7], hashes[1][:7], mode)
		}

	case strings.HasPrefix(line, "---"):
		file := strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(line, "---")), "a/")
		fmt.Printf("📄 OLD FILE MARKER: Original version is '%s'", file)

	case strings.HasPrefix(line, "+++"):
		file := strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(line, "+++")), "b/")
		fmt.Printf("📄 NEW FILE MARKER: Modified version is '%s'", file)

	case strings.HasPrefix(line, "@@"):
		fmt.Print("📍 HUNK HEADER: ")
		explainHunkHeader(line)

	case strings.HasPrefix(line, "-"):
		content := strings.TrimPrefix(line, "-")
		fmt.Printf("❌ REMOVED: '%s'", content)

	case strings.HasPrefix(line, "+"):
		content := strings.TrimPrefix(line, "+")
		fmt.Printf("✅ ADDED: '%s'", content)

	case strings.HasPrefix(line, " ") || line == "":
		if line == "" {
			fmt.Print("⚪ CONTEXT: Empty line (unchanged)")
		} else {
			content := strings.TrimPrefix(line, " ")
			fmt.Printf("⚪ CONTEXT: '%s' (unchanged)", content)
		}

	default:
		fmt.Print("❓ UNKNOWN: Not standard diff format")
	}
}

func explainHunkHeader(line string) {
	// Parse @@ -23,7 +23,8 @@ func main() {
	// Find the content between @@
	parts := strings.Split(line, "@@")
	if len(parts) < 2 {
		fmt.Print("Malformed hunk header")
		return
	}

	ranges := strings.TrimSpace(parts[1])
	context := ""
	if len(parts) > 2 {
		context = strings.TrimSpace(parts[2])
	}

	// Split into old (-) and new (+) ranges
	rangeParts := strings.Fields(ranges)
	if len(rangeParts) >= 2 {
		oldRange := strings.TrimPrefix(rangeParts[0], "-")
		newRange := strings.TrimPrefix(rangeParts[1], "+")

		oldParts := strings.Split(oldRange, ",")
		newParts := strings.Split(newRange, ",")

		oldStart := oldParts[0]
		oldCount := "1"
		if len(oldParts) > 1 {
			oldCount = oldParts[1]
		}

		newStart := newParts[0]
		newCount := "1"
		if len(newParts) > 1 {
			newCount = newParts[1]
		}

		fmt.Printf("Old file: line %s, %s lines shown | New file: line %s, %s lines shown",
			oldStart, oldCount, newStart, newCount)

		if context != "" {
			fmt.Printf(" | Context: %s", context)
		}
	}
}
