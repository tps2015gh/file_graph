package state

import "sync"

// SharedState holds the data that needs to be accessible across different parts of the application,
// such as between the MCP server and the web handlers.
type SharedState struct {
	mu                 sync.RWMutex
	highlightedNodeID  string
}

// GlobalState is the singleton instance of SharedState.
var GlobalState = &SharedState{}

// SetHighlightedNodeID safely sets the ID of the node to be highlighted in the UI.
func (s *SharedState) SetHighlightedNodeID(nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.highlightedNodeID = nodeID
}

// GetHighlightedNodeID safely retrieves the ID of the node to be highlighted.
// It returns the ID and a boolean indicating if an ID was present.
// It also clears the ID after reading to prevent it from being processed multiple times.
func (s *SharedState) GetHighlightedNodeID() (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.highlightedNodeID == "" {
		return "", false
	}
	nodeID := s.highlightedNodeID
	// Clear the ID after reading
	s.highlightedNodeID = ""
	return nodeID, true
}
