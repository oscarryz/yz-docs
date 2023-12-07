http://docs.witheve.com/v0.3/tutorials/quickstart/
```javascript
w: commit {name: "Greetings"; message: "Hello world"}
m: search("Greetings").message
bind{type:"ui/text"; text: m}

```