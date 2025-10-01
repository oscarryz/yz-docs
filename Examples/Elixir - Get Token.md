#example

[Writing Assertive Code With Elixir](https://dashbit.co/blog/writing-assertive-code-with-elixir)

First attempt
```js
`
Yz>get_token "foo=bar&token=value&bar=baz"
...value
`
get_token {
  string String
  parts: string.split '&'
  parts.find {
    pair String
    key value : pair.split '=' // should fail if split returns != 2
    key == token && { value }
  }
}

```
