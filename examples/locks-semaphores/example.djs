let lockModule = require('./lock.js')
let acquireLock = lockModule.acquireLock
let acquireSemaphore = lockModule.acquireSemaphore

// Simulate a shared resource (bank account balance)
let balance = 1000

// Helper to simulate async work
function sleep(ms) {
  return new Promise(function(resolve) {
    setTimeout(resolve, ms)
  })
}

// Example 1: Using locks to protect critical sections
async function transferMoney(from, to, amount) {
  console.log(`\n💰 Starting transfer: $${amount} from ${from} to ${to}`)
  
  // Acquire lock for the account
  let lock = await acquireLock('account-lock') or |err| {
    console.error('❌ Failed to acquire lock:', err.message)
    return
  }
  defer lock.release()

  // Critical section - reading and modifying shared resource
  let currentBalance = balance
  console.log(`   Current balance: $${currentBalance}`)
  
  if (currentBalance < amount) {
    console.log(`   ❌ Insufficient funds`)
    return
  }

  // Simulate processing time
  await sleep(100)
  
  balance = currentBalance - amount
  console.log(`   ✅ Transfer completed. New balance: $${balance}`)
}

// Example 2: Using semaphores to limit concurrent operations
async function processJob(jobId, semaphore) {
  console.log(`\n🔧 Job ${jobId}: Waiting for semaphore...`)
  
  let sem = await acquireSemaphore(semaphore, 2) or |err| {
    console.error(`❌ Job ${jobId}: Failed to acquire semaphore:`, err.message)
    return
  }
  defer sem.release()

  console.log(`   Job ${jobId}: Processing...`)
  await sleep(200)
  console.log(`   Job ${jobId}: ✅ Completed`)
}

// Example 3: Nested locks with automatic cleanup
async function complexOperation(resourceId) {
  console.log(`\n🔄 Complex operation on resource ${resourceId}`)
  
  let lockA = await acquireLock(`resource-${resourceId}-a`) or |err| {
    console.error('❌ Failed to acquire lock A:', err.message)
    return
  }
  defer lockA.release()

  console.log(`   Step 1: Lock A acquired`)
  await sleep(50)

  let lockB = await acquireLock(`resource-${resourceId}-b`) or |err| {
    console.error('❌ Failed to acquire lock B:', err.message)
    return
  }
  defer lockB.release()

  console.log(`   Step 2: Lock B acquired`)
  await sleep(50)
  
  console.log(`   Step 3: ✅ Complex operation completed`)
  // Both locks will be released automatically in reverse order (LIFO)
}

// Main execution
(async function main() {
  console.log('=== DJS Locks and Semaphores Example ===\n')
  console.log('Initial balance: $' + balance)

  // Test 1: Sequential transfers with locks
  console.log('\n--- Test 1: Sequential Transfers ---')
  await transferMoney('Account A', 'Account B', 200)
  await transferMoney('Account C', 'Account D', 300)
  await transferMoney('Account E', 'Account F', 600) // This should fail

  // Test 2: Concurrent jobs with semaphore (max 2 concurrent)
  console.log('\n\n--- Test 2: Concurrent Jobs with Semaphore (max 2) ---')
  let jobs = Promise.all([
    processJob(1, 'worker-pool'),
    processJob(2, 'worker-pool'),
    processJob(3, 'worker-pool'),
    processJob(4, 'worker-pool')
  ])
  await jobs

  // Test 3: Nested locks
  console.log('\n\n--- Test 3: Nested Locks ---')
  await complexOperation('X')

  console.log('\n\n=== All tests completed! ===')
})()
