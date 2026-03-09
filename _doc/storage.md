# Data Storage: No Database

This project uses a **stateless, on-the-fly computation** approach - no database required.

## How It Works

1. **Scan on Demand**: When you click "Scan" or use `-startpath`, the server reads the directory in real-time.
2. **Metadata Extraction**: For each file, it extracts:
   - File size, modification time
   - Filename, extension
   - Content hash (SHA-256)
3. **Vector Embedding**: Each file is converted into a 26D vector representing its "identity"
4. **Relation Calculation**: Similarity between files is calculated on-the-fly using cosine similarity
5. **In-Memory**: All data exists only in server memory while running

## Data Flow

```
User Request (Scan)
       ↓
Directory Walker (filepath.WalkDir)
       ↓
Metadata Extractor (size, time, hash)
       ↓
Vector Builder (26D embedding)
       ↓
Similarity Calculator (N² comparison)
       ↓
JSON Response to Frontend
       ↓
Canvas Visualization
```

## Why No Database?

- **Simplicity**: No external dependencies (SQLite, PostgreSQL, etc.)
- **Speed**: Fresh data on every scan - no sync issues
- **Portable**: Just copy the .exe and run anywhere
- **Lightweight**: No storage overhead

## Limitations

- Data is lost when server stops
- Large directories take time to scan each time
- No historical data or versioning

## Future Options

If persistent storage is needed later:
- Export graph as JSON
- Cache vectors to disk
- Add SQLite for history tracking
