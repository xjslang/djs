# DJS - Defers for JavaScripts

JavaScript with [defers](https://go.dev/tour/flowcontrol/12).

```js
function init() {
  let db = openDb()
  defer closeDb(db)

  let file = openFile('myfile.txt')
  defer {
    console.log('closing file')
    closeFile(file)
  }

  // db and file operations ...
}
```
