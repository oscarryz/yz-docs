#example

[Scrape](https://www.reddit.com/r/ProgrammingLanguages/s/CAHHEo0Xil)


```js
html: golang.org.x.net.html
http: std.net.http

main: {
  resp: http.get("https://en.wikipedia.org/wiki/Go").or {
     err Err
     panic(err)
  }
  defer {resp.body.close()}{
     doc : html.parse(resp.body).or {
       err Err
       panic(err)
     }

     title : doc.find("title").first_child.data
     print(title)
  }
}
```
