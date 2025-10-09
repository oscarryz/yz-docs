#example

```typescript
// Tagged Union Types for modelling state that can be in one of many shapes
type State =
  | { type: "loading" }
  | { type: "success", value: number }
  | { type: "error", message: string };

declare const state: State;
if (state.type === "success") {
  console.log(state.value);
} else if (state.type === "error") {
  console.error(state.message);
}

```

Would be
```js
State: {
    type  String
    value Int
    message String
}
loading: State{ type: 'loading'}
success: State{ type: 'success'}
error: State{ type: 'error'}
state State
...
when(
    {state.type =='sucess'}, {print('{state.type}'},)
    {state.type =='error'}, {print('{state.message}'},)

})
```

```js
// https://stackoverflow.com/questions/71948940/copying-a-static-variable-in-copy-constructor

Park: {
    tickets [T]()icket
}
Ticket: {
    id Int
}
count Int
new_ticket: {
    count = count + 1
    Ticket: {count}
}
copy: {
    t Ticket
    Ticket{t.id}
}
```
