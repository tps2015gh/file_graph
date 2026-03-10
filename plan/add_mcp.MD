# Plan: Integrating Model Context Protocol (MCP)

This document outlines the plan to integrate the Model Context Protocol (MCP) into the File Graph Visualizer application.

## 1. Objective

The primary goal is to expose the file graph's data and core functionalities to Large Language Models (LLMs) through a standardized MCP interface. This will enable programmatic interaction with the application's capabilities, such as querying graph data and triggering actions that can be reflected in the web UI.

## 2. Key Features & Tools

We will implement an MCP server that provides the following tools:

-   `getNodeCount()`: Returns the total number of nodes (files) in the most recent scan.
-   `getHeartbeat()`: Returns the server's operational status and timestamp, reusing the existing heartbeat logic.
-   `getNodeStatus(path: string)`: Retrieves detailed information for a specific node (file) given its path.
-   `searchAndHighlight(query: string)`: Searches for a node that matches the query and flags it to be highlighted in the web UI.

## 3. Architectural Changes

### 3.1. New MCP Package

-   A new package will be created at `internal/mcp` to encapsulate all MCP-related logic, including server setup, tool definitions, and handlers.

### 3.2. MCP SDK Integration

-   We will integrate the official Go SDK for MCP: `github.com/modelcontextprotocol/go-sdk`.
-   The dependency will be added to the `go.mod` file.

### 3.3. Transport Method

-   The initial implementation will use `StdioTransport` (Standard I/O). This is the most common transport for desktop LLM clients and local tool use, allowing the MCP server to run as a separate, callable process.

### 3.4. UI & MCP Communication

-   To allow the MCP `searchAndHighlight` tool to update the UI, a shared state mechanism is required.
    1.  **Shared State**: A global variable (e.g., `mcp_highlighted_node_id`) will be held in a shared package, accessible by both the MCP server and the web handlers. This variable will store the ID of the node to be highlighted. It must be protected by a mutex for safe concurrent access.
    2.  **New API Endpoint**: A new web API endpoint, `/api/highlighted-node`, will be created. The web frontend will poll this endpoint periodically.
    3.  **Frontend Logic**: The JavaScript on the frontend will be updated to:
        -   Fetch from `/api/highlighted-node` every few seconds.
        -   If a node ID is returned, it will use the existing graph logic to find and highlight that node.
        -   Once highlighted, it can clear the state by making a request or the backend can automatically clear it after being read once.

## 4. Implementation Steps

1.  **Dependency Management**: Add `github.com/modelcontextprotocol/go-sdk` to `go.mod`.
2.  **Create MCP Package**: Create the directory and initial files for `internal/mcp/mcp.go`.
3.  **Shared State**: Create a new package, e.g., `internal/state`, to hold the shared `highlighted_node_id` and its access mutex.
4.  **Implement MCP Server**:
    -   In `internal/mcp/mcp.go`, set up the `mcp.NewMCPServer`.
    -   Implement the `getNodeCount` tool by calling the existing `scanner` functions.
    -   Implement the `getHeartbeat` tool.
    -   Implement the `getNodeStatus` tool.
    -   Implement the `searchAndHighlight` tool, which will find a node and update the shared state variable.
5.  **Implement UI Communication Backend**:
    -   In `internal/handlers/handlers.go`, create the `HandleHighlightedNode` function.
    -   This handler will read the shared state variable and return its value as JSON.
6.  **Update Main Application**:
    -   In `main.go`, add a new command-line flag (e.g., `-mcp`) to launch the application in MCP server mode instead of web server mode. This keeps the concerns separate and follows the `StdioTransport` model.
7.  **Update Frontend**:
    -   In `index.html`, add a JavaScript `setInterval` function to poll the `/api/highlighted-node` endpoint.
    -   Add logic to process the response and call the existing graph visualization functions to highlight the specified node.
