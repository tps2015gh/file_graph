# Red Energy Gravity: Search-Driven Attractor Physics

## 1. Theory: The "Search Vortex" Concept
Unlike standard semantic colonies that form naturally via pairwise cosine similarity, **Red Energy Gravity** is a synthetic, high-intensity force field triggered by user queries. It introduces a "Centralized Attractor" model into the distributed N-body simulation to facilitate rapid data discovery.

### Physical Characteristics:
- **Offset Genesis**: To maintain the "Star Colony" context, the attractor is born at `(Center_of_Results.x + 300, Center_of_Results.y)`. This creates a visual "migration" from the old colony to the new search-defined space, allowing the user to compare the search results against the original cluster.
- **Draggable Singularity**: The user can manually reposition the attractor, effectively "steering" the search results across the data universe to investigate relationships with other clusters.
- **Vortex Dynamics**: As nodes approach the attractor, a tangential velocity component (swirl) is added. This prevents "node stacking" and creates a more fluid, dynamic clustering behavior that feels "alive."

## 2. Force Mathematical Model
The total force $\vec{F}_{total}$ applied to a search result node $n$ by the Red Energy Attractor $A$ is the vector sum of attraction and orbital swirl:

$$ \vec{F}_{total} = \vec{F}_{attraction} + \vec{F}_{swirl} $$

### Attraction Component (Linear Decay)
The attraction force pulls nodes toward the singularity with a strength that diminishes as they approach the center, allowing them to settle into a stable orbit:

$$ \vec{F}_{attraction} = 	ext{energy} \cdot (1 - \frac{d}{R}) \cdot \hat{u} $$

Where:
- $d$ is the Euclidean distance between node $n$ and attractor $A$.
- $R$ is the attractor radius (2000px).
- $\hat{u}$ is the unit vector from $n$ to $A$.
- $	ext{energy}$ is the dynamic amplitude of the vortex (120 units).

### Swirl (Vortex) Component
Activated when $d < 0.2R$ to create a "swirling" motion:

$$ \vec{F}_{swirl} = (1 - \frac{d}{0.2R}) \cdot 2 \cdot \hat{v} $$

Where $\hat{v}$ is the perpendicular vector $(-\Delta y, \Delta x)$ providing the orbital "spin" or vortex effect.

## 3. Persistent Highlighting & Labeling
Red Energy nodes are elevated to a higher visual tier in the rendering pipeline:
- **Red Priority**: Node color is locked to `#ff4d4d` with a high-intensity glow.
- **Label Persistence**: Labels are forced to render regardless of global optimization settings (like `nodeCount < 100`), ensuring immediate identification.
- **Opacity Lock**: Opacity is maintained at `1.0` even when other nodes are selected, preventing search results from fading into the background.

## 4. Vortex Visualization (UI Display)
The gravity source is rendered as a multi-layered aesthetic element:
- **Core Singularity**: A solid `#ff4d4d` circle with a `25px` shadow blur, representing the densest part of the energy field.
- **Interaction Ring**: A white dashed circle (`setLineDash([5, 5])`) at $1.5R$ radius, signaling to the user that the element is "Grabable" and interactive.
- **Atmospheric Glow**: A radial gradient starting from $50\%$ red opacity at the center to $0\%$ at $3R$, creating a volumetric "atmosphere" for the attractor.
- **Dynamic Feedback**: A textual indicator `DRAG RED ENERGY` is anchored above the core whenever the user is zoomed in enough ($k > 0.5$) to facilitate immediate UX understanding.

---

### Author Signature
**Written by**: Gemini 3
**Model**: Google Gemini 2.0 Flash (Advanced CLI Integration)
**Role**: Senior AI Architect & Physics Simulation Specialist
**Date**: 2026-03-10
**Instance ID**: 69ca6229-1ccc-4e80-8f3c-86a9ad1af00c

---
**Related Documents**:
- [`force_dynamics_correction.md`](force_dynamics_correction.md) - Base physics engine
- [`concept.md`](concept.md) - 26D Vector Space theory
