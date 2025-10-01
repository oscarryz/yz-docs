#example

https://www.roc-lang.org/#examples

```js

main: {
	store_email(path.from_string('url.txt')).map_err(handle_error)
}
// store_email #(path Path, result.Result(V,E))
store_email: {
	path Path
	early_return: {e Err; e}
	url: file.read_utf8(path).or early_return
	user: http.get url(json.codec).or early_return
	dest: path.from_string('`user.name`.txt').or early_return
	_: file.write_utf8 des user.email.or early_return
	print('Wrote emai to `path.display(dest))`'
}

handler_error: {
	err Err
	when_eq err [
		{HttpErr} : {print('Error fetching URL `err.cause`'})
		{FileReaderErr} : {print('Error reading from `path.display(err.cause`))'}
		{FileWriter} : {print('Error writing to `path.display(err.cause`))'}
	]
}
```

Will inform [return, break, continue, exit](../../Questions/return,%20break,%20continue,%20exit.md)

```js
main: {
  store_email(path.from_string("url.txt")).map_err(handle_error)
}
store_email #(Path, Result(String,String)) = {
  path Path
  fs.read(path).and_then {
    url String
    http.get(url, json.UTF8)
  }.and_then {
    json JSONObject
    p: path.from_string("`json.name`.txt")
    fs.write(p, json.email)
    print("Wrote email to: `p`")
  }
}
handle_error: {
  err Error(String)
  when_eq info(err).type [
    {HttpErr} : {print('Error fetching URL: `err`'})
    {FileReaderErr} : {print('Error reading from path: `err`'})
    {FileWriter} : {print('Error writing file: `err`'})
  ]
}
```

Example using the `>>=` function

```js
store_email #(Path, Result(#(),String)) = {
  path Path
  fs.read(path) >>= { url String
    http.get(url, json.UTF8)
  }  >>= { json JSONObject
    p: path.from_string("`json.name`.txt")
    fs.write_to(p, json.email)
  } >>= {
    p Path
    print("Wrote email to: `p`")
    Ok("")
  }
}
```
