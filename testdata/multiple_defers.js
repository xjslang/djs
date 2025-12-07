function resourceManager() {
  console.log('acquiring resource 1')
  defer {
    console.log('releasing resource 1')
  }
  
  console.log('acquiring resource 2')
  defer {
    console.log('releasing resource 2')
  }
  
  console.log('acquiring resource 3')
  defer {
    console.log('releasing resource 3')
  }
  
  console.log('doing work')
}

resourceManager()