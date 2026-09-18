---
type: example
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#example
Hola

```js

name: io.args[0]
file: files.open(name, 'r').or({
    print('Unable to open file')
    exit!
})

bytes: [Bytes]()
rc: file.read_into(bytes)
while({ rc != files.eof }, {
    io.stdout.write(bytes)
    rc = file.read_into(bytes)
})
file.close()

```


