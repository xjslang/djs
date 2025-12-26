# Timers and Intervals Example

This example demonstrates how DJS's `defer` statement simplifies cleanup of timers and intervals, preventing memory leaks and ensuring proper resource management.

## Features Demonstrated

- **`defer` for automatic timer cleanup**: Ensures timers are cleared even if errors occur
- **Interval management**: Health checks, progress reporters, heartbeats
- **Timeout handling**: Automatic cleanup of scheduled timeouts
- **Nested timers**: Multiple timers with LIFO cleanup order
- **Error scenarios**: Proper cleanup when operations timeout or fail

## What the Example Does

### Test 1: Health Monitoring
- Sets up a health check that runs every second
- Uses `defer` to guarantee interval cleanup
- Runs for 3.5 seconds then automatically stops

### Test 2: Multiple Timers
- Creates multiple timeouts and an interval
- All are automatically cleared using `defer`
- Demonstrates that deferred timers never execute

### Test 3: Progress Reporter
- Shows download progress using an interval
- Updates every 10% until completion
- Interval is automatically cleared when done

### Test 4: Timeout Scenarios
- Tests operations with timeout constraints
- One completes successfully, another times out
- Timeouts are always cleaned up with `defer`

### Test 5: Service with Heartbeat
- Simulates a service that sends periodic heartbeats
- Heartbeat continues while service works
- Automatically stops when service completes

### Test 6: Nested Timers
- Creates three nested intervals
- All are cleaned up in reverse order (LIFO)
- Demonstrates proper cleanup sequence

## Key DJS Features

### Automatic Interval Cleanup
```djs
let intervalId = setInterval(function() {
  checkHealth()
}, 1000)
defer clearInterval(intervalId)

await longRunningOperation()
// Interval cleared automatically, even on errors
```

### Automatic Timeout Cleanup
```djs
let timeoutId = setTimeout(function() {
  console.log('Timeout!')
}, 5000)
defer clearTimeout(timeoutId)

await quickOperation()
// Timeout cleared automatically, won't execute
```

### LIFO Cleanup Order
```djs
let timer1 = setInterval(doWork1, 1000)
defer clearInterval(timer1)

let timer2 = setInterval(doWork2, 1000)
defer clearInterval(timer2)

// Cleanup order: timer2 first, then timer1
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
=== DJS Timers and Intervals Example ===
Demonstrates automatic cleanup of timers using defer

🏥 Starting health monitoring...
[2025-12-26T...] Health check #1 - Status: running
[2025-12-26T...] Health check #2 - Status: running
[2025-12-26T...] Health check #3 - Status: running
✅ Monitoring completed - interval will be cleared automatically

⏰ Setting up multiple timers...
   All timers set, waiting 2 seconds...
   Tick 1
   Tick 2
   Tick 3
   Tick 4
✅ Function ending - all timers will be cleaned up automatically

📥 Downloading package.zip...
   Progress: 10%
   Progress: 20%
   Progress: 30%
   Progress: 40%
   Progress: 50%
   Progress: 60%
   Progress: 70%
   Progress: 80%
   Progress: 90%
   Progress: 100%
✅ Download of package.zip completed!

⏱️  Running Quick task with 2000ms timeout...
   ✅ Quick task completed successfully

⏱️  Running Slow task with 1500ms timeout...
   ❌ Timeout! Slow task took too long
   Slow task was cancelled due to timeout

💓 Starting service with heartbeat...
   💓 Heartbeat 1
   Working on task 1...
   💓 Heartbeat 2
   Working on task 2...
   💓 Heartbeat 3
   Working on task 3...
   💓 Heartbeat 4
✅ Service completed - heartbeat stopped automatically

🪆 Testing nested timers with LIFO cleanup...
   Outer interval set
   Middle interval set
   Inner interval set
✅ All timers will be cleared in reverse order (LIFO): Inner → Middle → Outer

=== All tests completed! ===
Notice how all intervals and timers were cleaned up automatically.
```

## Comparison with Standard JavaScript

### Without DJS (manual cleanup)
```javascript
const intervalId = setInterval(() => checkHealth(), 1000);
try {
  await longRunningOperation();
} finally {
  clearInterval(intervalId); // Must remember to clear!
}
```

### With DJS (automatic cleanup)
```djs
let intervalId = setInterval(function() {
  checkHealth()
}, 1000)
defer clearInterval(intervalId)

await longRunningOperation()
// Interval cleared automatically!
```

## Common Use Cases

### Progress Tracking
Automatically stop progress updates when operation completes:
```djs
let progress = setInterval(updateProgress, 100)
defer clearInterval(progress)
await performOperation()
```

### Heartbeat/Keep-Alive
Automatically stop heartbeat when connection closes:
```djs
let heartbeat = setInterval(sendHeartbeat, 5000)
defer clearInterval(heartbeat)
await handleConnection()
```

### Timeout Protection
Automatically cleanup timeouts whether operation succeeds or fails:
```djs
let timeout = setTimeout(handleTimeout, 10000)
defer clearTimeout(timeout)
await criticalOperation()
```

## Benefits

1. **No memory leaks**: Timers are always cleaned up
2. **Cleaner code**: No nested try-finally blocks
3. **Guaranteed cleanup**: Even if errors occur
4. **LIFO order**: Multiple timers cleaned up in reverse order
5. **Easier to reason about**: Setup and cleanup are visually close

## Important Notes

- All timers are cleared when the function exits, regardless of how it exits (return, error, completion)
- Multiple `defer` statements execute in LIFO (Last-In-First-Out) order
- Deferred timer clears happen even if the timer hasn't fired yet
- This prevents common bugs where timers continue running after they should have stopped
