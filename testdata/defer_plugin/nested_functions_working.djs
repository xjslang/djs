function outerFunction() {
  console.log('outer function start')
  
  defer {
    console.log('outer function cleanup')
  }
  
  function innerFunction() {
    console.log('inner function start')
    
    defer {
      console.log('inner function cleanup')
    }
    
    console.log('inner function end')
  }
  
  innerFunction()
  console.log('outer function end')
}

outerFunction()