function openDb(path) {
  if (path == 'fail') {
    throw 'Database connection failed'
  }
  return { connected: true, path: path }
}

function testSuccessCase() {
  let db = openDb('mydata.db') or |err| {
    console.log('Error:', err)
    return
  }
  console.log('Connected to:', db.path)
}

testSuccessCase()
