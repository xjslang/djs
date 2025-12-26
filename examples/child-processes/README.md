# Child Processes Example

This example demonstrates how DJS's `defer` statement simplifies management of child processes, ensuring they are properly terminated even if the parent process exits unexpectedly.

## Features Demonstrated

- **`defer` for automatic process cleanup**: Ensures child processes are killed when parent exits
- **Multiple concurrent processes**: Managing several child processes simultaneously
- **Process timeouts**: Automatically killing processes that run too long
- **Output handling**: Capturing and processing child process output
- **Nested processes**: Multiple processes with LIFO cleanup order
- **Error handling with `or`**: Clean error handling for failed processes

## What the Example Does

### Test 1: Simple Command
- Runs a basic `ls -la` command
- Automatically cleans up the child process
- Shows basic process execution

### Test 2: Long-Running Process
- Starts a worker that would run for 10 seconds
- Exits after only 2 seconds
- Child process is automatically killed with `defer`

### Test 3: Multiple Concurrent Processes
- Spawns 3 worker processes simultaneously
- All run concurrently for 2 seconds
- All are automatically cleaned up

### Test 4: Process with Timeout
- Starts a long-running worker
- Sets a 2.5 second timeout
- Automatically kills process when timeout is reached
- Both timeout and process are cleaned up with `defer`

### Test 5: Processing Output
- Executes a simple echo command
- Captures and processes the output
- Shows how to handle stdout

### Test 6: Nested Processes
- Starts three processes at different times
- Demonstrates LIFO cleanup order
- Shows that `defer` handles nested resources correctly

### Test 7: Failed Command
- Tries to execute a non-existent command
- Properly handles the error with `or`
- Shows error handling patterns

## Key DJS Features

### Automatic Process Termination
```djs
let child = spawn('long-running-command', ['arg1', 'arg2'])
defer child.kill()

await waitForProcess(child)
// Process killed automatically if parent exits early
```

### Multiple Processes
```djs
let child1 = spawn('command1', [])
defer child1.kill()

let child2 = spawn('command2', [])
defer child2.kill()

// Both processes cleaned up automatically in LIFO order
```

### Process with Timeout
```djs
let child = spawn('worker', [])
defer child.kill()

let timeout = setTimeout(function() {
  child.kill()
}, 5000)
defer clearTimeout(timeout)

// Both timer and process are cleaned up
```

### Error Handling
```djs
let child = spawn('command', [])
defer child.kill()

await waitForProcess(child) or |err| {
  console.error('Process failed:', err.message)
  return
}
// Process cleaned up even on error
```

## Running the Example

```bash
# Run the transpiled JavaScript directly
npm start

# Or build from DJS source and run
npm run dev
```

## Expected Output

```
=== DJS Child Processes Example ===
Demonstrates automatic cleanup of child processes using defer

📋 Example 1: Running simple command (ls)
   total 32
   drwxr-xr-x  5 user  staff   160 Dec 26 15:30 .
   drwxr-xr-x  8 user  staff   256 Dec 26 15:25 ..
   -rw-r--r--  1 user  staff  1234 Dec 26 15:30 example.djs
   ...
   ✅ Command completed

⏳ Example 2: Long-running process (will be terminated early)
   [LongWorker] Started (will run for 10000ms)
   [LongWorker] Working... (500ms elapsed)
   [LongWorker] Working... (1000ms elapsed)
   [LongWorker] Working... (1500ms elapsed)
   ⚠️  Parent function ending - child process will be killed automatically

🔀 Example 3: Multiple concurrent processes
   [Worker-1] Started (will run for 2000ms)
   [Worker-2] Started (will run for 2000ms)
   [Worker-3] Started (will run for 2000ms)
   [Worker-1] Working... (500ms elapsed)
   [Worker-2] Working... (500ms elapsed)
   [Worker-3] Working... (500ms elapsed)
   ...
   [Worker-1] Completed successfully
   [Worker-2] Completed successfully
   [Worker-3] Completed successfully
   ✅ All processes completed successfully

⏱️  Example 4: Process with timeout
   [TimeoutWorker] Started (will run for 10000ms)
   [TimeoutWorker] Working... (500ms elapsed)
   [TimeoutWorker] Working... (1000ms elapsed)
   [TimeoutWorker] Working... (1500ms elapsed)
   [TimeoutWorker] Working... (2000ms elapsed)
   ⏰ Timeout reached! Killing process...
   ✅ Process terminated due to timeout (as expected)

📊 Example 5: Processing command output
   Received output: "Hello from child process!"
   Output length: 26 characters
   ✅ Output processed successfully

🪆 Example 6: Nested processes with LIFO cleanup
   Started outer process
   [Outer] Started (will run for 4000ms)
   Started middle process
   [Middle] Started (will run for 3000ms)
   Started inner process
   [Inner] Started (will run for 2000ms)
   [Outer] Working... (1000ms elapsed)
   [Middle] Working... (800ms elapsed)
   [Inner] Working... (600ms elapsed)
   [Inner] Working... (1200ms elapsed)
   [Inner] Working... (1800ms elapsed)
   [Inner] Completed successfully
   ⚠️  Function ending - processes will be killed in LIFO order: Inner → Middle → Outer

❌ Example 7: Handling failed command
   ✅ Caught error as expected: spawn nonexistent-command ENOENT

=== All tests completed! ===
Notice how all child processes were cleaned up automatically.
```

## Comparison with Standard JavaScript

### Without DJS (manual cleanup)
```javascript
const child = spawn('long-process', []);
try {
  await waitForProcess(child);
} finally {
  child.kill(); // Must remember to kill!
}
```

### With DJS (automatic cleanup)
```djs
let child = spawn('long-process', [])
defer child.kill()

await waitForProcess(child)
// Process killed automatically!
```

## Common Use Cases

### Data Processing Pipeline
```djs
let child = spawn('data-processor', ['input.dat'])
defer child.kill()

await waitForProcess(child)
```

### Background Worker
```djs
let worker = spawn('node', ['worker.js'])
defer worker.kill()

await mainApplication()
// Worker terminated when app exits
```

### Process with Resource Limits
```djs
let child = spawn('memory-intensive-task', [])
defer child.kill()

let timeout = setTimeout(function() {
  child.kill('SIGTERM')
}, 30000)
defer clearTimeout(timeout)

await waitForProcess(child)
```

## Benefits

1. **No orphaned processes**: Child processes are always terminated
2. **Cleaner code**: No nested try-finally blocks
3. **Guaranteed cleanup**: Even if parent crashes or exits early
4. **LIFO order**: Multiple processes cleaned up in reverse order
5. **Resource safety**: Prevents runaway processes consuming resources

## Implementation Notes

- The `worker.js` script is a simple Node.js script that simulates long-running work
- Uses standard Unix commands (`ls`, `echo`) that work on macOS/Linux
- On Windows, you might need to adjust commands (e.g., use `dir` instead of `ls`)
- The example properly handles SIGTERM and SIGINT signals in child processes
- All processes are killed gracefully when the parent function exits

## Important Notes

- Child processes are killed when the function exits, regardless of how it exits
- Multiple `defer` statements execute in LIFO (Last-In-First-Out) order
- Killing a process sends SIGTERM by default (can be changed with `child.kill('SIGKILL')`)
- Always handle process errors with `or` blocks to prevent unhandled promise rejections
