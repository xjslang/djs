# Temporary Files Example

This example demonstrates how DJS's `defer` and `or` constructs simplify working with temporary files and directories.

## Features Demonstrated

- **`defer` for automatic cleanup**: Ensures temporary directory is removed even if errors occur
- **`or` for error handling**: Clean error handling without try-catch nesting
- **File operations**: Creating, reading, and writing temporary files
- **Directory management**: Working with temporary directories

## What the Example Does

1. Creates a temporary directory in the system temp folder
2. Uses `defer` to guarantee cleanup of the directory (even on errors)
3. Creates multiple files in the temp directory:
   - `data1.json` - JSON data file
   - `data2.txt` - Plain text file
   - `output.log` - Log file
4. Reads and processes the JSON file
5. Lists all files in the temp directory
6. Automatically cleans up the entire temp directory when done

## Key DJS Features

### Defer Statement
```djs
defer fs.rmdir(tmpDir, { recursive: true })
```
Guarantees the temp directory is cleaned up when the function exits, similar to Go's defer.

### Or Expression
```djs
let tmpDir = fs.mkdtemp(path.join(os.tmpdir(), 'djs-example-')) or |err| {
  console.error('Failed to create temp directory:', err.message)
  return
}
```
Provides clean error handling without deep try-catch nesting.

## Running the Example

```bash
# Run the transpiled JavaScript directly
npm start

# Or build from DJS source and run
npm run dev
```

## Expected Output

```
Created temporary directory: /tmp/djs-example-xxxxx
✓ Created: /tmp/djs-example-xxxxx/data1.json
✓ Created: /tmp/djs-example-xxxxx/data2.txt
✓ Read data: { name: 'Alice', age: 30, city: 'NYC' }
✓ Created log: /tmp/djs-example-xxxxx/output.log
✓ Files in temp directory: [ 'data1.json', 'data2.txt', 'output.log' ]

All operations completed successfully!
Temporary directory will be cleaned up automatically.
```

## Comparison with Standard JavaScript

### Without DJS (nested try-catch)
```javascript
const fs = require('fs').promises;
let tmpDir;
try {
  tmpDir = await fs.mkdtemp(path.join(os.tmpdir(), 'example-'));
  try {
    await fs.writeFile(file1, data);
    // ... more operations
  } finally {
    await fs.rmdir(tmpDir, { recursive: true });
  }
} catch (err) {
  console.error('Failed:', err.message);
}
```

### With DJS (clean and readable)
```djs
let tmpDir = fs.mkdtemp(path.join(os.tmpdir(), 'example-')) or |err| {
  console.error('Failed:', err.message)
  return
}
defer fs.rmdir(tmpDir, { recursive: true })

fs.writeFile(file1, data) or |err| {
  console.error('Failed:', err.message)
  return
}
// ... more operations - cleanup is automatic!
```

## Benefits

1. **No manual cleanup tracking**: `defer` handles it automatically
2. **Readable error handling**: `or` blocks are clearer than nested try-catch
3. **Guaranteed cleanup**: Even if an error occurs, the temp directory is removed
4. **Less boilerplate**: Cleaner, more maintainable code
