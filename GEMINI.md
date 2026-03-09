# Project Gemini: File Graph Visualizer (Star Colony)

This file documents the core philosophy, theory, and architecture of the project for Gemini AI Team Assistance.

## 1. Core Philosophy
The project's guiding philosophy is to transform the abstract, hierarchical structure of a filesystem into a tangible, high-dimensional data universe.
- **Files as Stars**: Every file is treated as a "star" with unique properties.
- **Relationships as Colonies**: Files with similar attributes and content signatures are drawn together by simulated gravity, forming "Star Colonies" that reveal hidden relationships.
- **Exploration over Hierarchy**: The goal is to enable data discovery through visual patterns rather than just navigating a rigid folder tree.

## 2. Theoretical Framework
The system is built on three core theoretical components:
1.  **26-Dimensional Vector Embedding**: We use a hand-crafted feature extractor to map each file into a 26D vector space. This vector represents a file's "profile" based on its identity, time patterns, name signature, and content hash.
2.  **Cosine Similarity**: To measure the relationship between two files, we calculate the cosine similarity of their vectors. This metric is ideal as it measures the *angle* (profile similarity) rather than the magnitude (raw size/age).
3.  **Force-Directed Physics with Cooling**: A "Shake & Brake" physics simulation is used for visualization. It starts with high energy to "shake" nodes apart and then "brakes" by applying increasing friction as the simulation cools, allowing the graph to settle quickly and peacefully.

## 3. Program Architecture & Data Folders
The program is a modular Go web application with a JavaScript/Canvas frontend.

- **Backend (`internal/`)**: The Go backend is structured into packages for handling models, scanning, embedding, logging, and API requests.
- **Frontend (`index.html`)**: A single HTML file containing the UI, visualization logic (Canvas), and interaction handlers (Zoom, Pan, Drag, Controls).

### "Data Parsing" Folders

Within this project, there is one main folder that contains generated data:
- **`logs/`**: This folder contains runtime diagnostic logs (`app.log`). It is crucial for debugging what the server is doing during a scan or when an error occurs. This folder is a data output, not source code, and is therefore ignored by Git.

Any other folders, such as the `C:\xampp...` path you provided, are considered **external data inputs** for the program to scan. They are not part of the project's codebase and should not be added to its `.gitignore`.
