```javascript
HelloWorldModel: {
  saying String
}
model: HelloWorldModel (
  saying: 'Hello world'
)


win: Frame (
    title: '$(model.saying) JavaFX'
    width: 200
    content: TextField(
        value: model.saying
    )
    visible: true
)
```

