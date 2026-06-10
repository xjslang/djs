// comments are preserved
function init() {
  let db = openDb()
  defer closeDb(db) // ensure closing db properly

  let file = openFile('myfile.txt')
  defer {
    console.log('closing file')
    closeFile(file) // ensure closing file properly
  }

  // db and file operations ...
}
