
```js
// Given a type `HelloWorldModel`
HelloWorldModel : {
  saying String
}

// We can create an instance using `(args)` e.g.
model: HelloWorldModel (
  saying: 'Hello world'
)


// We can name attributes and have nested initilization
// e.g. the property `content` below
win: Frame (
    title: '`model.saying` YzFx'
    width: 200
    content: TextField(
        value: model.saying
    )
    visible: true
)
```

