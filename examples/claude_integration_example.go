package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BhardwajShrey/zizou/internal/client"
	"github.com/BhardwajShrey/zizou/internal/diff"
)

func main() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("   CLAUDE API INTEGRATION EXAMPLE")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Example 1: Basic Configuration
	fmt.Println("1️⃣  CREATING CONFIGURATION")
	fmt.Println("────────────────────────────────────────────")

	// Load config from environment
	config, err := client.NewConfigFromEnv()
	if err != nil {
		log.Printf("⚠️  Config from env failed: %v", err)
		log.Println("   Using manual configuration instead...")

		// Manual configuration (for demo purposes)
		config = client.DefaultConfig()
		config.APIKey = os.Getenv("ANTHROPIC_API_KEY_ZIZOU")

		if config.APIKey == "" {
			log.Fatal("❌ ANTHROPIC_API_KEY_ZIZOU environment variable not set")
		}
	}

	fmt.Printf("   ✅ Model: %s\n", config.Model)
	fmt.Printf("   ✅ Max Tokens: %d\n", config.MaxTokens)
	fmt.Printf("   ✅ Timeout: %v\n", config.Timeout)
	fmt.Printf("   ✅ Max Retries: %d\n", config.MaxRetries)
	fmt.Printf("   ✅ Rate Limit: %d req/min\n", config.RateLimit)
	fmt.Println()

	// Example 2: Create Enhanced Client
	fmt.Println("2️⃣  CREATING ENHANCED CLIENT")
	fmt.Println("────────────────────────────────────────────")

	enhancedClient, err := client.NewEnhancedClient(config)
	if err != nil {
		log.Fatalf("❌ Failed to create client: %v", err)
	}

	fmt.Println("   ✅ Enhanced client created")
	fmt.Println("   • Rate limiting: enabled")
	fmt.Println("   • Retry logic: enabled")
	fmt.Println("   • Error handling: enhanced")
	fmt.Println()

	// Example 3: Create Reviewer Client
	fmt.Println("3️⃣  CREATING REVIEWER CLIENT")
	fmt.Println("────────────────────────────────────────────")

	reviewerClient, err := client.NewReviewerClient(config)
	if err != nil {
		log.Fatalf("❌ Failed to create reviewer: %v", err)
	}

	fmt.Println("   ✅ Reviewer client created")
	fmt.Println()

	// Example 4: Parse a Test Diff
	fmt.Println("4️⃣  PARSING TEST DIFF")
	fmt.Println("────────────────────────────────────────────")

	testDiff := `diff --git a/app/main.go b/app/main.go
index 30dcdee..2cb2e68 100644
--- a/app/main.go
+++ b/app/main.go
@@ -23,7 +23,8 @@ func main() {

     image := os.Args[2]

-    EnterNewJail(os.Args[3], image)
+    jailPath := EnterNewJail(os.Args[3], image)
+    defer os.Remove(jailPath)

     cmd := exec.Command(command, args...)
     cmd.Stderr = os.Stderr`

	parser := diff.NewParser()
	parsedDiff, err := parser.Parse(testDiff)
	if err != nil {
		log.Fatalf("❌ Failed to parse diff: %v", err)
	}

	stats := parsedDiff.Stats()
	fmt.Printf("   ✅ Parsed diff successfully\n")
	fmt.Printf("   • Files changed: %d\n", stats.Files)
	fmt.Printf("   • Lines added: +%d\n", stats.LinesAdded)
	fmt.Printf("   • Lines removed: -%d\n", stats.LinesRemoved)
	fmt.Println()

	// Example 5: Send for Review (if API key is set)
	fmt.Println("5️⃣  SENDING FOR CLAUDE REVIEW")
	fmt.Println("────────────────────────────────────────────")

	if config.APIKey == "" || config.APIKey == "your-api-key-here" {
		fmt.Println("   ⚠️  Skipping API call (no valid API key)")
		fmt.Println("   Set ANTHROPIC_API_KEY_ZIZOU to test live review")
	} else {
		fmt.Println("   🚀 Sending diff to Claude...")

		ctx := context.Background()
		result, err := reviewerClient.ReviewDiff(ctx, parsedDiff)

		if err != nil {
			log.Printf("   ❌ Review failed: %v\n", err)
		} else {
			fmt.Println("   ✅ Review completed!")
			fmt.Println()

			// Display results
			fmt.Println("6️⃣  REVIEW RESULTS")
			fmt.Println("────────────────────────────────────────────")

			if len(result.Comments) == 0 {
				fmt.Println("   ✅ No issues found!")
			} else {
				fmt.Printf("   Found %d comment(s):\n\n", len(result.Comments))

				for i, comment := range result.Comments {
					fmt.Printf("   [%d] %s:%d\n", i+1, comment.File, comment.Line)
					fmt.Printf("       Severity: %s | Category: %s\n", comment.Severity, comment.Category)
					fmt.Printf("       %s\n\n", comment.Message)
				}
			}

			if result.Summary != "" {
				fmt.Println("   Summary:")
				fmt.Printf("   %s\n", result.Summary)
			}
		}
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("   EXAMPLE COMPLETE!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Example 6: Direct API Call
	fmt.Println()
	fmt.Println("💡 TIP: You can also use the enhanced client directly:")
	fmt.Println()
	if enhancedClient != nil {
		fmt.Println("   response, err := enhancedClient.SendMessage(ctx, \"Your prompt here\")")
		fmt.Println("   if err != nil {")
		fmt.Println("       // Handles retries, rate limiting automatically")
		fmt.Println("   }")
	}
	fmt.Println()
}
