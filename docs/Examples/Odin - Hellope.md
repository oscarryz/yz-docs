https://odin-lang.org/

```yz
main: { 
	program : "+ + * 😃 - /"
	accumulator: 0
	
	program.split("").each({ token String 
		match 
		 { token == '+' => accumulator += 1 },
	     { token == '-' => accumulator -= 1 },
		 { token == '*' => accumulator *= 2 },
		 { token == '/' => accumulator /= 2 },
		 { token == '😃' => accumulator += accumulator },
	})
	print('The program "${program}" calculated the value ${accumulator}")
	
}
```
