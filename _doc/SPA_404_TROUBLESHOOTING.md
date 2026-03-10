# Single-Page Application (SPA) 404 Error Troubleshooting

## Issue Description
Users encounter "404 Page Not Found" errors when trying to access URLs other than the root path.

## Why This Happens

### Technical Architecture
This application is built as a **Single-Page Application (SPA)** where:
- All routing is handled client-side by JavaScript
- The Go server serves `index.html` for ALL unmatched routes
- No server-side routing exists beyond the API endpoints

### Valid vs Invalid URLs

| URL Pattern | Status | Explanation |
|-------------|--------|-------------|
| `http://localhost:8080/` | ✅ Works | Root path served correctly |
| `http://localhost:8080/index.html` | ❌ 404 | Direct file access not supported |
| `http://localhost:8080/graph` | ❌ 404 | Client-side route, server doesn't know |
| `http://localhost:8080/search` | ❌ 404 | Client-side route, server doesn't know |

## Server Routing Logic

```go
func ServeIndex(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)  // This causes 404 for non-root paths
        return
    }
    http.ServeFile(w, r, "index.html")
}
```

## Solutions

### For Users
1. **Always use the root URL**: `http://localhost:8080/`
2. **Bookmark the correct URL**: Don't bookmark sub-pages
3. **Use browser navigation**: Back/forward buttons work within the SPA

### For Developers (Future Enhancement)
To support direct URLs, implement client-side routing:

```javascript
// Hash-based routing
window.addEventListener('hashchange', () => {
    const route = window.location.hash.slice(1);
    // Handle route changes
});

// Or history API routing
window.addEventListener('popstate', () => {
    const route = window.location.pathname;
    // Handle route changes
});
```

## Related Files

- `main.go` - Server routing configuration
- `index.html` - Client-side application
- `internal/handlers/handlers.go` - HTTP request handling

## Common Scenarios

### Scenario 1: Browser Refresh on Subpage
**Problem**: User refreshes browser on `/graph` path
**Solution**: Redirect to root path in JavaScript

### Scenario 2: Shared URLs
**Problem**: User shares `http://localhost:8080/search?q=test`
**Solution**: Use URL parameters instead of paths: `http://localhost:8080/?search=test`

### Scenario 3: Development Environment
**Problem**: IDE/debugger tries to open specific files
**Solution**: Always start from root URL

## Troubleshooting Steps

1. **Check current URL** - Ensure it's exactly `http://localhost:8080/`
2. **Clear browser cache** - Sometimes old routes are cached
3. **Restart server** - Ensure server is running correctly
4. **Check console errors** - Look for JavaScript errors

---

**Note**: This is intentional SPA behavior, not a bug. The application is designed to work entirely from the root path.