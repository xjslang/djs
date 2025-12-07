function complexCleanup() {
  let step = 1
  
  defer {
    console.log('step 5: final cleanup')
  }
  
  console.log('step 1: initializing')
  step = 2
  
  defer {
    console.log('step 4: secondary cleanup, step was:', step)
  }
  
  console.log('step 2: processing')
  step = 3
  
  defer {
    console.log('step 3: immediate cleanup, step was:', step)
  }
  
  console.log('main work completed')
}

complexCleanup()