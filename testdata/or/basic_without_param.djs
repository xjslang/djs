function openDb(path) {
  if (path == 'fail') {
    throw 'Database connection failed'
  }
  return { connected: true, path: path }
}

function testWithoutErrorParam() {
  let db = openDb('fail') or {
    console.log('Cannot connect to database')
    return
  }
  console.log('Database connected')
}

testWithoutErrorParam()
