```js
// native binding — wraps the goslug Go library
`go_source: "vendor/slug_binding.go"`
Slugger: {
  slugify #(text String, String)
  truncate #(text String, max Int, String)
}

// macro — derive JSON serialization for Post
`JSON: { pretty: true }`
Post: {
  title String
  slug  String
  body  String
}
  
`project: {
   name: "blog_tool"
   dependencies: [
   goslug: { url: "https://github.com/example/goslug", ref: "v1.2.0" }
]
}`
main: {
  s:    Slugger()
  post: Post(
    title: "Hello, World!",
    slug:  s.slugify("Hello, World!"),
    body:  "My first post."
  )
  print(post.to_json())
}

// vendor/slug_binding.go
//yz:bind Slugger slugify #(text String, String)
func SluggerSlugify(text std.String) std.String { ... }

//yz:bind Slugger truncate #(text String, max Int, String)
func SluggerTruncate(text std.String, max std.Int) std.String { ... }
```