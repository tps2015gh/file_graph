# Drag Behavior & Link Relationships Explained

## Current Dragging Behavior

### What You're Seeing:
When you drag a file node, **all semantically related files move together** as a cohesive group.

## How Link Detection Works

### 1. **Direct Connections (Current Logic)**
```
File A --similar--> File B --similar--> File C
```
If you drag **File A**:
- File A moves ✅ (dragged node)
- File B moves ✅ (directly connected)
- File C does NOT move ❌ (not directly connected to A)

### 2. **Why "Far Away" Nodes Might Move**
If files share **high similarity scores**, they form **clusters** that move together:

```
Cluster 1: File A -- File B -- File C
Cluster 2: File D -- File E -- File F
```

Even if File A and File D seem "far apart" visually, if they have **high similarity**, they belong to the same semantic cluster.

## Link Strength & Threshold

### Similarity Threshold: 75%
- Files with **>75% similarity** are considered "related"
- Similarity calculated using 26D vector embeddings
- Based on file content, structure, and metadata

### What Creates Links:
1. **Similar file extensions** (.js, .html, .css)
2. **Similar content patterns** 
3. **Related functionality**
4. **Directory proximity**

## When Dragging Makes Sense vs. Confusion

### ✅ Intended Use Cases:
- Moving **related functionality components** together
- Exploring **semantic clusters** 
- Understanding **file relationships** visually

### ❌ Potential Confusion:
- Files that "look unrelated" but share deep semantic connections
- Large clusters moving together unexpectedly

## Customizing Link Detection

You can adjust the similarity threshold:
- **Higher threshold (e.g., 80%)**: Fewer links, more focused clusters
- **Lower threshold (e.g., 60%)**: More links, larger clusters

## Visual Indicators

### During Drag:
- 🔵 **Dragged node**: Original blue highlight
- 🟢 **Connected nodes**: Green highlight
- 📝 **Labels**: All connected node names visible

### Normal View:
- Only selected node shows relationships
- Gray links indicate similarity strength

## Common Questions

### Q: Why do files from different folders move together?
A: They share high semantic similarity despite folder separation.

### Q: Can I move individual files only?
A: No - the intent is to explore relationships. Use search/filter for individual files.

### Q: How to identify why files are linked?
A: Select a node and check the "Related Files" panel for similarity scores.

## Technical Details

### Link Creation:
- 26-dimensional vector comparison
- Cosine similarity algorithm
- Automatic threshold-based filtering

### Drag Algorithm:
- Finds all nodes with **direct similarity links**
- Applies same movement delta to connected nodes
- Maintains relative positioning

---

**Tip**: Use the search function (`Ctrl+F`) to find specific files without dragging clusters.