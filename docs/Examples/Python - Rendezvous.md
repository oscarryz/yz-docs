#example

https://rosettacode.org/wiki/Rendezvous#Python

```js
Printer: {
	name String
	backup: Option.None() // Option(Printer)
	ink_level: 5
	output: io.stdout
    // print({String Result})
	print: {
	    msg String
		ink_level > 0 ? {
			output.print('(${name}): ${msg}')
			ink_level = ink_level - 1
			Result.Ok()
		}, {
			backup.Some ? { p Printer
			   p.print(msg)
			}, {
			   Result.Err('Out of ink error ${name}')
			}
		}
	}
}
main: {
	reserve: Printer('reserve')
	main: Printer('main', reserve)
    humpty_lines: [
        "Humpty Dumpty sat on a wall.",
        "Humpty Dumpty had a great fall.",
        "All the king's horses and all the king's men,",
        "Couldn't put Humpty together again."
    ]

    goose_lines: [
        "Old Mother Goose,",
        "When she wanted to wander,",
        "Would ride through the air,",
        "On a very fine gander.",
        "Jack's mother came in,",
        "And caught the goose soon,",
        "And mounting its back,",
        "Flew up to the moon."
    ]
    print_humpty: {
        humpty_lines.each({
	        line String
            main.print(line).or({
	            print('\t Humpty Dumpty out of ink!')
				break
            })
        })
    }
    print_goose: {
        goose_lines.each({
	        line String
            main.print(line).or({
	            print('\t Mother Goose out of ink!')
				break
            })
        })
    }
    print_goose()
    print_humpty()
}
```