# Architecture & Implementation Details

This document explains how the File Graph Visualizer works.

## Backend (Go)

- **Scanner**: The `scanDirectory` function uses `filepath.WalkDir` to traverse the directory tree.
- **Metadata Extraction**:
    - `FileNode` stores attributes like name, size, last 3 digits of size, last 4 chars of filename (without extension), and SHA-256 hash.
    - `calculateHash` reads the file content to compute its SHA-256 checksum.
- **Embedding / Similarity Logic**:
    - `buildVector` normalizes various attributes (size, name hash, timestamps) into a 4D vector.
    - `calculateRelations` performs an N² comparison between nodes using `cosineSimilarity` on their vectors, supplemented by explicit attribute matches (e.g., matching the last 4 characters of a filename or the last 3 digits of the size).
- **API**:
    - `/api/scan`: Triggers the scanning and relation calculation process and returns JSON nodes and links.

## Frontend (JavaScript/Canvas)

- **Force-Directed Graph**:
    - Uses a custom physics simulation implemented on a 2D HTML5 Canvas.
    - **Repulsion**: Nodes push each other away to prevent overlaps.
    - **Attraction (Links)**: Similar nodes attract each other, forming "colonies" or clusters.
    - **Gravity**: Nodes are pulled toward the center of the screen to keep them within view.
- **Interactive UI**:
    - Clicking a node highlights it and displays its metadata in the side panel.
    - The "Related Files" table lists nodes that have a high similarity score with the selected node.
    - Search: Entering a keyword in the search bar finds the corresponding node and displays its related files.

## Metadata Mapping

The user's requested attributes are mapped as follows:
- `name`: `node.name`
- `size`: `node.size`
- `last 3digit of size`: `node.size % 1000`
- `last 4 char of file name (no extension)`: `strings.TrimSuffix(name, ext)[-4:]`
- `create date time`: Currently uses `ModTime` as a fallback.
- `modify date time`: `info.ModTime()`
- `hash value`: `sha256(file_content)`
