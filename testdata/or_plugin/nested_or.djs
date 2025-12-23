function riskyOperation(shouldFail) {
  if (shouldFail) {
    throw 'Operation failed'
  }
  return 'success'
}

function testNestedOr() {
  let result1 = riskyOperation(true) or |err1| {
    console.log('First error:', err1)
    let result2 = riskyOperation(false) or |err2| {
      console.log('Second error:', err2)
      return
    }
    console.log('Fallback succeeded:', result2)
    result1 = 'fallback-' + result2
  }
  
  console.log('Final result:', result1)
}

testNestedOr()
