# HTTP Resource Cleanup

This example demonstrates how `or` and `async/await` constructs simplify HTTP resource management in real-world scenarios.

## Run the Example

```bash
# install djs
go install github.com/xjslang/djs@latest

# Terminal 1: Start the test server
npm run server

# Terminal 2: Transpile and run the example
npm start

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

### 4. **HTTP Connection Management**

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