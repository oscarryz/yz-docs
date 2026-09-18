---
type: example
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#example

From https://gowebexamples.com/static-files/

```go

    fs: http.file_server(http.dir('assets/'))
	    http.handle('/static/'(http.strip_prefix('/static/'), fs))

    http.listen_and_serve(':8080')
```


```js
data: io.read('/tmp/something.txt')

data.each({ byte Int
    match {
        byte == '0xCAFEBABE' => init()
    }, {
        byte == '0xDEADBEEF' => finish()
    }
})

```