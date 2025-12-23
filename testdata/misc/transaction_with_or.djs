function beginTransaction(id) {
  if (id == 0) {
    throw 'Invalid transaction ID'
  }
  console.log('Transaction started:', id)
  return { id: id, active: true }
}

function rollback(tx) {
  console.log('Rolling back transaction:', tx.id)
}

function commit(tx) {
  console.log('Committing transaction:', tx.id)
}

function processTransaction(id) {
  let tx = beginTransaction(id) or |err| {
    console.log('Failed to start transaction:', err)
    return
  }
  
  defer {
    if (tx.active) {
      rollback(tx)
    }
  }
  
  console.log('Processing transaction:', tx.id)
  
  tx.active = false
  commit(tx)
}

processTransaction(123)
console.log('---')
processTransaction(0)
