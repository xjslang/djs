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

# Clean generated files
npm run clean
```

## Comparison with Standard JavaScript

### Without DJS (manual cleanup)

```javascript
try {
  let db = sqlite('mydata.db');
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

### With DJS (automatic cleanup and error handling)
```js
let db = sqlite('mydata.db') or |err| {
  console.log('Cannot connect to database', err)
  return
}
defer db.close()
```

DJS makes this cleaner and less error-prone.