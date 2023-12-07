https://www.roc-lang.org/#examples

```javascript

main: {
	store_email(path.from_striing 'url.txt').or(handle_error)

}
// store_email {path Path result.Result}
store_email: {
	path Path
	early_return: {e Err; return e}
	url: file.read_utf8(path).or early_return 
	user: http.get url(json.codec).or early_return
	dest: path.from_string('$(user.name).txt').or early_return
	_: file.write_utf8 des user.email.or early_return 
	print 'Wrote emai to ${path.display(dest)}'
}
handler_error: {
	err Err
	when_eq err [
		{HttpErr} : {print 'Error fetching URL $(err.cause)'}
		{FileReaderErr} : {print 'Error reading from $(path.display(err.cause))'}
		{FileWriter} : {print 'Error writing to $(path.display(err.cause))'}
	]
}
```

Will inform [return, break, continue, exit](../../Questions/return,%20break,%20continue,%20exit.md)
