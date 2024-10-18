Kind of want to do enums but they would have to:
- Have structural typing too
- Have different methods
- Stuff

```js
//rejected idea
Result {
   Ok {
      v
    }
    Err {
      cause
    }
}
```
#answered The current mechanism should suffice. Use instances and / or structural typing where applicable nd move on
```js
dow: {
    Day {
       name Stirng
       desc String =""
     }
     MONDAY: Day('monday')
     TUESDAY: Day('tuesday')
      ...
     FRIDAY {
         name : 'friday'
         desc: 'Im in love'
         end_of_week: True()
     }
  }
 ...
 today dow.Day

  when_eq today [
     { dow.FRIDAY } : { dow.FRIDAY }
  ]
```