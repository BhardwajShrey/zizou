package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/shreybhardwaj/zizou/internal/cache"
	"github.com/shreybhardwaj/zizou/internal/client"
	"github.com/shreybhardwaj/zizou/internal/diff"
	"github.com/shreybhardwaj/zizou/internal/output"
	"github.com/shreybhardwaj/zizou/internal/review"
	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFmt  string
	apiKey     string
	cacheDir   string
	noCache    bool
)

var rootCmd = &cobra.Command{
	Use:   "zizou",
	Short: "AI-powered code review CLI tool",
	Long: `Zizou is a Go-based code review CLI tool that:
- Reads git diffs from stdin or a file
- Sends the diff to Claude API for analysis
- Returns structured review comments with severity levels
- Caches duplicate results to avoid API calls`,
	RunE: runReview,
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "file", "f", "", "Input file containing git diff (defaults to stdin)")
	rootCmd.Flags().StringVarP(&outputFmt, "output", "o", "text", "Output format (text, json, markdown)")
	rootCmd.Flags().StringVar(&apiKey, "api-key", "", "Claude API key (or set ANTHROPIC_API_KEY env var)")
	rootCmd.Flags().StringVar(&cacheDir, "cache-dir", "", "Cache directory (defaults to ~/.zizou/cache)")
	rootCmd.Flags().BoolVar(&noCache, "no-cache", false, "Disable caching")
}

func Execute() error {
	return rootCmd.Execute()
}

func runReview(cmd *cobra.Command, args []string) error {
	// Read diff input
	var diffContent string
	var err error

	if inputFile != "" {
		diffContent, err = readFromFile(inputFile)
	} else {
		diffContent, err = readFromStdin()
	}
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Parse diff
	parser := diff.NewParser()
	parsedDiff, err := parser.Parse(diffContent)
	if err != nil {
		return fmt.Errorf("failed to parse diff: %w", err)
	}

	// Initialize cache
	var cacheStore cache.Cache
	if !noCache {
		if cacheDir == "" {
			homeDir, _ := os.UserHomeDir()
			cacheDir = fmt.Sprintf("%s/.zizou/cache", homeDir)
		}
		cacheStore, err = cache.NewFileCache(cacheDir)
		if err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}
	} else {
		cacheStore = cache.NewNoOpCache()
	}

	// Get API key
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if apiKey == "" {
		return fmt.Errorf("API key required: set --api-key flag or ANTHROPIC_API_KEY environment variable")
	}

	// Initialize Claude client
	claudeClient := client.NewClaudeClient(apiKey)

	// Create reviewer
	reviewer := review.NewReviewer(claudeClient, cacheStore)

	// Perform review
	comments, err := reviewer.Review(cmd.Context(), parsedDiff)
	if err != nil {
		return fmt.Errorf("review failed: %w", err)
	}

	// Format and output results
	formatter := output.NewFormatter(outputFmt)
	result, err := formatter.Format(comments)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	fmt.Println(result)
	return nil
}

func readFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readFromStdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
