https://github.com/letsgetrusty/combinators/blob/master/src/final.rs


```js
Student : {
	name String
	gpa Decimal
}

main: { 
	s1 : "Bogdan 3.1"
    s2 : "Wallace 2.3"
    s3 : "Lidiya 3.5"
    s4 : "Kyle 3.9"
    s5 : "Anatoliy 4.0"
    
    students : [s1, s2, s3, s4, s5]
    
    good_students : students.map({
        s Student
	    p : s.split(" ")
	    name : p[0]
	    gpa  : p[1]
	    Some(Student(name, gpa)) 
    })
    .filter({s Student ; s.gpa >= 3.5})
    .collect()
    
    good_students.each({one Student; println("`one`")})
}
```