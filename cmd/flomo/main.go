package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"mcp-server-flomo/pkg/flomo"
)

func main() {
	logger := log.New(os.Stdout, "[CLI] ", log.LstdFlags|log.Lmsgprefix)
	startTime := time.Now()

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Printf("Warning: Could not load .env file: %v", err)
		logger.Println("Will try to use environment variables directly")
	} else {
		logger.Println("Successfully loaded .env file")
	}

	// Get API URL from environment
	apiURL := os.Getenv("FLOMO_API_URL")
	if apiURL == "" {
		logger.Println("Error: FLOMO_API_URL environment variable is not set")
		fmt.Println("Please set it in your .env file or environment")
		os.Exit(1)
	}
	logger.Printf("Using Flomo API URL: %s", apiURL)

	// Define flags
	var (
		content string
		tags    string
		verbose bool
	)

	flag.StringVar(&content, "content", "", "Note content (required)")
	flag.StringVar(&content, "c", "", "Note content (shorthand)")
	flag.StringVar(&tags, "tags", "", "Comma-separated tags (optional)")
	flag.StringVar(&tags, "t", "", "Comma-separated tags (shorthand)")
	flag.BoolVar(&verbose, "verbose", false, "Show verbose output")
	flag.BoolVar(&verbose, "v", false, "Show verbose output (shorthand)")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -c \"This is a note\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -content \"This is a note\" -tags \"work,todo\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  echo \"This is a note\" | %s\n", os.Args[0])
	}

	flag.Parse()
	logger.Printf("Parsed command line flags (content length: %d, tags: %s)", len(content), tags)

	// If no content provided via flags, try to read from stdin
	if content == "" {
		logger.Println("No content provided via flags, checking stdin")
		// Check if there's input from pipe
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			logger.Println("Reading content from stdin")
			// Read from stdin
			var input strings.Builder
			buffer := make([]byte, 1024)
			for {
				n, err := os.Stdin.Read(buffer)
				if n > 0 {
					input.Write(buffer[:n])
				}
				if err != nil {
					break
				}
			}
			content = strings.TrimSpace(input.String())
			logger.Printf("Read %d characters from stdin", len(content))
		}
	}

	// Validate content
	if content == "" {
		logger.Println("Error: No content provided")
		fmt.Println("Error: Note content is required")
		flag.Usage()
		os.Exit(1)
	}

	// Process tags
	if tags != "" {
		logger.Printf("Processing tags: %s", tags)
		tagList := strings.Split(tags, ",")
		for _, tag := range tagList {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				content += " #" + tag
				logger.Printf("Added tag: %s", tag)
			}
		}
	}

	// Create Flomo client
	logger.Println("Creating Flomo client")
	client := flomo.NewClient(apiURL)

	// Send note
	logger.Println("Sending note to Flomo")
	resp, err := client.WriteNote(content)
	if err != nil {
		logger.Printf("Error sending note: %v", err)
		fmt.Printf("Error sending note: %v\n", err)
		os.Exit(1)
	}

	// Print success message with details
	duration := time.Since(startTime)
	logger.Printf("Note sent successfully (took %v)", duration)
	
	fmt.Println("\nNote sent successfully! 🎉")
	fmt.Printf("Created at: %s\n", resp.Memo.CreatedAt)
	if len(resp.Memo.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(resp.Memo.Tags, ", "))
	}
	fmt.Printf("View at: https://v.flomoapp.com/mine/?memo_id=%s\n", resp.Memo.Slug)

	if verbose {
		fmt.Printf("\nDetailed information:\n")
		fmt.Printf("- Source: %s\n", resp.Memo.Source)
		fmt.Printf("- Creator ID: %d\n", resp.Memo.CreatorID)
		fmt.Printf("- Response code: %d\n", resp.Code)
		fmt.Printf("- Response message: %s\n", resp.Message)
		fmt.Printf("- Total time: %v\n", duration)
	}
} 