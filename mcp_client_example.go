// mcp_client_example.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// This should be the path to the compiled server executable.
	// We'll take it from the command line for flexibility.
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run mcp_client_example.go <path_to_server_executable>")
	}
	serverExecutablePath := os.Args[1]

	ctx := context.Background()

	// 1. Create a transport that runs the server executable as a subprocess
	//    and connects to its stdin/stdout. The `-mcp` flag is passed as an argument.
	transport, err := mcp.NewCommandTransport(ctx, serverExecutablePath, "-mcp")
	if err != nil {
		log.Fatalf("Failed to create command transport: %v", err)
	}
	defer transport.Close()

	// 2. Create a new MCP client that communicates over this transport.
	client := mcp.NewClient(transport, nil) // nil for default ClientOptions
	defer client.Close()

	// 3. Define the input and output structs for the tool we want to call.
	//    Even if there are no inputs, we define an empty struct.
	type GetHeartbeatInput struct{}
	type GetHeartbeatOutput struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
	}

	// 4. Call the remote tool.
	log.Println("Calling tool 'getHeartbeat' on the server...")
	var result GetHeartbeatOutput
	_, err = client.CallTool(ctx, "getHeartbeat", GetHeartbeatInput{}, &result)
	if err != nil {
		log.Fatalf("Failed to call tool: %v", err)
	}

	// 5. Print the result!
	fmt.Println("========================================")
	fmt.Printf("Successfully received response from MCP server:
")
	fmt.Printf("  Status:    %s
", result.Status)
	fmt.Printf("  Timestamp: %s
", result.Timestamp)
	fmt.Println("========================================")
}
