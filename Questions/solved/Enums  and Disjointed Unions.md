

#answered  [Type variants (similar to SumTypes)](Type%20variants%20(similar%20to%20SumTypes).md)



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


