async function main() {
  console.log("before")
  await Promise.resolve(true)
  console.log("after")
}
main()
