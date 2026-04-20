#open-question 

Fist inclination is to use https to get dependencies, following a similar approach to Deno 

The dependency description would be in a .yz class that follows the `Project` interface, defining dictionaries or arrays with the things the project needs (see [[Code organization]] for examples) and these dependencies are not in the source code e.g. 

```js
// some_file.yz
mix foo
mix bar.baz

foo()
bar

// project.yz
dependencies: [
	"foo": {url:"https://foo.com/v1"; version:"v1.2.3"]},
	"bar": ...
]
```

But there is a lot to design: 

- How to ensure safety (lock.yz ? separate `.lock` file)
- Compare different approaches
- What is desirable for Yz considering it might not grow to much or anything at all