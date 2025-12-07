function outerFunction() {
  console.log('outer function start')
  
  defer {
    console.log('outer function cleanup')
  }
  
  console.log('doing some work')
  console.log('outer function end')
}

function innerFunction() {
  console.log('inner function start')
  
  defer {
    console.log('inner function cleanup')
  }
  
  console.log('inner function end')
}

outerFunction()
innerFunction()