#example
```js
User: {
 name String
 age Int
 is_registered Bool
 can_register: { age > 16 }
 register: {
    is_registered = true
 }
}
main: {
  s: '[{"name":"Frodo", "age": 25}, {"name":"Bobby", "age":10}]'
  users:  json.decode([User]() s) .or {
     printerr 'Failed to parse json'
  }
  users.for_each : { index Int user User 
      print('{index} {user.name}')
      (user.can_register == false).? {
           print('Cannot register {user.name}, they are too young')
           continue
      }
      user.register()
  }
  print('')
  print(json.encode users)
}
```