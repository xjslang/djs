function processWithCleanup() {
  let counter = 0
  let multiplier = 2
  
  defer {
    console.log('final counter value:', counter * multiplier)
  }
  
  counter = 5
  console.log('processing with counter:', counter)
  
  multiplier = 3
  console.log('updated multiplier:', multiplier)
}

processWithCleanup()