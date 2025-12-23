function openDb(path) {
  if (path == 'fail') {
    throw 'Connection error'
  }
  return { connected: true, path: path }
}

function testOrWithFallbackValue() {
  let db = openDb('fail') or |err| {
    console.log('Using fallback:', err)
    db = { connected: false, path: 'fallback.db' }
  }
  
  if (db.connected) {
    console.log('DB status: connected')
  } else {
    console.log('DB status: fallback')
  }
  console.log('DB path:', db.path)
}

testOrWithFallbackValue()
