#example

[Notes on 1ML](https://shonfeder.github.io/themata/programming/notes-on-1ml.html)


```js
// Signature / interface
stack #(
	empty #(T, Stack(T))
	Stack #(
		T
		push     #(a T)
		pop      #(Opt(T))
		is_empty #(Bool)
		size     #(Int)
	)
)
// usage

s : stack.empty(String)
s.push("hola")
s.is_empty() // false
s.size() // 1


// implmentation
stack = {
	empty = {
		T
		Stack(T)
	}
	Stack = {
		T
		array [T] = [T]()
		push = {
			a T
			array.<<(a)
		}
		pop = {
			array.len() > 0  ? {
				v: opt.Some(array[array.len()-1])
				array.remove(array.len())
				v
			}, {
				opt.None()
			}
		}

		is_empty = {
			array.len() == 0
		}
		size = {
			array.len()
		}
	}
}
```
