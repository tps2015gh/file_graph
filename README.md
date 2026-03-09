# File Graph Visualizer

A standalone Go web application to scan directories, extract file metadata, calculate similarity embeddings, and visualize relationships as a 2D interactive graph.

## Features

- **File Scanner**: Recursively scans folders for metadata.
- **Metadata Extraction**:
  - Filename and last 4 characters (excluding extension).
  - Size and last 3 digits of size.
  - Creation and modification timestamps.
  - SHA-256 hash value.
- **Similarity Calculation**: Normalizes metadata into a vector to calculate relationship proximity.
- **Interactive Visualization**: 2D Canvas graph with "gravity" effects, resembling a colony of stars.
- **Search**: Interactive search to find and highlight related files.

## Tech Stack

- **Backend**: Go (standard library)
- **Frontend**: HTML5, Vanilla JavaScript, CSS, Canvas API

## Getting Started

1.  Clone the repository.
2.  Run `go run main.go`.
3.  Open `http://localhost:8080` in your browser.

## Authors

- **Main Developer**: The User
- **Assistant**: Gemini CLI
