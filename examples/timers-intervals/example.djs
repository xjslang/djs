// Helper to simulate async work
function sleep(ms) {
  return new Promise(function(resolve) {
    setTimeout(resolve, ms)
  })
}

// Simulate a health check system
let healthCheckCount = 0
let systemStatus = 'running'

function checkHealth() {
  healthCheckCount++
  console.log(`[${new Date().toISOString()}] Health check #${healthCheckCount} - Status: ${systemStatus}`)
}

// Example 1: Auto-cleanup of intervals with defer
async function monitorSystemHealth() {
  console.log('\n🏥 Starting health monitoring...')
  
  let intervalId = setInterval(checkHealth, 1000)
  defer clearInterval(intervalId)

  // Simulate long-running operation
  await sleep(3500)
  
  console.log('✅ Monitoring completed - interval will be cleared automatically')
}

// Example 2: Multiple timers with automatic cleanup
async function multipleTimers() {
  console.log('\n⏰ Setting up multiple timers...')
  
  let timer1 = setTimeout(function() {
    console.log('   Timer 1: This should NOT appear (cleared by defer)')
  }, 5000)
  defer clearTimeout(timer1)

  let timer2 = setTimeout(function() {
    console.log('   Timer 2: This should NOT appear (cleared by defer)')
  }, 10000)
  defer clearTimeout(timer2)

  let counter = 0
  let intervalId = setInterval(function() {
    counter++
    console.log(`   Tick ${counter}`)
  }, 500)
  defer clearInterval(intervalId)

  console.log('   All timers set, waiting 2 seconds...')
  await sleep(2000)
  
  console.log('✅ Function ending - all timers will be cleaned up automatically')
}

// Example 3: Progress reporter with interval
async function downloadFile(filename, durationMs) {
  console.log(`\n📥 Downloading ${filename}...`)
  
  let progress = 0
  let progressInterval = setInterval(function() {
    progress = progress + 10
    if (progress <= 100) {
      console.log(`   Progress: ${progress}%`)
    }
  }, durationMs / 10)
  defer clearInterval(progressInterval)

  await sleep(durationMs)
  
  console.log(`✅ Download of ${filename} completed!`)
}

// Example 4: Timeout with error handling
async function operationWithTimeout(taskName, workMs, timeoutMs) {
  console.log(`\n⏱️  Running ${taskName} with ${timeoutMs}ms timeout...`)
  
  let timeoutId = setTimeout(function() {
    console.log(`   ❌ Timeout! ${taskName} took too long`)
    systemStatus = 'timeout'
  }, timeoutMs)
  defer clearTimeout(timeoutId)

  // Simulate work
  await sleep(workMs)
  
  if (systemStatus == 'timeout') {
    console.log(`   ${taskName} was cancelled due to timeout`)
    systemStatus = 'running'
    return
  }
  
  console.log(`   ✅ ${taskName} completed successfully`)
}

// Example 5: Heartbeat with cleanup on error
async function serviceWithHeartbeat() {
  console.log('\n💓 Starting service with heartbeat...')
  
  let heartbeatCount = 0
  let heartbeat = setInterval(function() {
    heartbeatCount++
    console.log(`   💓 Heartbeat ${heartbeatCount}`)
  }, 800)
  defer clearInterval(heartbeat)

  // Simulate service work
  for (let i = 1; i <= 3; i++) {
    console.log(`   Working on task ${i}...`)
    await sleep(1000)
  }
  
  console.log('✅ Service completed - heartbeat stopped automatically')
}

// Example 6: Nested timers with LIFO cleanup
async function nestedTimers() {
  console.log('\n🪆 Testing nested timers with LIFO cleanup...')
  
  let outerInterval = setInterval(function() {
    console.log('   [Outer] This should NOT appear')
  }, 2000)
  defer clearInterval(outerInterval)
  console.log('   Outer interval set')

  await sleep(500)

  let middleInterval = setInterval(function() {
    console.log('   [Middle] This should NOT appear')
  }, 1500)
  defer clearInterval(middleInterval)
  console.log('   Middle interval set')

  await sleep(500)

  let innerInterval = setInterval(function() {
    console.log('   [Inner] This should NOT appear')
  }, 1000)
  defer clearInterval(innerInterval)
  console.log('   Inner interval set')

  await sleep(800)
  
  console.log('✅ All timers will be cleared in reverse order (LIFO): Inner → Middle → Outer')
}

// Main execution
(async function main() {
  console.log('=== DJS Timers and Intervals Example ===')
  console.log('Demonstrates automatic cleanup of timers using defer\n')

  // Test 1: Health monitoring with interval
  await monitorSystemHealth()
  await sleep(1000)

  // Test 2: Multiple timers
  await multipleTimers()
  await sleep(1000)

  // Test 3: Progress reporter
  await downloadFile('package.zip', 2000)
  await sleep(500)

  // Test 4: Timeout scenarios
  await operationWithTimeout('Quick task', 1000, 2000) // Should complete
  await sleep(500)
  await operationWithTimeout('Slow task', 3000, 1500) // Should timeout
  await sleep(500)

  // Test 5: Service with heartbeat
  await serviceWithHeartbeat()
  await sleep(1000)

  // Test 6: Nested timers
  await nestedTimers()

  console.log('\n=== All tests completed! ===')
  console.log('Notice how all intervals and timeouts were cleaned up automatically.')
})()
