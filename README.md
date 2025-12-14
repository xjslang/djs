# DeferJavaScript (DJS)

A JavaScript dialect that adds Go-style `defer` statements, error-handling `or` blocks, and strict equality by default.

## Features

- **`defer`**: Execute cleanup code when functions exit (Go-style)
- **`or` blocks**: Elegant error handling fallbacks
- **Strict equality**: `==` behaves like `===`

## Installation

```bash
go install github.com/xjslang/djs@latest
```

Or build from source:
```bash
git clone https://github.com/xjslang/djs.git
cd djs
go build -o djs
```

## Usage

### Execute directly
```bash
djs script.djs
```

### Transpile to JavaScript
```bash
# Basic transpilation
djs -o output.js script.djs

# With external source map
djs -o output.js --sourcemap script.djs

# With inline source map
djs -o output.js --inline-sourcemap script.djs

# With source content embedded
djs -o output.js --sourcemap --inline-sources script.djs

# Custom paths
djs -o output.js --sourcemap \
  --map-root "/maps/" \
  --source-root "/src/" \
  script.djs
```

## Language Examples

### Defer statement
```javascript
function processFile(path) {
    let file = openFile(path);
    defer closeFile(file);
    
    // File closes automatically when function exits
    return parseContent(file);
}
```

### Or blocks (error handling)
```javascript
let data = fetchData() or {
    console.log("Using fallback");
    return defaultData;
};
```

### Strict equality
```javascript
// In DJS, == works like ===
if (value == "123") {  // Only matches string "123"
    console.log("Strict comparison");
}
```

## Source Map Options

- `--sourcemap`: Generate external `.map` file
- `--inline-sourcemap`: Embed source map as base64
- `--inline-sources`: Include source content in map
- `--map-root`: Set location for `.map` file
- `--source-root`: Set prefix for source files

## Development

```bash
# Run tests
mage test

# Lint code
mage lint
```
