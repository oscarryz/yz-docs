#example

https://twitter.com/pomarlang/status/1763877680280187098/photo/1

Using [Type variants](Type%20variants.md) + [Conditional Bocs](Conditional%20Bocs.md)

```js
HttStatus: {
    Ok(),
    ClientError(error_message String)    
}

HttpResult : { 
    status HttpStatus
}
make_http_request #(HttpStatus) {
    HttpStatus.ClientError("Invalid Request")
}

match_with_iflets: {
    // With function call 
    response : make_http_request()
    match response {
        Ok => print("Ok")
    }, {
        ClientError => print(response.error_message)
    }
    
    // With field access
    result : HttpResult(make_http_request())
    match result.status {
        Ok => print("Ok")
    }, {
        ClientError => print(response.error_message)
    }
}
```