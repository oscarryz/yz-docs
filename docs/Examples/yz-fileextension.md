---
type: example
---

> freeCodeCamp.org on X: "For today's coding challenge, you'll build a little extension extractor. Given a string representing a filename, you'll return the extension of the file. Check it out on the freeCodeCamp mobile app." / X https://share.google/4jVsRLePpcZgjGU5a

```js
get_extension #(filename String, Option(String)) {
   
   extension : None(String)
   filename.length.down_to(1).each({ 
      i Int
      filename.at(i - 1) == "." ? {
         extension = Some(filename.substring(i - 1)) 
      }
   })
   extension
}
```
      