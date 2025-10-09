#example


https://ballerina.io/why-ballerina/concurrent/

[First example](https://ballerina.io/why-ballerina/concurrent/#:~:text=import%20ballerina/io,y)


For a list of persons and a list of quantities we want to calculate the average. 

In the example, the sum for the quantities and the total number of employees want to be calculated concurrently. 

Also, we want to calculate the message `Employed members: N` and `Average : A`

So, we want to calculate concurrently as much as needed, until we need the value

```js
Person : {
	name String
	employed Bool
}

process #(members [Person](), quantities [Int]()) {
	// creates two variables and continues the flow
	employee_count, total_message: employee_count(members)
	// employee count is passed but still not needed, so it does't block
	avg, avg_message : calculate_average(quantities, employee_count)
	print(total_message)
	print(avg_message)
}
employee_count #(members [Person]()) {
	// count is calculated async, but still not used
    count: members.filter({p Person; p.employed}).len()
	
	// tbd if this constitutes usage, but if doesn't 
	// it can keep going and this value
	"Employed members: `count`"
}
calculate_average(quantities [Int](), employed_cont Int ) {
	total: quantities.sum()
	// employed_count is needed here, the above can run aysn
	// but the flow will stop here to let the calculation finish
	// tbd, what determines "usage", is because it is inside a`match`? 
	// is it because it is compared vs. 0? 
	// is it because it used as arg to `/` ? (it should b)
	match { 
		employed_count == 0 => 0
	}, {
	    total / employed_count	
	}
}
```

```js
// Same as above but comments removed
Person : {
	name String
	employed Bool
}

process #(members [Person](), quantities [Int]()) {
	employee_count, total_message: employee_count(members)
	avg, avg_message : calculate_average(quantities, employee_count)
	print(total_message)
	print(avg_message)
}
employee_count #(members [Person]()) {
    count: members.filter({p Person; p.employed}).len()
	"Employed members: `count`"
}
calculate_average(quantities [Int](), employed_cont Int ) {
	total: quantities.sum()
	match { 
		employed_count == 0 => 0
	}, {
	    total / employed_count	
	}
}
```

```js

// Same as above bur more of a transliteration of the ballerina sample
// Sept 12, 2025
Person : {
	name String
	employed Bool
}

process #(members [Person](), quantities [Int]()) {

	employed_count, count_msg :  { 
	    employedMembers: members.filter({p Person; p.employed})
		count : employedMembers.len()
		"Employed members: `count`"
    }()
	
	avg, avg_msg : {
		total: quantities.sum()
		match { 
			employed_count == 0 => 0
		}, {
			total / employed_count
		}
		"Average: `{avg}`"
	}()
	print(count_msg)
	print(avg_msg)
	

}
```
