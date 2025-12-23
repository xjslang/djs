function saveData(data) {
  if (data == 'bad') {
    throw 'Cannot save data'
  }
  console.log('Data saved:', data)
  return true
}

function cleanup() {
  console.log('Starting cleanup')
  
  let result = saveData('backup') or |err| {
    console.log('Backup failed:', err)
  }
  
  console.log('Cleanup finished')
}

function testOrInsideDefer() {
  console.log('Begin operation')
  
  defer cleanup()
  
  console.log('Main operation')
}

testOrInsideDefer()
