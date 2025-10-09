#feature
A single uppercase letter will mean the data type is generic:

```js
Box: {
	data T // T is generic
}
word : Box('hello') // `T` is ~~bound~~ instantiated to String
s String = word.data  // ok
// i Int = word.data // compilation error

number: Box(1) // 'T' is bound to Int
i Int = number.data // ok 
// s String = number.data // compilation error 

```

It can usually be inferred when assigning values, like `Box('hello')` above, but sometimes 
we don't want/can create the instance with values, in that case a lonely type can be set and then in the declaration be defined. 


```js
Box : {
	T // box is generic 
	data T // we might not want to pass this when building it
}
//word: Box() // compilation error, what's the value of T
word: Box(String) // initializes T to string
word.data = 'hello' // ok

// in declarations: 
a_number Box(Int) // a number will be initalized to a Box of int
a_generic_box Box(T) // the `a_generic_box` will be initialized to something in the future
a_generic_box = Box(Bool) // 

````


The generic data for a variable can be used in other elements of the block like internal variables. 
The actual type of the generic is preserved
```js

Node : {
	data T
	left Node(T) = optional.empty()
	right Node(T) = optiona.empty()
}

int_root: Node(1)
int_root.left = Node(2)
int_root.right = Node(3)

s_root: Node('hi')
s_root.left = Node('left')
// compilation error
// s_root.right = Node(1) 
```

```js
create_array: { 
	T
	[T]()
}
a [Int] = create_array(Int) 
```

The generics are flexible and allow to use any method from the given type, the compiler will verify the argument / assignment matches the methods used 

```js
say_hi: {
  data T
  print("Hello `data.name()`)
}
thing: {
  name : {"thing"}
}
blah: {}
say_hi(thing) // works, prints "Hello thing"
say_hi(blah) // doesn't compile: blah #() doesn't have a name #(String) method
```




#answered  
