let childProcess = require('child_process')
let spawn = childProcess.spawn
let path = require('path')

// Helper to wait for a process to complete
function waitForProcess(child) {
  return new Promise(function(resolve, reject) {
    child.on('exit', function(code) {
      if (code == 0) {
        resolve(code)
      } else {
        reject(new Error(`Process exited with code ${code}`))
      }
    })
    child.on('error', function(err) {
      reject(err)
    })
  })
}

// Helper to simulate async work
function sleep(ms) {
  return new Promise(function(resolve) {
    setTimeout(resolve, ms)
  })
}

// Example 1: Simple command execution with automatic cleanup
async function runSimpleCommand() {
  console.log('\n📋 Example 1: Running simple command (ls)')
  
  let child = spawn('ls', ['-la'])
  defer child.kill()
  
  child.stdout.on('data', function(data) {
    console.log(`   ${data.toString().trim()}`)
  })
  
  await waitForProcess(child) or |err| {
    console.error('   ❌ Command failed:', err.message)
    return
  }
  
  console.log('   ✅ Command completed')
}

// Example 2: Long-running process with automatic termination
async function runLongProcess() {
  console.log('\n⏳ Example 2: Long-running process (will be terminated early)')
  
  let workerPath = path.join(__dirname, 'worker.js')
  let child = spawn('node', [workerPath, '10000', '500', 'LongWorker'])
  defer child.kill()
  
  child.stdout.on('data', function(data) {
    console.log(`   ${data.toString().trim()}`)
  })
  
  // Let it run for only 2 seconds, then exit
  await sleep(2000)
  
  console.log('   ⚠️  Parent function ending - child process will be killed automatically')
}

// Example 3: Multiple concurrent processes
async function runMultipleProcesses() {
  console.log('\n🔀 Example 3: Multiple concurrent processes')
  
  let workerPath = path.join(__dirname, 'worker.js')
  
  let child1 = spawn('node', [workerPath, '2000', '500', 'Worker-1'])
  defer child1.kill()
  
  child1.stdout.on('data', function(data) {
    console.log(`   ${data.toString().trim()}`)
  })
  
  let child2 = spawn('node', [workerPath, '2000', '500', 'Worker-2'])
  defer child2.kill()
  
  child2.stdout.on('data', function(data) {
    console.log(`   ${data.toString().trim()}`)
  })
  
  let child3 = spawn('node', [workerPath, '2000', '500', 'Worker-3'])
  defer child3.kill()
  
  child3.stdout.on('data', function(data) {
    console.log(`   ${data.toString().trim()}`)
  })
  
  // Wait for all to complete
  await Promise.all([
    waitForProcess(child1),
    waitForProcess(child2),
    waitForProcess(child3)
  ]) or |err| {
    console.error('   ❌ One or more processes failed:', err.message)
    return
  }
  
  console.log('   ✅ All processes completed successfully')
}

// Example 4: Process with timeout
async function runWithTimeout() {
  console.log('\n⏱️  Example 4: Process with timeout')
  
  let workerPath = path.join(__dirname, 'worker.js')
  let child = spawn('node', [workerPath, '10000', '500', 'TimeoutWorker'])
  defer child.kill()
  
  child.stdout.on('data', function(data) {
    console.log(`   ${data.toString().trim()}`)
  })
  
  // Set a timeout
  let timeoutReached = false
  let timeoutId = setTimeout(function() {
    timeoutReached = true
    console.log('   ⏰ Timeout reached! Killing process...')
    child.kill()
  }, 2500)
  defer clearTimeout(timeoutId)
  
  await waitForProcess(child) or |err| {
    if (timeoutReached) {
      console.log('   ✅ Process terminated due to timeout (as expected)')
      return
    }
    console.error('   ❌ Process failed:', err.message)
    return
  }
  
  console.log('   ✅ Process completed before timeout')
}

// Example 5: Processing command output
async function processOutput() {
  console.log('\n📊 Example 5: Processing command output')
  
  let child = spawn('echo', ['Hello from child process!'])
  defer child.kill()
  
  let output = ''
  child.stdout.on('data', function(data) {
    output = output + data.toString()
  })
  
  await waitForProcess(child) or |err| {
    console.error('   ❌ Command failed:', err.message)
    return
  }
  
  console.log(`   Received output: "${output.trim()}"`)
  console.log(`   Output length: ${output.trim().length} characters`)
  console.log('   ✅ Output processed successfully')
}

// Example 6: Nested processes with LIFO cleanup
async function nestedProcesses() {
  console.log('\n🪆 Example 6: Nested processes with LIFO cleanup')
  
  let workerPath = path.join(__dirname, 'worker.js')
  
  let outer = spawn('node', [workerPath, '4000', '1000', 'Outer'])
  defer outer.kill()
  console.log('   Started outer process')
  
  outer.stdout.on('data', function(data) {
    console.log(`   ${data.toString().trim()}`)
  })
  
  await sleep(500)
  
  let middle = spawn('node', [workerPath, '3000', '800', 'Middle'])
  defer middle.kill()
  console.log('   Started middle process')
  
  middle.stdout.on('data', function(data) {
    console.log(`   ${data.toString().trim()}`)
  })
  
  await sleep(500)
  
  let inner = spawn('node', [workerPath, '2000', '600', 'Inner'])
  defer inner.kill()
  console.log('   Started inner process')
  
  inner.stdout.on('data', function(data) {
    console.log(`   ${data.toString().trim()}`)
  })
  
  // Wait only for inner to complete
  await waitForProcess(inner) or |err| {
    console.error('   ❌ Inner process failed:', err.message)
  }
  
  console.log('   ⚠️  Function ending - processes will be killed in LIFO order: Inner → Middle → Outer')
  await sleep(500)
}

// Example 7: Error handling with failed command
async function handleFailedCommand() {
  console.log('\n❌ Example 7: Handling failed command')
  
  // ls on a non-existent directory will fail with exit code 1
  let child = spawn('ls', ['/nonexistent/directory/path'])
  defer child.kill()
  
  child.stderr.on('data', function(data) {
    console.log(`   Error output: ${data.toString().trim()}`)
  })
  
  await waitForProcess(child) or |err| {
    console.log(`   ✅ Caught error as expected: ${err.message}`)
    return
  } or |err| {
    console.log(`   ✅ Caught error as expected: ${err.message}`)
    return
  }
  
  console.log('   This should not appear')
}

// Main execution
(async function main() {
  console.log('=== DJS Child Processes Example ===')
  console.log('Demonstrates automatic cleanup of child processes using defer\n')

  await runSimpleCommand()
  await sleep(500)

  await runLongProcess()
  await sleep(1000)

  await runMultipleProcesses()
  await sleep(500)

  await runWithTimeout()
  await sleep(500)

  await processOutput()
  await sleep(500)

  await nestedProcesses()
  await sleep(500)

  await handleFailedCommand()

  console.log('\n=== All tests completed! ===')
  console.log('Notice how all child processes were cleaned up automatically.')
})()
