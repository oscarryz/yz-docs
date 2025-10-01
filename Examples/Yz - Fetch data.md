#example
```js
http: std.net.http
Result: std.result.Result

fetch_data #(url String, Result(String,Error)) = {
	http.get(url).or_else {
		Error("unable to retrieve data")
	}
}

main: {
	urls: [
		"https://api.example.com/data1"
		"https://api.example.com/data2"
	]
	results: urls.map(fetch_data)

	esults.each({ 
		r Result(String,Error)
		r.if_ok { 
		  value String
		  print("The result is `value`")
		}
	}
}
```