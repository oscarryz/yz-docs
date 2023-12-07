https://vosca.dev/p/17d52c0759

```javascript
// Experimental, if no type is defined, the variable is generic and will be binded on the first usage
List: {
    data []
    push: {
        val
        data << val
    }
    pop: {
        data.last()
    }
}
list_of: {
   List{} 
}
main: {
    string_list: List{} // no data
    string_list.push 'hello' // bebomes a list of String
    string_list.push 'world'

    last_string: string_list.pop() // returns a String 

    bool_list: list_of() // Still nothing :/ not sure 
    bool_list.push true // list of String
}

```