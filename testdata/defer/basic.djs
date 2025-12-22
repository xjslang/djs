function openDb() {
  console.log('opening db')
  return {
    close: function () {
      console.log('closing db')
    }
  }
}

function foo() {
  let db = openDb()
  defer {
    db.close()
  }
}

foo()