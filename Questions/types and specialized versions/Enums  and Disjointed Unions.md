```js
HttpStatus : { 
  Ok()
  ClientError(s String)
}
HttpResult : {
  status HttpStatus
}
make_http_request: {
  HttpResult(ClientError("Invalid Request"))
}

r : make_http_request()
match r.status { 
  Ok: print("Ok")
  ClientError(msg): print("Error: `msg`") 
}

HttpStatus: {
  Ok:{}
  ClientError: {
    s String
  }
}

```


