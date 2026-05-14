#example

[Notes on 1ML](https://shonfeder.github.io/themata/programming/notes-on-1ml.html)


```

```js
// Signature / interface
stack #(
	empty #(T, Stack(T))
	Stack #(
		T
		push     #(a T)
		pop      #(Option(T))
		is_empty #(Bool)
		size     #(Int)
	)
)
// usage

s: stack.empty(String)
s.push("hola")
s.is_empty() // false
s.size() // 1


// implementation
stack: {
	empty: {
		T
		Stack(T)
	}
	Stack: {
		T
		array: [T]()
		push: {
			a T
			array.push(a)
		}
		pop: {
			array.length() > 0 ? {
				v: Option.Some(array[array.length()-1])
				array.remove(array.length())
				v
			}, {
				Option.None()
			}
		}

		is_empty: {
			array.length() == 0
		}
		size: {
			array.length()
		}
	}
}
```
