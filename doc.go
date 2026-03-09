// Package main provides a standalone web application for file graph visualization.
// It scans directories, extracts file metadata, calculates similarity vectors,
// and visualizes these relations in a 2D interactive graph on the web browser.
//
// Features include:
// - Recursive file scanning with metadata extraction (name, size, hash, etc.).
// - Cosine similarity calculation between file attributes.
// - Force-directed 2D layout for relationship visualization.
// - Search and details panel for investigating file connections.
package main
