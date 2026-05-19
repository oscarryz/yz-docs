#example

Fat - http-server

https://gitlab.com/fatscript/fry/-/blob/main/sample/http_server.fat

```js
Chunk: fat.type.Chunk
http:  fat.http


routes: [
  Route(
    '/hello*',
    get: HttpResponse(body: 'world')
  ),
  Route(
    '/json',
    post #(HttpRequest, HttpResponse) {
        HttpResponse(
          body: { "message": "hello" }
        )
    }
  ),
  Route(
    '/to-json',
    get: {
      request HttpRequest
      HttpResponse(
        body: request.params
      )
    }
  ),
  Route(
     '/binary',
     get: {
      _ HttpRequest
      HttpResponse(
        body: Chunk([1, 2, 3, 4, 5, 6])
      )
    }
  )
]
http.listen(8080, routes)

```