#example
From: http://rosettacode.org/wiki/24_game

```js
24_game: {
  // Gets 4 random numbrs from 1 ..9 inclusive
  get_random_numbers: {
    [random.get(1, 9), 
     random.get(1, 9), 
     random.get(1, 9), 
     random.get(1, 9)]
  }
  // Ask user for an aritmetic operation
  prompt_user: {
    numbers [] Int
    expression: input("Enter an aritmetic expression to add up to 24 using `numbers[0]` `numbers[1]` `numbers[2]` and `numbers[3]`.
    You can only use the following operators: + - / *")
    eval expression 
  }

  // Eval user aritmetic operation
  eval: {
    expression String 
    stack util.Stack{}
    operations: [
        '+': {a Int; b Int; a + b}, 
        '-': {a Int; b Int; a - b}, 
        '*': {a Int; b Int; a * b}, 
        '/': {a Int; b Int; a / b}
    ]
      
    expression.each({ item String
        match {
          ['+', '-', '*', '/'].contains(item) => {
            a: stack.pop()
            b: stack.pop() 
            stack.push(operators[item](a, b))
          }
        }, {
          [' ', "\n"].contains(item) => {
            // skip whitespace
          }
        }, {
          allowed_numbers.contains(item) => { 
            stack.push(numbers.parseInt(item))
          }
        }, {
          panic("Invalid character `item`")
        }
    }
    stack.pop()
  }  
  
  main: {
    numbers:get_random_numbers()
    expression: prompt_user numbers
    result: eval expression  numbers
    print(result == 24 ? {
      "You got it!"
    } {
      "I'm sorry that doesn't evaluate to 24"
    })
}

```

Actual go implementation 
```go
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	numbers := getRandomNumbers()
	expression := promptUser(numbers)
	if r := evaluate(expression, numbers); r == 24 {
		fmt.Println("You got it")
	} else {
		fmt.Printf("Sorry, you've got %d\n", r)
	}
}

func getRandomNumbers() []int {
	var numbers []int
	rand.Seed(time.Now().UnixNano())
	numbers = append(numbers, rand.Intn(9)+1)
	numbers = append(numbers, rand.Intn(9)+1)
	numbers = append(numbers, rand.Intn(9)+1)
	numbers = append(numbers, rand.Intn(9)+1)
	//return []int{2, 4, 5, 1}
	return numbers
}
func promptUser(numbers []int) string {
	fmt.Printf("You know the drill %v: ", numbers)
	reader := bufio.NewReader(os.Stdin)
	if line, err := reader.ReadString('\n'); err != nil {
		fmt.Println(err)
		return "0"
	} else {
		return line
	}
}
func evaluate(expression string, numbers []int) int {
	var stack []int
	for _, r := range expression {
		// only \s or \d or [+-*/]
		// eg.  1 2 + 3 *  4 / -> (1 + 2) * 3 / 4
		// 1 2 3 4 + -
		// 2 3 *
		// fmt.Printf("stack %v\n", stack)
		c := string(r)
		if c == "+" || c == "-" || c == "*" || c == "/" {
			b := stack[len(stack)-1]
			stack = stack[0 : len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[0 : len(stack)-1]
			switch c {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				stack = append(stack, a/b)
			}
		} else {
			switch c {
			case " ":
				fallthrough
			case "\n":
				continue
			default:

				//fmt.Printf("default [%s]\n", c)
				number, _ := strconv.Atoi(c)
				isValid := false
				for _, vn := range numbers {
					if number == vn {
						isValid = true
						break
					}
				}
				if isValid {
					stack = append(stack, number)
				} else {
					fmt.Printf("Dude, c'mon: %d\n", number)
					return -1
				}

			}
		}
	}
	return stack[0]

}


```

```js
24_game: {
  // Gets 4 random numbrs from 1 ..9 inclusive
  get_random_numbers: {
    [random.get(1, 9), 
    random.get(1, 9), 
    random.get(1, 9), 
    random.get(1, 9)]
  }
  // Ask user for an aritmetic operation
  prompt_user: {
    numbers [] Int
    expression: input("Enter an aritmetic expression to add up to 24 using `numbers[0]`, `numbers[1]`, `numbers[2]` and `numbers[3]`.
    You can only use the following operators: + - / *")
    eval expression 
  }

  // Eval user aritmetic operation
  eval: {
    expression String 
    stack util.Stack{}
    take2: { a: stack.pop(); b: stack.pop() }
    expression.each({ item String
      (item == '+' || item == '-' || item == '*' || item == '/') ? {
        a, b: take2()
        stack.push(a + b)
      } {
        (item == ' ' || item == '\n') ? {
          continue
        } {
          stack.push(numbers.parse(item))
        }
      }
    }
    stack.pop()
  }
  main: {
    numbers:get_random_numbers()
    expression: prompt_user numbers
    result: eval expression  numbers
    result == 24 ? {
      "You got it!"
    } {
      "I'm sorry that doesn't evaluate to 24"
    }
}
```

Possible  `when`  definition 
```js
when: {
  conditions [{Bool}]{}
  conditions.for_each { 
    cond   {Bool}
    action {}
    cond() ? { action(); return}
  }
}

match {
  ['+', '-', '*', '/'].contains(item) => {
    a, b: take2()
    stack.push(operators[item](a, b))
  }
}, {
  [' ', "\n"].contains(item) => {
    continue
  }
}, {
  allowed_numbers.contains(item) => { 
    stack.push(numbers.parseInt(item))
  }
}, {
  print("Invalid character `item`")
  panic()
}

eval: {
    expression String 
    stack util.Stack{}
    operations: [
        '+': {a Int; b Int; a + b}, 
        '-': {a Int; b Int; a - b}, 
        '*': {a Int; b Int; a * b}, 
        '/': {a Int; b Int; a / b}
    ]
    eval2: { 
        a: stack.pop()
        b: stack.pop() 
        stack.push operators[item](a b)
    }
    continue: {}
      
    expression.each({ item String
        match {
          ['+', '-', '*', '/'].contains(item) => {
            eval2()
          }
        }, {
          [' ', "\n"].contains(item) => {
            continue()
          }
        }, {
          allowed_numbers.contains(item) => { 
            stack.push(numbers.parseInt(item))
          }
        }, {
          panic("Invalid character `item`")
        }
    stack.pop()
  }
```

