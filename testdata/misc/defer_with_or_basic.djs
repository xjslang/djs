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

function testDeferWithOr() {
  let db = openDb('mydata.db') or |error| {
    console.log('Something was wrong:', error)
    return
  }
  
  defer closeDb(db)
  
  console.log('Working with database:', db.name)
}

testDeferWithOr()
