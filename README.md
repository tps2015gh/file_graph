# File Graph Visualizer (Star Colony)

A standalone Go web application to scan directories, extract file metadata, calculate similarity embeddings (26D Vector Space), and visualize relationships as an interactive 2D interactive graph.

## Project Roles & Direction
- **Lead Programmer & Director**: The User
- **AI Team Assistance**: Gemini CLI (Powered by multiple models)

## Core Concept
This program treats the filesystem as a high-dimensional universe where files are stars. Related files cluster together into "Star Colonies" based on their metadata and content signatures.

- **[Theoretical Concept (ภาษาไทย/English)](_theory/concept.md)**
- **[User Manual & Features (ภาษาไทย/English)](_doc/manual.md)**

## Features
- **26D Embedding Vector**: Extracts size, name, time, and content hash (SHA-256).
- **Physics Simulation**: "Shake & Brake" stabilizer for peaceful node placement.
- **Interactive UX**: Zoom, Pan, Drag & Bounce, and real-time Scan Progress.
- **Server Management**: Restart/Kill server and batch loop directly from the UI.

## Getting Started
1. Clone the repository.
2. Run `RUN.bat`.
3. Open `http://localhost:8080` in your browser.

## License
Distributed under the **MIT License**. See `LICENSE` for more information.
