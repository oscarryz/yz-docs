#example
Hola

```js

name: io.args[0]
file: files.open(name 'r').or {
    print('Unable to open file')
    exit!
}

bytes []Bytes
while { rc: file.read_into bytes ; rc != files.eof } {
    io.stdout.write bytes
}
file.close()

```


