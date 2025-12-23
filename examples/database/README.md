# Database Connection Management

This example demonstrates how `defer` simplifies database cleanup in real-world scenarios.

## Setup

```bash
# install djs
go install github.com/xjslang/djs@latest

# install libraries and create database
npm install
npm run setup
```

## Run the Example

### Basic Example

```bash
# Quick start: transpile and run
npm start

# Or step by step:
npm run build  # Transpile example.djs → example.js
node example.js
```

### Advanced Example

The advanced example demonstrates multiple queries and more complex error handling:

```bash
npm run start:advanced

# Or step by step:
npm run build:advanced  # Transpile advanced.djs → advanced.js
node advanced.js
```

### Cleanup

```bash
# Clean generated files
npm run clean
```

## What it Demonstrates

### 1. **Automatic Resource Cleanup**
The `defer` statement ensures the database connection is closed regardless of success or errors:

```djs
let db = sqlite('mydata.db') or |err| {
  console.log('Cannot connect to database', err)
  return
}
defer db.close()  // ← Always executed before function returns
```

### 2. **Error Handling with `or` Blocks**
DJS provides cleaner error handling compared to traditional try-catch:

```djs
let db = sqlite('mydata.db') or |err| {
  console.log('Cannot connect to database', err)
  return
}
```

This transpiles to proper try-catch blocks automatically.

### 3. **Real-world Database Pattern**
This is a common pattern in production code:
- Open connection
- Register cleanup with defer
- Execute queries
- Automatic cleanup even on errors

### Traditional JavaScript Equivalent

Without DJS, you'd need verbose try-finally blocks:

```javascript
let db;
try {
  db = sqlite('mydata.db');
  try {
    let stmt = db.prepare('SELECT * FROM users WHERE active = ?');
    let users = stmt.all(1);
    console.log(`Found ${users.length} active users`);
  } finally {
    db.close();
  }
} catch (err) {
  console.log('Cannot connect to database', err);
}
```

DJS makes this cleaner and less error-prone.

## Expected Output

```
Found 4 active users
```

## Files

- `example.djs` - Basic example with defer and or
- `advanced.djs` - Advanced example with multiple queries
- `setup.js` - Database initialization script
- `package.json` - Node.js dependencies and scripts
- `.gitignore` - Ignore generated files

## How It Works

The DJS transpiler converts the `.djs` files to regular JavaScript:

**DJS Source (`example.djs`):**
```djs
let db = sqlite('mydata.db') or |err| {
  console.log('Cannot connect to database', err)
  return
}
defer db.close()
```

**Transpiled JavaScript (`example.js`):**
```javascript
let defers = [];
try {
  let db;
  try {
    db = sqlite("mydata.db")
  } catch(err) {
    console.log("Cannot connect to database", err);
    return
  }
  defers.push(() => { db.close() });
  // ... rest of the code
} finally {
  // Execute all defers in reverse order
  for(let i = defers.length; i > 0; i--) {
    try { defers[i-1]() } catch(e) { console.log(e) }
  }
}
```

This shows how DJS provides cleaner syntax while maintaining full JavaScript compatibility.
