function openDb(name) {
  if (name == 'fail') {
    throw 'Database error: ' + name
  }
  console.log('Opened:', name)
  return { name: name }
}

function closeResource(name) {
  console.log('Closed:', name)
}

function testMultipleDeferOr() {
  let db1 = openDb('db1') or |err| {
    console.log('Failed db1:', err)
    return
  }
  
  defer closeResource(db1.name)
  
  let db2 = openDb('db2') or |err| {
    console.log('Failed db2:', err)
    return
  }
  
  defer closeResource(db2.name)
  
  console.log('Both databases ready')
}

testMultipleDeferOr()
