# Locks and Semaphores Example

This example demonstrates how DJS's `defer` statement simplifies resource management in concurrent programming scenarios, specifically with locks and semaphores.

## Features Demonstrated

- **`defer` for automatic lock release**: Ensures locks are always released, even on errors
- **Mutual exclusion with locks**: Protects critical sections from concurrent access
- **Concurrency limiting with semaphores**: Controls the number of concurrent operations
- **Nested locks**: Multiple locks with automatic LIFO cleanup
- **`or` for error handling**: Clean error handling without try-catch blocks

## What the Example Does

### Test 1: Sequential Transfers with Locks
Demonstrates protecting a shared resource (bank account balance) using locks:
- Acquires a lock before modifying the balance
- Prevents race conditions in concurrent access
- Automatically releases the lock using `defer`

### Test 2: Concurrent Jobs with Semaphores
Shows how to limit concurrent operations:
- Creates 4 jobs but only allows 2 to run concurrently
- Uses a semaphore with capacity of 2
- Automatically releases semaphore slots using `defer`

### Test 3: Nested Locks
Demonstrates complex operations requiring multiple locks:
- Acquires two locks in sequence
- Both locks are automatically released in reverse order (LIFO)
- Shows how `defer` handles nested resource cleanup

## Key DJS Features

### Automatic Lock Release
```djs
let lock = await acquireLock('resource-123') or |err| {
  console.error('Failed to acquire lock:', err.message)
  return
}
defer lock.release()

// Critical section
await criticalOperation()
// Lock released automatically, even if operation fails
```

### Semaphore Control
```djs
let sem = await acquireSemaphore('worker-pool', 2) or |err| {
  console.error('Failed to acquire semaphore:', err.message)
  return
}
defer sem.release()

await processJob()
// Semaphore slot released automatically
```

### Nested Resource Management
```djs
let lockA = await acquireLock('resource-a') or |err| { return }
defer lockA.release()

let lockB = await acquireLock('resource-b') or |err| { return }
defer lockB.release()

// Both locks released automatically in reverse order (B then A)
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
=== DJS Locks and Semaphores Example ===

Initial balance: $1000

--- Test 1: Sequential Transfers ---

💰 Starting transfer: $200 from Account A to Account B
🔒 Lock "account-lock" acquired
   Current balance: $1000
   ✅ Transfer completed. New balance: $800
🔓 Lock "account-lock" released

💰 Starting transfer: $300 from Account C to Account D
🔒 Lock "account-lock" acquired
   Current balance: $800
   ✅ Transfer completed. New balance: $500
🔓 Lock "account-lock" released

💰 Starting transfer: $600 from Account E to Account F
🔒 Lock "account-lock" acquired
   Current balance: $500
   ❌ Insufficient funds
🔓 Lock "account-lock" released


--- Test 2: Concurrent Jobs with Semaphore (max 2) ---

🔧 Job 1: Waiting for semaphore...
🔧 Job 2: Waiting for semaphore...
🔧 Job 3: Waiting for semaphore...
🔧 Job 4: Waiting for semaphore...
🎫 Semaphore "worker-pool" acquired (1/2)
   Job 1: Processing...
🎫 Semaphore "worker-pool" acquired (2/2)
   Job 2: Processing...
   Job 1: ✅ Completed
🎟️  Semaphore "worker-pool" released (1/2)
🎫 Semaphore "worker-pool" acquired from queue (2/2)
   Job 3: Processing...
   Job 2: ✅ Completed
🎟️  Semaphore "worker-pool" released (1/2)
🎫 Semaphore "worker-pool" acquired from queue (2/2)
   Job 4: Processing...
   Job 3: ✅ Completed
🎟️  Semaphore "worker-pool" released (1/2)
   Job 4: ✅ Completed
🎟️  Semaphore "worker-pool" released (0/2)


--- Test 3: Nested Locks ---

🔄 Complex operation on resource X
🔒 Lock "resource-X-a" acquired
   Step 1: Lock A acquired
🔒 Lock "resource-X-b" acquired
   Step 2: Lock B acquired
   Step 3: ✅ Complex operation completed
🔓 Lock "resource-X-b" released
🔓 Lock "resource-X-a" released


=== All tests completed! ===
```

## Comparison with Standard JavaScript

### Without DJS (manual cleanup)
```javascript
const lock = await acquireLock('resource-123');
try {
  await criticalOperation();
} finally {
  lock.release(); // Must remember to release!
}
```

### With DJS (automatic cleanup)
```djs
let lock = await acquireLock('resource-123') or |err| {
  console.error('Failed:', err.message)
  return
}
defer lock.release()

await criticalOperation()
// Lock released automatically!
```

## Benefits

1. **No forgotten releases**: `defer` guarantees cleanup happens
2. **Cleaner code**: No nested try-finally blocks
3. **LIFO order**: Multiple defers execute in reverse order (like Go)
4. **Error safety**: Resources are released even if errors occur
5. **Easier to reason about**: Acquisition and release are visually close

## Implementation Notes

- The lock and semaphore implementations in `lock.js` are simplified for demonstration
- In production, consider using libraries like `async-mutex` or `semaphore`
- The example uses `defer` to ensure proper resource cleanup in all scenarios
- Semaphores limit concurrency while locks provide mutual exclusion
