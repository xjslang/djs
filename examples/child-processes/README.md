# Child Processes Example

This example demonstrates how DJS's `defer` statement simplifies management of child processes, ensuring they are properly terminated even if the parent process exits unexpectedly.

## Running the Example

```bash
# install djs
go install github.com/xjslang/djs@latest

# Run the transpiled JavaScript directly
npm start

# Or build from DJS source and run
npm run dev
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
```js
let child = spawn('long-process', [])
defer child.kill()

await waitForProcess(child)
// Process killed automatically!
```

## Important Notes

- Child processes are killed when the function exits, regardless of how it exits
- Multiple `defer` statements execute in LIFO (Last-In-First-Out) order
- Killing a process sends SIGTERM by default (can be changed with `child.kill('SIGKILL')`)
- Always handle process errors with `or` blocks to prevent unhandled promise rejections
