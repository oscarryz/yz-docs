---
type: solved
generated: { by: "oscarryz", at: 2024-10-23T12:30:27-05:00 }
---
#solved  [Type variants](../../Features/Type%20variants.md)



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


