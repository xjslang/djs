# Defer for JavaScript (DJS)

DJS is a JavaScript dialect that simplifies resource management using the `defer` and `or` constructs, eliminating the verbosity of `try/catch/finally`. It is ideal for the DevOps Community, where proper resource cleanup is critical.

DJS is built on top of [XJS](https://github.com/xjslang/xjs), a super-fast, minimalist, and extensible JavaScript parser powered by plugins.

Below is an example of how DJS transpiles to standard JavaScript:

**Original code:**
```js
// input.djs
let sqlite = require('better-sqlite3')

function main() {
  let db = sqlite('mydata.db') or |err| {
    console.log('Cannot connect to database', err)
    return
  }
  defer db.close()

  // prepare and execute queries
  let stmt = db.prepare('SELECT * FROM users WHERE active = ?');
  let users = stmt.all(1);
  console.log(`Found ${users.length} active users`);
}

main()
```

**Transpiled code:**
```js
// output.js
let sqlite = require("better-sqlite3");

function main() {
  let defers = [];
  try {
    let db;
    try {
      db = sqlite("mydata.db");
    } catch (err) {
      console.log("Cannot connect to database", err);
      return;
    }
    defers.push(() => {
      db.close();
    });
    let stmt = db.prepare("SELECT * FROM users WHERE active = ?");
    let users = stmt.all(1);
    console.log(`Found ${users.length} active users`);
  } finally {
    for (let i = defers.length; i > 0; i--) {
      try {
        defers[i - 1]();
      } catch (e) {
        console.log(e);
      }
    }
  }
}

main();
```

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
