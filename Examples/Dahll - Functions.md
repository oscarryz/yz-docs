#example

https://dhall-lang.org/



```js
make_user: {
  user Text
  {
    home: '/home/{user}'
    private_key: '/{home}/.ssh/.id_ed25519'
    public_key: '{private_key}.pub'
  }
}

print(to_json [)
  make_user 'bill'
  make_user 'jane'
]

```

```json
[
  {
    "home": "/home/bill",
    "privateKey": "/home/bill/.ssh/id_ed25519",
    "publicKey": "/home/bill/.ssh/id_ed25519.pub"
  },
  {
    "home": "/home/jane",
    "privateKey": "/home/jane/.ssh/id_ed25519",
    "publicKey": "/home/jane/.ssh/id_ed25519.pub"
  }
]
```
