function getValue(shouldFail) {
  if (shouldFail) {
    throw 'Value error'
  }
  return 42
}

function testOrInExpression() {
  getValue(true) or |err| {
    console.log('Caught error:', err)
  }
  
  console.log('Continued execution')
}

testOrInExpression()
