Iterate over an array or `Iterable`

```js
array [String] = some_array
for {i Int { i < array.len() } {i = i + 1} {
    e: array(i)
    print("$(e)")
}}

for (index Int, condition (Bool), increment (),  body ()) = {
   condition() ? {
       body()
       increment()
       index = index + 1
   }
   for(index condition increment body)
}
// 
for (i Int, { i < array.len() }, {i.++()} {
    e: array[i]
    print("$(e)")
})

```