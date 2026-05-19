#example

[Writing Assertive Code With Elixir](https://dashbit.co/blog/writing-assertive-code-with-elixir)

First attempt
```js
`
!:[Example]
example: {
	value: get_token("foo=bar&token=value&bar=baz")
	// value is Option.Ok("value")
	not_found: get_token("foo=bar&token=val=ue&bar=baz")
	// not_found is Option.None()
}
`
get_token : {
	s String
	s.split("&")
	.map({e String; e.split("=")})
	.find({pair [String]; pair[0] == "token" && { pair.len() == 2}})
	.map({ pair [String]; pair[1]})
}
value Option(String) = get_token("foo=bar&token=value&bar=baz")
```