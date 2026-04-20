#example 

https://github.com/letsgetrusty/combinators/blob/master/src/final.rs


```js
Student: {
	name String
	gpa Decimal
}

main: {
	s1: "Bogdan 3.1"
    s2: "Wallace 2.3"
    s3: "Lidiya 3.5"
    s4: "Kyle 3.9"
    s5: "Anatoliy 4.0"

    students: [s1, s2, s3, s4, s5]

	students.map({ s String; s.split(" ") })
	  .map({ p [String]
		name: p.get(0).or("")
		gpa: decimal.parse(p.get(1).or("0.0")).or(0.0)
		Option.Some(Student(name, gpa))
  	  })
	   .flatten()
	   .filter({ student Student; student.gpa >= 3.5 })
	   .collect()
	   .each({ one Student; println("`one`") })
}

```