function riskyOperation() {
  let resource = 'critical_resource'
  
  defer {
    console.log('cleaning up:', resource)
  }
  
  console.log('start operation')
  console.log('attempting risky task')
  console.log('operation completed')
}

riskyOperation()