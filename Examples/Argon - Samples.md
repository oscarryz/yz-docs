#example

[arlang.io](https://arlang.io)

```js
NOBLE_GAS = ['Helium', 'Neon', 'Argon', 'Krypton', 'Xenon']

group_by_name_length : enum.group_by({s String; s.len()})
groups:  group_by_name_length(NOBLE_GAS)
print("`groups`")

```

```js

NameCounter: {
	name String
	counter Int
	inc_counter: {
		counter = counter.++()
	}
}
obj_counters : [
	NameCounter('Alice', 0)
	NameCounter('Bob', 0)
]
obj_counter.each({ itm NameCounter
	itm.inc_counter()
}

// Also this without the inc_counter method
// and using a map instead
obj_counters : [
	{ 'Alice': 0 }
	{ 'Bob': 0 }
]
obj_counters.map  {
   itm NameCounter
   itm().1 = itm().1 + 1
}
// although, why would you?

```


```js
task: {
	id T
	{
		print("Hello from task `id`")
	}()
}
i.to(10000).do({ i Int, task(i) })
```
