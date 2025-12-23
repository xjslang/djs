function openDb(path) {
  if (path == 'fail') {
    throw 'Database error: ' + path
  }
  return { connected: true, path: path }
}

function readConfig(path) {
  if (path == 'missing') {
    throw 'Config not found'
  }
  return { loaded: true, config: 'data' }
}

function testMultipleOr() {
  let db = openDb('fail') or |err| {
    console.log('DB Error:', err)
    return
  }
  
  let config = readConfig('missing') or {
    console.log('Config error (no param)')
    return
  }
  
  console.log('All loaded successfully')
}

testMultipleOr()
