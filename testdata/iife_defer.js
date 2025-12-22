function openDb() {
  console.log('opening database')
  return {
    close: function() {
      console.log('closing database')
    }
  }
}

(function main () {
  let db = openDb()
  defer db.close()
})()
