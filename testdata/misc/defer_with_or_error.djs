function openDb(name) {
  if (name == 'fail') {
    throw 'Cannot open database'
  }
  console.log('Database opened:', name)
  return { name: name, connected: true }
}

function closeDb(db) {
  console.log('Closing database:', db.name)
}

function testDeferWithOrError() {
  let db = openDb('fail') or |error| {
    console.log('Error occurred:', error)
    return
  }
  
  defer closeDb(db)
  
  console.log('This should not be printed')
}

testDeferWithOrError()
