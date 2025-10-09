#example
https://vosca.dev/p/17d52c0759

```js
// Experimental, if no type is defined, the variable is generic and will be binded on the first usage
List: {
    T
    data [T]()
    push: {
        val T
        data << val
    }
    pop: {
        data.last()
    }
}
list_of: {
   T
   List(T)
}
main: {
    string_list: List(String) 
    string_list.push 'hello' 
    string_list.push 'world'

    last_string: string_list.pop() // returns a String 

    bool_list: list_of(Bool) 
    bool_list.push true 
}

```