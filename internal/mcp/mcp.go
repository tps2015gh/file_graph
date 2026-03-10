package mcp

import (
	"context"
	"file_graph/internal/models"
	"file_graph/internal/scanner"
	"file_graph/internal/state"
	"log"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RunServer initializes and starts the MCP server.
func RunServer(ctx context.Context) error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "FileGraphServer",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "getNodeCount",
		Description: "Gets the total number of nodes from the last scan",
	}, GetNodeCount)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "getHeartbeat",
		Description: "Checks the operational status of the server",
	}, GetHeartbeat)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "getNodeStatus",
		Description: "Retrieves details for a specific file path",
	}, GetNodeStatus)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "searchAndHighlight",
		Description: "Searches for a node and highlights it in the UI",
	}, SearchAndHighlight)

	log.Println("Starting MCP server on stdin/stdout...")
	return server.Run(ctx, &mcp.StdioTransport{})
}

// --- Tool Handlers ---

// GetNodeCount Tool
type GetNodeCountInput struct{} // Empty input struct
type GetNodeCountOutput struct {
	Count int `json:"count"`
}

func GetNodeCount(ctx context.Context, req *mcp.CallToolRequest, input GetNodeCountInput) (*mcp.CallToolResult, GetNodeCountOutput, error) {
	nodes, _ := scanner.GetCachedNodes()
	count := len(nodes)
	return &mcp.CallToolResult{}, GetNodeCountOutput{Count: count}, nil
}

// GetHeartbeat Tool
type GetHeartbeatInput struct{} // Empty input struct
type GetHeartbeatOutput struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func GetHeartbeat(ctx context.Context, req *mcp.CallToolRequest, input GetHeartbeatInput) (*mcp.CallToolResult, GetHeartbeatOutput, error) {
	now := time.Now().Format("15:04:05")
	return &mcp.CallToolResult{}, GetHeartbeatOutput{Status: "alive", Timestamp: now}, nil
}

// GetNodeStatus Tool
type GetNodeStatusInput struct {
	Path string `json:"path" jsonschema:"the absolute path of the file to query"`
}
type GetNodeStatusOutput struct {
	Node  *models.FileNode `json:"node"`
	Found bool             `json:"found"`
}

func GetNodeStatus(ctx context.Context, req *mcp.CallToolRequest, input GetNodeStatusInput) (*mcp.CallToolResult, GetNodeStatusOutput, error) {
	nodes, _ := scanner.GetCachedNodes()
	for _, node := range nodes {
		if node.Path == input.Path {
			nodeCopy := node
			return &mcp.CallToolResult{}, GetNodeStatusOutput{Node: &nodeCopy, Found: true}, nil
		}
	}
	return &mcp.CallToolResult{}, GetNodeStatusOutput{Node: nil, Found: false}, nil
}

// SearchAndHighlight Tool
type SearchAndHighlightInput struct {
	Query string `json:"query" jsonschema:"the search query to find a matching node"`
}
type SearchAndHighlightOutput struct {
	Found   bool   `json:"found"`
	NodeID  string `json:"node_id,omitempty"`
	Message string `json:"message"`
}

func SearchAndHighlight(ctx context.Context, req *mcp.CallToolRequest, input SearchAndHighlightInput) (*mcp.CallToolResult, SearchAndHighlightOutput, error) {
	nodes, _ := scanner.GetCachedNodes()
	for _, node := range nodes {
		if strings.Contains(strings.ToLower(node.Path), strings.ToLower(input.Query)) {
			state.GlobalState.SetHighlightedNodeID(node.ID)
			return &mcp.CallToolResult{}, SearchAndHighlightOutput{Found: true, NodeID: node.ID, Message: "Node found and set for highlighting."}, nil
		}
	}
	return &mcp.CallToolResult{}, SearchAndHighlightOutput{Found: false, Message: "No node found matching the query."}, nil
}
