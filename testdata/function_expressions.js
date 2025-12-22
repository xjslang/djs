// Test function expression assigned to variable
let assigned = function() {
  console.log('assigned start')
  defer {
    console.log('assigned defer')
  }
  console.log('assigned middle')
}
assigned()
