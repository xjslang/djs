function mixedDefers() {
  let resource1 = {
    cleanup: function() { console.log('cleaning resource1') }
  }
  let resource2 = {
    cleanup: function() { console.log('cleaning resource2') }
  }
  
  console.log('starting operations')
  
  // Defer with single statement
  defer resource1.cleanup()
  
  console.log('middle operation')
  
  // Defer with block
  defer {
    resource2.cleanup()
    console.log('additional cleanup')
  }
  
  console.log('ending operations')
}

mixedDefers()