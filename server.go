package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
    "mcp-server-flomo/pkg/flomo"
)

func main() {
    logger := log.New(os.Stdout, "[Server] ", log.LstdFlags|log.Lmsgprefix)

    // Load .env file if it exists
    if err := godotenv.Load(); err != nil {
        logger.Printf("Warning: Could not load .env file: %v", err)
        logger.Println("Will try to use environment variables directly")
    } else {
        logger.Println("Successfully loaded .env file")
    }

    // Get Flomo API URL from environment variable
    flomoAPIURL := os.Getenv("FLOMO_API_URL")
    if flomoAPIURL == "" {
        logger.Println("Error: FLOMO_API_URL environment variable is not set")
        fmt.Println("Please set it in your .env file or environment variables")
        os.Exit(1)
    }
    logger.Printf("Using Flomo API URL: %s", flomoAPIURL)

    // Create Flomo client
    flomoClient := flomo.NewClient(flomoAPIURL)

    // Create a new MCP server
    s := server.NewMCPServer(
        "mcp-server-flomo",
        "1.0.0",
        server.WithResourceCapabilities(true, true),
        server.WithLogging(),
        server.WithRecovery(),
    )
    logger.Println("Created MCP server instance")

    // Add write_note tool
    writeNoteTool := mcp.NewTool("write_note",
        mcp.WithDescription("Write note to flomo"),
        mcp.WithString("content",
            mcp.Required(),
            mcp.Description("Text content of the note with markdown format"),
        ),
    )
    logger.Println("Created write_note tool definition")

    // Add the write_note handler
    s.AddTool(writeNoteTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        content := request.Params.Arguments["content"].(string)
        logger.Printf("Received write_note request with content length: %d", len(content))

        // Write note to Flomo
        resp, err := flomoClient.WriteNote(content)
        if err != nil {
            logger.Printf("Error writing note: %v", err)
            return nil, fmt.Errorf("failed to write note: %v", err)
        }

        // Create memo URL
        memoURL := fmt.Sprintf("https://v.flomoapp.com/mine/?memo_id=%s", resp.Memo.Slug)
        resultText := fmt.Sprintf("Successfully wrote note to Flomo.\nMemo URL: %s\nCreated at: %s", 
            memoURL, resp.Memo.CreatedAt)

        logger.Printf("Successfully wrote note. Memo URL: %s", memoURL)
        return mcp.NewToolResultText(resultText), nil
    })
    logger.Println("Registered write_note tool handler")

    // Start the server
    logger.Println("Starting MCP server...")
    if err := server.ServeStdio(s); err != nil {
        logger.Printf("Server error: %v", err)
        os.Exit(1)
    }
} 