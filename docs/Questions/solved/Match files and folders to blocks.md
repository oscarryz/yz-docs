In project structure I basically equated files to blocks, let's make sure that's what we want
See [Code organization](Code%20organization.md)


e.g. in a `/main.yz` file
```js
// main.yz
print foo.name
print foo.bar.name
```

This:

```
./foo.yz
./foo/bar.yz
```
```js
// foo.yz
name String
//./foo/bar.yz
name String
```
Would be the same as this:
```js
//./foo.yz
name String
bar : {
   name String
}
```

Notes, you cannot add the file `foo/bar.yz` anymore as the `bar` block is already defined, you can add a sibling `./foo/baz.yz` though
#answered
