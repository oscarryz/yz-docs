#example

https://harelang.org/

```js
main: {
	greetings: [
		'Hello, world!'
		'¡Hola Mundo!'
		"Γειά σου Κόσμε!"
		"Привіт, світ!"
		"こんにちは世界！"
	]
	greetings.each({
	    greet String
		print(greet)
	}
}
```

https://harelang.org/tutorials/introduction/

```js
// fs.create(path, perms) can an error
path: '/tmp/xyz.txt'
file: fs.create(path, fs.ErrorHandler {[
    {f File; f == fs.noaccess}: {"Error opening: {path}. Access denied"}
    {f File; f == fs.error   }: {"Error opening: {path}. {info(file)}"}
]})
// Updated
path: '/tmp/xyz.txt'
file: fs.create(path).or {
   e FileError
   panic("Error opening: `path`. `e`)
}
```


Blocks that errors can provide a handling error block that will be invoked if an error is generated

```js
fs: {
    create: {
        path String
        // do something
        error_handler ErrorHandler
        file File = ....
        file.noaccess ? {error_handler.handle(file)}
    }
    ErrorHandler: {
       scenarios [{File,Boolean}:{}]
       handle: { file File
           scenarios.keys().for_each({k {File,Boolean}
                k(file).?({
                    scenarios[k]();
                })
           })
       }
    }
}
// Updated
fs: {
   create {path String; Result}
   create = {
		path String
		file File = ...
		ok(file)
		// if error
		error("cannot open file")
	}
}
```

```js
keep: true
while ({keep},{
    line : bufio.scanline(os.stdin)
    what: when([
        {line == u8}:{line}
        {line == io.EOF}: { keep = false}
    ])
    append(lines, strings.fromutf8(what))
})
```
