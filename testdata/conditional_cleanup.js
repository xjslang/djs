function fileProcessor() {
  let filename = 'data.txt'
  let fileHandle = null
  
  defer {
    if (fileHandle) {
      console.log('closing file:', filename)
    }
  }
  
  console.log('opening file:', filename)
  fileHandle = { name: filename, open: true }
  
  if (fileHandle.open) {
    console.log('processing file:', filename)
    console.log('file processed successfully')
  }
}

fileProcessor()