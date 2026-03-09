# File Graph Visualizer (Star Colony)

A standalone Go web application to scan directories, extract file metadata, calculate similarity embeddings (26D Vector Space), and visualize relationships as an interactive 2D interactive graph.

## Project Roles & Direction
- **Lead Programmer & Director**: The User
- **AI Team Assistance** (Last Updated: 2026-03-09):
  - Gemini CLI (Powered by multiple models)
  - OpenCode (Advanced AI coding assistant)
  - Big Pickle (Specialized AI model)

## Core Concept
This program treats the filesystem as a high-dimensional universe where files are stars. Related files cluster together into "Star Colonies" based on their metadata and content signatures.

### Documentation
- **Theoretical Concept**: [English](_theory/concept.md#english-version) / [ภาษาไทย](_theory/concept.md#ภาษาไทย)
- **User Manual & Features**: [English](_doc/manual.md#english-version) / [ภาษาไทย](_doc/manual.md#ภาษาไทย)

## Features
- **26D Embedding Vector**: Extracts size, name, time, and content hash (SHA-256).
- **Smart Link Detection**: Folder proximity, filename prefix/suffix, number proximity (file1 ↔ file2).
- **Physics Simulation**: "Shake & Brake" stabilizer for peaceful node placement.
- **Interactive UX**: Zoom, Pan, Drag & Bounce, and real-time Scan Progress.
- **Server Management**: Restart/Kill server and batch loop directly from the UI.
- **Memory Optimized**: Configurable batch size for low-RAM systems.

## Getting Started
1. Clone the repository.
2. Run `RUN.bat`.
3. Open `http://localhost:8080` in your browser.

## Command Line Options

```bash
file_graph_server.exe [options]
```

### Options
- `-startpath=<folder>` - Initial folder to scan on startup
- `-port=<number>` - Port to listen on (default: 8080)
- `-threshold=<0.0-1.0>` - Similarity threshold (default: 0.75, lower = more links)
- `-batch=<number>` - Batch size for scanning (default: 1000, lower = less memory)
- `-low_ram` - Low memory mode (~300 batch, 0.6 threshold)
- `-ram8g` - 8GB RAM profile (~500 batch, 0.65 threshold) - **Recommended for your system**
- `-ram16g` - 16GB RAM profile (~800 batch, 0.7 threshold)

### Link Detection Bonuses

| Condition | Bonus |
|-----------|-------|
| Same folder | +0.40 |
| Same number (file1.txt = file1.txt) | +0.30 |
| Adjacent numbers (file1.txt ↔ file2.txt) | +0.20 |
| Same extension (.go, .txt) | +0.15 |
| Same name suffix (3+ chars) | +0.03 per char |
| Same name prefix (3+ chars) | +0.05 per char |
| Same last 4 chars | +0.10 |
| Same size last 3 digits | +0.05 |

### Examples

```bash
# Scan current folder (default)
file_graph_server.exe

# Scan specific folder (use double quotes)
file_graph_server.exe -startpath="C:\My Projects"

# Custom port
file_graph for paths with spaces_server.exe -port=9000

# 8GB RAM mode (RECOMMENDED for your system)
file_graph_server.exe -startpath=C:\myproject -ram8g

# Low RAM mode (for systems with <2GB free)
file_graph_server.exe -startpath=C:\myproject -low_ram

# Manual settings
file_graph_server.exe -startpath=C:\myproject -batch=500 -threshold=0.6
```

## License
Distributed under the **MIT License**. See `LICENSE` for more information.
