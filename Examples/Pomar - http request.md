#example

https://twitter.com/pomarlang/status/1763877680280187098/photo/1

TBD: enums, asymetric structural typing?
```js
http_status: {
	Ok: {
	}
	ClientError: {
		error_msg String
	}
}
HttpResult: {
	ok Option(Ok)
	error Option(ClientError)
}
make_http_request: {
	http_status.ClientError("Invalid request")
}
match_with_iflets: {
	response: make_http_request()
	match {
		response == http_status.Ok         => println('Ok')
	}, {
		response == http_status.ClientError => println(response.error_msg)
	}

	result: HttpResult(make_http_request())
	match {
		result.status == http_status.Ok          => println('Ok')
	}, {
		result.status == http_status.ClientError => println(response.error_msg)
	}
}
```