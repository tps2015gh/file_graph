# Model Context Protocol (MCP) Usage for File Graph Visualizer

This document explains how to interact with the File Graph Visualizer's MCP server, detailing the available tools, their usage, and example scenarios.

## 1. Running the MCP Server

The File Graph Visualizer can be launched in MCP server mode by providing the `-mcp` flag. This will make it listen for MCP commands on standard input/output (stdio).

```bash
go run main.go -mcp
```

Or if you have a compiled executable:

```bash
./file_graph_server.exe -mcp
```

## 2. Interacting with the MCP Server

You will need an MCP client to interact with the server. Many LLM interfaces have built-in MCP client capabilities. For programmatic interaction, you can use an MCP client library in your preferred language or a command-line tool like `mcp-cli` (if available).

The MCP server expects and returns JSON objects over stdio according to the MCP specification.

## 3. Available MCP Tools

The following tools are exposed by the File Graph Visualizer MCP server:

### 3.1. `getNodeCount`

**Description**: Returns the total number of file nodes discovered in the last scan.
**Input**: None
**Output**:
```json
{
  "count": 123
}
```
**Example Scenario**: An LLM might use this to quickly get an overview of the dataset size before performing more detailed queries.

### 3.2. `getHeartbeat`

**Description**: Checks the operational status of the server.
**Input**: None
**Output**:
```json
{
  "status": "alive",
  "timestamp": "15:04:05"
}
```
**Example Scenario**: An LLM or an orchestrator could use this to verify that the File Graph Visualizer server is running and responsive.

### 3.3. `getNodeStatus`

**Description**: Retrieves detailed information for a specific file node given its absolute path.
**Input**:
```json
{
  "path": "C:\dev\file_graph\main.go"
}
```
**Output (if found)**:
```json
{
  "found": true,
  "node": {
    "id": "C:\dev\file_graph\main.go",
    "name": "main.go",
    "path": "C:\dev\file_graph\main.go",
    "is_folder": false,
    "size": 5120,
    "modified_at": "2026-03-11T10:30:00Z",
    "created_at": "2026-03-11T10:30:00Z",
    "size_last_3": 120,
    "name_last_4": "main",
    "hash": "a1b2c3d4e5f6...",
    "vector": [0.1, 0.2, 0.3, ...]
  }
}
```
**Output (if not found)**:
```json
{
  "found": false,
  "node": null
}
```
**Example Scenario**: An LLM, after being asked about a specific file, could use this tool to retrieve its metadata (size, modification date, hash) to provide a detailed answer to the user.

### 3.4. `searchAndHighlight`

**Description**: Searches for a file node whose path or name contains the provided query string. If a node is found, its ID is sent to the web frontend for highlighting in the visualizer.
**Input**:
```json
{
  "query": "index.html"
}
```
**Output (if found)**:
```json
{
  "found": true,
  "node_id": "C:\dev\file_graph\index.html",
  "message": "Node found and set for highlighting."
}
```
**Output (if not found)**:
```json
{
  "found": false,
  "message": "No node found matching the query."
}
```
**Example Scenario**: An LLM could use this tool to help a user visually locate a file or related files in the graph. For instance, if a user asks "Show me where the main HTML file is," the LLM could call `searchAndHighlight` with "index.html" and the node would be highlighted in the browser.

## 4. UI Reflection

When the `searchAndHighlight` tool is successfully invoked, the `File Graph Visualizer`'s web interface (running in parallel) will periodically poll the `/api/highlighted-node` endpoint. If a `highlightNodeId` is returned, the UI will automatically select that node, display its details, and pan the view to center on it. This provides a direct visual feedback loop from LLM interaction to the user interface.

## 6. How AI Interprets MCP Tools ("Manual for AI")

For an AI (specifically, a Large Language Model or LLM) to understand and effectively use these MCP tools, it relies on two primary pieces of information provided directly by the MCP server:

1.  **Natural Language Descriptions**: Each tool is registered with a clear, concise human-readable description (e.g., "Gets the total number of nodes from the last scan"). The LLM uses its natural language understanding to match a user's intent to the purpose of the tool.

2.  **Structured Schema (Input/Output Definitions)**: The MCP server also provides a formal definition (schema, often JSON Schema) for each tool's input parameters and its expected output. This schema, derived from the Go struct fields and `json` / `jsonschema` tags in our code, tells the LLM:
    *   What arguments the tool expects (e.g., `path: string`).
    *   What data types those arguments should be.
    *   What the structure of the returned data will be.

Therefore, the "manual" that an AI reads is not a separate document, but rather this rich, structured metadata that the MCP server dynamically provides about each tool's capabilities and interface. The LLM processes this information to decide which tool to call, how to construct the input, and how to interpret the output.
