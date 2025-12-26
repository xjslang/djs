let fs = require('fs').promises
let os = require('os')
let path = require('path')

(function main() {
  // Create a temporary directory
  let tmpDir = fs.mkdtemp(path.join(os.tmpdir(), 'djs-example-')) or |err| {
    console.error('Failed to create temp directory:', err.message)
    return
  }
  defer fs.rmdir(tmpDir, { recursive: true })

  console.log('Created temporary directory:', tmpDir)

  // Create multiple temporary files
  let file1 = path.join(tmpDir, 'data1.json')
  let file2 = path.join(tmpDir, 'data2.txt')
  let file3 = path.join(tmpDir, 'output.log')

  // Write to first file
  let data1 = { name: 'Alice', age: 30, city: 'NYC' }
  fs.writeFile(file1, JSON.stringify(data1, null, 2)) or |err| {
    console.error('Failed to write file1:', err.message)
    return
  }
  console.log('✓ Created:', file1)

  // Write to second file
  fs.writeFile(file2, 'This is a temporary text file\nLine 2\nLine 3') or |err| {
    console.error('Failed to write file2:', err.message)
    return
  }
  console.log('✓ Created:', file2)

  // Read and process the JSON file
  let content = fs.readFile(file1, 'utf8') or |err| {
    console.error('Failed to read file1:', err.message)
    return
  }
  
  let parsedData = JSON.parse(content)
  console.log('✓ Read data:', parsedData)

  // Create a log entry
  let logEntry = `[${new Date().toISOString()}] Processed ${parsedData.name}'s data\n`
  fs.writeFile(file3, logEntry) or |err| {
    console.error('Failed to write log:', err.message)
    return
  }
  console.log('✓ Created log:', file3)

  // List all files in the temp directory
  let files = fs.readdir(tmpDir) or |err| {
    console.error('Failed to list directory:', err.message)
    return
  }
  console.log('✓ Files in temp directory:', files)

  console.log('\nAll operations completed successfully!')
  console.log('Temporary directory will be cleaned up automatically.')
})()
