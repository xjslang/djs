# HTTP Resource Cleanup

This example demonstrates how `or` and `async/await` constructs simplify HTTP resource management in real-world scenarios.

> **Note:** There's currently a known issue with using `defer` inside `async` functions. This example uses manual cleanup via a `cleanup()` function to demonstrate the async/await and or constructs. Once the issue is resolved, defer statements can be added for automatic resource cleanup.

## Setup

```bash
# install djs
go install github.com/xjslang/djs@latest
```

## Run the Example

```bash
# Terminal 1: Start the test server
npm run server

# Terminal 2: Transpile and run the example
npm start
```

This will automatically start the server, run the example, and stop the server.

### Cleanup

```bash
# Clean generated files
npm run clean
```

## What it Demonstrates

### 1. **Async/Await for Sequential Operations**

```djs
async function fetchUserData() {
  // First request
  let userResponse = await makeRequest('/api/users/123', agent)
  let userData = JSON.parse(userResponse.data)

  // Second request using data from first
  let postsResponse = await makeRequest(`/api/users/${userData.id}/posts`, agent)
  let posts = JSON.parse(postsResponse.data)
}
```

The `await` keyword pauses execution until the Promise resolves, allowing you to write asynchronous code in a sequential, readable manner.

### 2. **Clean Error Handling with `or`**

```djs
let userResponse = await makeRequest(url, agent) or |err| {
  console.error('Request failed:', err.message)
  cleanup()  // Ensure cleanup on error
  return  // Early exit on error
}
```

The `or` construct provides elegant error handling for async operations, making the happy path more prominent while ensuring cleanup happens even on errors.

### 3. **Resource Management**

```djs
// Create resources
let agent = new http.Agent({ keepAlive: true })
let logFile = fs.openSync('request.log', 'w')

// Use resources with async/await
await fetchUserData()

// Manual cleanup (defer in async functions has a known issue)
function cleanup() {
  agent.destroy()
  fs.closeSync(logFile)
}
```

Resources are created once and cleaned up properly after async operations complete. Future support for `defer` in `async` functions will make this automatic.

### 3. **HTTP Connection Management**

- Creates HTTP agent with connection pooling
- Uses file logging to track requests
- Manual cleanup ensures all resources are freed
- Async/await makes sequential HTTP requests readable and maintainable

## The Test Server

The example includes a simple HTTP server (`server.js`) that provides:

- `GET /api/users/:id` - Fetch user data
- `GET /api/users/:id/posts` - Fetch user's posts
- `GET /api/health` - Health check endpoint

## How It Works

1. **Resource Allocation**: Creates an HTTP agent and log file
2. **Async Requests**: Makes HTTP requests using `await` for sequential execution
3. **Error Handling**: Uses `or` blocks to handle request failures gracefully
4. **Manual Cleanup**: Calls `cleanup()` function to free resources

The transpiled JavaScript code uses `async/await` for Promise handling and try-catch for `or` blocks. Once `defer` support in `async` functions is fixed, cleanup will be automatic.

## Expected Output

```
Starting HTTP cleanup example...
Fetching data from local server (http://localhost:3000)

✓ User fetched: Alice Johnson (alice@example.com)
✓ Found 3 posts:

  - Getting Started with DJS
  - HTTP Resource Management
  - Error Handling Made Easy

✓ All requests completed successfully!

Cleaning up resources...
✓ HTTP agent destroyed
✓ Log file closed

Note: Check request.log for the full request log
```

## Error Scenarios

If the server is not running, you'll see:

```
Failed to fetch user data: connect ECONNREFUSED 127.0.0.1:3000
Make sure to run: npm run server

Cleaning up resources...
✓ HTTP agent destroyed
✓ Log file closed
```

The `or` block catches the error and provides a helpful message, while the `cleanup()` function ensures all allocated resources are freed, even when errors occur early in the function.
