
```js
HelloWorldModel: {
  saying String
}
model: HelloWorldModel (
  saying: 'Hello world'
)


win: Frame (
    title: '`model.saying` YzFx'
    width: 200
    content: TextField(
        value: model.saying
    )
    visible: true
)
```

