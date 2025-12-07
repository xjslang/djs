function databaseTransaction() {
  let transaction = { id: 123, active: true }
  
  defer {
    console.log('rollback transaction:', transaction.id)
  }
  
  defer {
    console.log('close connection for transaction:', transaction.id)  
  }
  
  console.log('begin transaction:', transaction.id)
  console.log('execute query in transaction:', transaction.id)
  
  // Simulate successful completion
  transaction.active = false
  console.log('commit transaction:', transaction.id)
}

databaseTransaction()