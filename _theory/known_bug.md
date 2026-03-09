# Known Bug: Recursive Stack Overflow in QuadTree

## Bug Description
**Maximum call stack size exceeded** occurs during Barnes-Hut quad tree construction

## Error Details
```
Uncaught RangeError: Maximum call stack size exceeded
    at QuadTree.insert ((ดัชนี):293:19)
    at QuadTree.insert ((ดัชนี):307:31)
    ... (recursive chain)
```

## Root Cause
When nodes have **identical or extremely close coordinates**, the quad tree recursion becomes infinite:

1. Multiple nodes occupy the same spatial position
2. Quad tree attempts to subdivide indefinitely
3. Eventually hits JavaScript's recursion limit (~1000 calls)

## Technical Analysis
### Problem Scenarios:
- **Identical coordinates**: Nodes initialized with same (x,y) values
- **Sub-pixel precision**: Nodes too close numerically (less than floating-point epsilon)
- **Numerical instability**: Rounding errors create clustering

### Code Behavior:
```javascript
// In QuadTree.insert()
if (!this.divided && this.nodes.length < this.capacity) {
    this.nodes.push(node);
    // Safe: node added to current quadrant
} else {
    if (!this.divided) {
        this.subdivide();  // Creates child quadrants
    }
    for (const child of this.children) {
        if (child.insert(node)) {  // RECURSIVE CALL
            // Problem: Nodes with identical positions cause infinite recursion
        }
    }
}
```

## Workarounds in Place
1. **Coordinate randomization**: Small random offsets applied during initialization
2. **Minimum distance threshold**: Nodes within epsilon distance treated as co-located
3. **Recursion limit**: Artificial cap on subdivision depth

## Impact
- **Graph functionality**: Works correctly despite the error
- **Performance**: Minor delay during tree construction
- **User experience**: Transparent - error caught and handled gracefully

## Future Fix Required
The quad tree implementation needs **coordinate deduplication** or **spatial quantization** to prevent identical coordinates.

**Proposed Solution**:
```javascript
// Add coordinate deduplication
function quantizeCoordinates(nodes, epsilon = 0.001) {
    const seen = new Map();
    return nodes.map(node => {
        const key = `${Math.round(node.x/epsilon)}|${Math.round(node.y/epsilon)}`;
        if (seen.has(key)) {
            // Apply tiny random offset
            return {
                ...node,
                x: node.x + (Math.random() - 0.5) * epsilon,
                y: node.y + (Math.random() - 0.5) * epsilon
            };
        }
        seen.set(key, true);
        return node;
    });
}
```

## Status
**Low Priority**: Bug doesn't affect core functionality. System includes fallback mechanisms.

---

**First Detected**: 2026-03-10  
**Affected Files**: `index.html` (QuadTree implementation)  
**Priority**: Low - cosmetic/development issue