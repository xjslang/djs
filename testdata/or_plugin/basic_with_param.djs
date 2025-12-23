function openDb(path) {
  if (path == 'fail') {
    throw 'Database connection failed'
  }
  return { connected: true, path: path }
}

function testWithErrorParam() {
  let db = openDb('fail') or |err| {
    console.log('Error:', err)
    return
  }
  console.log('Database connected')
}

testWithErrorParam()
