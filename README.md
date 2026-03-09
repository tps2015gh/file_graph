# File Graph Visualizer (Star Colony)

A standalone Go web application to scan directories, extract file metadata, calculate similarity embeddings (26D Vector Space), and visualize relationships as an interactive 2D interactive graph.

## Project Roles & Direction
- **Lead Programmer & Director**: The User
- **AI Team Assistance** (Last Updated: 2026-03-09):
  - Gemini CLI (Powered by multiple models)
  - OpenCode (Advanced AI coding assistant)
  - Minimax-M2.5-Free (Fast code generation model via OpenCode)

## Core Concept
This program treats the filesystem as a high-dimensional universe where files are stars. Related files cluster together into "Star Colonies" based on their metadata and content signatures.

### Documentation
- **Theoretical Concept**: [English](_theory/concept.md#english-version) / [ภาษาไทย](_theory/concept.md#ภาษาไทย)
- **User Manual & Features**: [English](_doc/manual.md#english-version) / [ภาษาไทย](_doc/manual.md#ภาษาไทย)

## Features
- **26D Embedding Vector**: Extracts size, name, time, and content hash (SHA-256).
- **Smart Link Detection**: Folder proximity, filename prefix/suffix, reduced number proximity (prevents over-clustering).
- **Physics Simulation**: "Shake & Brake" stabilizer for peaceful node placement.
- **Interactive UX**: Zoom, Pan, Drag & Bounce, and real-time Scan Progress.
- **Link Filter**: Slider to show only top 1-100% strongest links (prevents mesh sphere effect).
- **Zoom Controls**: Zoom+/Zoom- buttons (2x per click), Reset to 100%.
- **Shake Button**: Add random energy to break stable clusters.
- **Node Details**: Full path display with Copy Path button, Open in Explorer.
- **Server Management**: Restart/Kill server and batch loop directly from the UI.
- **Memory Optimized**: Configurable batch size for low-RAM systems (-low_ram, -ram8g, -ram16g).

## Screenshots

![File Graph Visualizer](_doc/Screenshot-2026-03-10-000725.png)

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
| Same number (file1.txt = file1.txt) | +0.15 |
| Adjacent numbers (file1.txt ↔ file2.txt) | +0.08 |
| Close numbers (diff ≤3) | +0.03 |
| Same extension (.go, .txt) | +0.15 |
| Same name suffix (3+ chars) | +0.03 per char |
| Same name prefix (3+ chars) | +0.05 per char |
| Same last 4 chars | +0.10 |
| Same size last 3 digits | +0.05 |

### UI Controls

- **Spacing slider**: Adjust repulsion force between nodes
- **Rotate slider**: Add rotation to spread nodes radially
- **Links slider (1-100%)**: Filter to show only top X% strongest links
- **Zoom+ / Zoom-**: Zoom in/out (2x per click, max 50x)
- **1:1**: Reset zoom to 100%
- **Shake**: Add random energy to break stable clusters
- **Reset**: Reset links filter to 100%

### Examples

```bash
# Show help
file_graph_server.exe

# Scan specific folder (use double quotes)
file_graph_server.exe -startpath="C:\My Projects"

# Custom port
file_graph_server.exe -port=9000

# 8GB RAM mode (RECOMMENDED for your system)
file_graph_server.exe -startpath=C:\myproject -ram8g

# Low RAM mode (for systems with <2GB free)
file_graph_server.exe -startpath=C:\myproject -low_ram

# Manual settings
file_graph_server.exe -startpath=C:\myproject -batch=500 -threshold=0.6
```

## Known Issues & Limitations

### Performance
- **Large folders may cause browser hang or slow performance**: Scanning directories with thousands of files can consume significant memory and CPU. Use `-low_ram` or `-ram8g` flags to reduce batch size and improve stability.

### UI/UX
- **Navigation bar overflow**: On smaller screens or when using many controls, the toolbar may extend beyond the visible area. Use the browser's horizontal scroll or reduce the number of visible controls.
- **Spacing and Rotate sliders**: These controls currently have limited effect on the force simulation. A future update will enhance their functionality.

### Functional Bugs
- **Browse folder button**: The folder browser dialog may not work correctly in all browsers. Please enter the folder path manually as a workaround.
- **Unicode/Thai language paths**: Some paths containing Thai characters or special Unicode characters may cause JSON parsing errors. Avoid scanning folders with non-ASCII names.

## License
Distributed under the **MIT License**. See `LICENSE` for more information.
