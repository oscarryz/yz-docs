
### High level requirements 

- Managed through `yz` tool e.g. `yz init|add|remove|update` 
- Stored in the high level `dependencies` object
- The `dependencies` object will specify a name the app will use and a http uri to get the resource
- If the resource is an object containing a `dependencies` object it will download it too
- A `.lock` file will be generated to validate on subsequent runs / installations 
- Resources will be cached and stored somewhere in the src path (defaults to `./external_libraries` )
- Should be secure

#### Example

```js
yz init hello_world
// Creates hello_world.js 
main: {
	print 'Hi'
}
```

If we then add a `dependencies` object
```js
dependencies:[
	{"http/body_parser" "https://raw.githubusercontent.com/foo/body-parser/v0.1.3/body_parser.yz"}
	{path:"util/more/thing" uri:"https://raw.github.com/foo/stuff/stuff/yz"}
]
```

And run `yz install` Yz will fetch the file and if it has a `dependencies` object itself, it will fetch their dependencies and will install in `./external_libraries`

Only after installing the object can be used under the specified path

```js
//hello_world.yz
main: {
    xyz: http.body_parser.something()
}
```


### Resources

[Dependency management in Deno](https://medium.com/deno-the-complete-reference/dependency-management-in-deno-48f1c91ad84d)
[deno/modules](https://docs.deno.com/runtime/manual/basics/modules/)
[How should I build a package manager?](https://www.reddit.com/r/ProgrammingLanguages/comments/m3wvsv/how_should_i_build_a_package_manager/)


