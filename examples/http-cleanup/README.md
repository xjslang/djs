# HTTP Resource Cleanup

This example demonstrates how the `defer` construct simplifies HTTP resource management in real-world scenarios with Node.js-style callbacks.

## Setup

```bash
# install djs
go install github.com/xjslang/djs@latest
```

## Run the Example

### Quick Start

```bash
# Terminal 1: Start the test server
npm run server

# Terminal 2: Transpile and run the example
npm start
```

### Step by Step

```bash
# Terminal 1: Start the HTTP server
npm run server

# Terminal 2: Build and run
npm run build  # Transpile example.djs → example.js
node example.js
```

### Automated Test

```bash
# Run the full integration test
npm test
```

This will automatically start the server, run the example, and stop the server.

### Cleanup

```bash
# Clean generated files
npm run clean
```

## What it Demonstrates

### 1. **Automatic Resource Cleanup with `defer`**
```djs
let agent = new http.Agent({ keepAlive: true, maxSockets: 10 })
defer agent.destroy()  // Guaranteed cleanup even if errors occur

// Make request
defer closeResponse(response)  // Cleanup response resources
```

The `defer` statement ensures that resources are cleaned up in reverse order when the function exits, regardless of whether it exits normally or due to an error. This is especially useful for managing HTTP connections, file handles, and other resources that need cleanup.

### 2. **Error Handling with Callbacks**
```djs
makeRequest(url, agent, function(response, error) {
  if (error) {
    console.error('Request failed:', error.message)
    return  // Early exit on error
  }
  // Process successful response
})
```

The example uses callback-style error handling (Node.js style with error-first callbacks), making it compatible with standard Node.js patterns while demonstrating how `defer` works with asynchronous operations.

### 3. **HTTP Connection Management**
- Creates HTTP agent with connection pooling
- Properly closes sockets after requests
- Destroys the agent when done
- All cleanup is automatic via `defer`

## The Test Server

The example includes a simple HTTP server (`server.js`) that provides:

- `GET /api/users/:id` - Fetch user data
- `GET /api/users/:id/posts` - Fetch user's posts
- `GET /api/health` - Health check endpoint

## How It Works

1. **Resource Allocation**: Creates an HTTP agent with connection pooling
2. **Request Execution**: Makes HTTP requests to fetch user data and posts
3. **Error Handling**: Uses `or` blocks to handle request failures gracefully
4. **Cleanup**: All resources (agent, sockets) are cleaned up automatically via `defer`

The transpiled JavaScript code uses try-finally blocks to implement the `defer` semantics, ensuring cleanup happens even when errors occur in the callback chain. The DJS source code remains clean and readable.

## Expected Output

```
Starting HTTP cleanup example...
Fetching data from local server (http://localhost:3000)

✓ User fetched: Alice Johnson (alice@example.com)
✓ Found 3 posts:

  - Getting Started with DJS
  - HTTP Resource Management
  - Error Handling Made Easy

✓ All resources cleaned up via defer statements
```

## Error Scenarios

If the server is not running, you'll see:
```
Failed to fetch user data: connect ECONNREFUSED 127.0.0.1:3000
Make sure to run: npm run server
```

The error callback catches the error and provides a helpful message, while `defer` ensures any allocated resources are still cleaned up.
