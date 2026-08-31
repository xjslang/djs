# DJS (Defers for JavaScripts)

JS with `defer(s)` built of top of [XJS](https://github.com/xjslang/xjs).

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

<details>
  <summary>Is transpiled to:</summary>

```js
function init() {
  let __defers_2a78ec64__ = [];
  try {
    {
      let db = openDb();
      __defers_2a78ec64__.unshift(() => {
        closeDb(db);
      });
      let file = openFile("myfile.txt");
      __defers_2a78ec64__.unshift(() => {
        {
          console.log("closing file");
          closeFile(file);
        }
      });
    }
  } finally {
    for (let defer of __defers_2a78ec64__) {
      try {
        defer();
      } catch (e) {
        console.error(e);
      }
    }
  }
}
```
</details>

## Why?

For language ergonomics. Writing a `try/finally` block for every resource we want to release is cumbersome and hinders the natural reading of the code.

> [!NOTE]
> The latest versions of JS already have the `using` declaration, but it is still cumbersome and only available in latest versions.

## Requirements & Install

You'll need [Go](https://go.dev/doc/install). Once installed, simply do this:

```bash
git clone <this repo>
go install ./cmd/djs
```

Examples of use:

```bash
djs <file.djs>         # compile DJS to JS
djs <file.djs> | node  # compile and run
djs -format <file.djs> # format DJS
djs -check <file.djs>  # check for errors
```
