function openConnection(name) {
  if (name == 'fail') {
    throw 'Connection failed'
  }
  console.log('Connected:', name)
  return { name: name }
}

function closeConnection(conn) {
  console.log('Disconnected:', conn.name)
}

function acquireLock(lockName) {
  console.log('Lock acquired:', lockName)
  return { lock: lockName }
}

function releaseLock(lock) {
  console.log('Lock released:', lock.lock)
}

function complexOperation() {
  let conn = openConnection('main') or |err| {
    console.log('Connection error:', err)
    return
  }
  
  defer closeConnection(conn)
  
  let lock = acquireLock('resource1') or |err| {
    console.log('Lock error:', err)
    return
  }
  
  defer releaseLock(lock)
  
  console.log('Performing operation')
}

complexOperation()
