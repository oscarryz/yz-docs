

https://rosettacode.org/wiki/Rendezvous#Python

```javascript
Printer {
	name String
	backup Optional = optional.none// Printer
	ink_level: 5
	output: io.stdout
    // print {String Result}
	print: {
	    msg String
		ink_level > 0  ? {
			output.print('($(name)): $(msg)')
			ink_level = ink_level.--()
			result.ok()
		} {
			backup ? { p Printer 
			   p.print(msg)
			} {
			   result.error('Out of ink error $(name)')
			}
		}
	}
}
main: {
	reserve: Printer('reserve')
	main: Printer('main' reserve)
    humpty_lines : [
        "Humpty Dumpty sat on a wall."
        "Humpty Dumpty had a great fall."
        "All the king's horses and all the king's men,"
        "Couldn't put Humpty together again."
    ]

    goose_lines : [
        "Old Mother Goose,"
        "When she wanted to wander,"
        "Would ride through the air,"
        "On a very fine gander."
        "Jack's mother came in,"
        "And caught the goose soon,"
        "And mounting its back,"
        "Flew up to the moon."
    ]
    print_humpty: {
        humpty_line.each { 
	        line String 
            main.print(line).or {
	            print('\t Humpty Dumpty out of ink!')
				break
            }
			
        }
    }
    print_goose: {
        goose_line.each { 
	        line String
            main.print(line).or {
	            print('\t Mother Goose out of ink!')
				break
            }
        }
    }
    print_goose()
    print_humpty()
}
```