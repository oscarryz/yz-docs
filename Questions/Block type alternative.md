What if the block type could use `::` and `;`

eg
```javascript
// current
walkie_talkie: {
  thing {talk {String} walk {String}} 
  print thing.walk()
  print thing.talk()
}
// this suggestion
walkie_talkie: {
  thing :: talk :: String; walk :: String
  print thing.walk()
  print thing.talk()
}
