# DJS - Defers for JavaScripts

This is a proof of concept for the [XJS](https://github.com/xjslang/xjs) package, a library capable of enriching JavaScript or even altering its default behavior. In this case, we've incorporated the [defers](https://go.dev/tour/flowcontrol/12) structure, popularized by the Go language.

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

## Install

Clone and install the repo:

```bash
git clone <this repo>
go install
```

And now you can use the `djs` command line:

```bash
djs <file.djs>         # compiles DJS to standard JS
djs -format <file.djs> # format DJS
```

## Example
The following command:
```bash
djs -stdin <<< "function foo() { defer closeDb() }"
```

outputs:
```js
function foo() { let __defers_123__ = [];
try {{__defers_123__.unshift(() => closeDb());}}
finally {for (let defer of __defers_123__)
{ try { defer() } catch (e) { console.error(e) }}}}
```

or formatted:
```js
function foo() {
  let __defers_123__ = [];
  try {
    {
      __defers_123__.unshift(() => closeDb());
    }
  } finally {
    for (let defer of __defers_123__) {
      try {
        // defers are called at the end of the block
        defer();
      } catch (e) {
        // If a "defer" fails, it displays an error and
        // does not prevent the other "defers" from executing.
        // This ensures that all defers are called.
        console.error(e);
      }
    }
  }
}
```

## How does it work?

The XJS library exports `NewBuilder` and `NewPrinter` functions to create custom parsers and printers.

```go
// here we are creating a parser builder
b := xjs.NewBuilder().
  Install(djsPlugin).
  Install(awaitPlugin).
  Install(strictEqualityPlugin).
  Install(anotherPlugin)
  // ...

p := b.Build(sourceCode) // builds the parser
result, err := js.ParseProgram(p) // returns the AST
```

Finally, we can transform the AST (Abstract Syntax Tree) to JavaScript code with our custom printer. For example:

```go
pr := xjs.NewPrinter()
pr.UsePrinter(compiler) // "instructs" the printer how to print our custom nodes
pr.Print(result)
jsCode, err := pr.Output() // Returns JS
```

Very simple! If you want to learn more, just take a look at the [djs](./djs) folder.
