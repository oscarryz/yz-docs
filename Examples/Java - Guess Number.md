#example

https://github.com/contextfreecode/procfun/blob/main/Guess.java

```js
// Pick a random number between 1 and 100
// ask to guess
// evaluate and report guess
// update answer and counter
high: 100
answer: pick_answer high
game = Game(answer high)
game.play()
print('Finished in `game.guesses`')
print('Total input errors `game.error_count`')
pick_answer: { high Int random.next_int high}
Game {
    answer      Int
    high        Int
    done        Bool
    guesses     Int
    error_count Int

    play: {
       done == false ? {
           guess: ask_multi()
           report guess
           update guess
           play()
       }
    }

    ask_multi: {
         guess: ask_guess()
         guess == numbers.invalid ? {
             print("I didn't understand")
             error_count = error_count + 1
             ask_multi()
         }
         guess
    }

    ask_guess: {
        text = input 'Guess a number between 1 and `high`'
        numbers.parse_int text
    }

    report: { guess Int
       description: when [
           { guess < answer}: {'too low'}
           { guess > answer}: {'too high'}
           { when.else     }: {'the answer!'}
       ]
       print('`guess` is `description`')

    }

    update: { gues Int
        guess == answer ? {
            done = true
        }
        guesses = guesses + 1
    }
}

```


https://github.com/contextfreecode/procfun/blob/main/guess.hs

```js
// Same as above, but trying to avoid keeping state and passing it around
// Good example on how default values and subsequent calls work together.
high: 100
answer: pick_answer high
guesses, errors: play answer, high
print('Finished in `guesses`')
print('Total input errors `error_count`')
pick_answer: {high Int; random.next_int high }
play: {
    answer  Int
    high    Int
    done    Bool
    guesses Int
    errors  Int

    guess error_count _ : ask_multi()
    errors_count != 0 ? {
        errors = errors + error_count
    }
    report guess answer
    done guesses: update guess answer guesses


    done == false ? {
       play(answer, high, done, guesses + 1, errors)
    }
}

ask_multi: {
    guess: ask_guess()
    error_count: 0
    guess == numbers.invalid ? {
      print("I didn't understand")
      ask_multi(error_count: error_count + 1)
    }
}

ask_guess: {
    numbers.parse_int  input 'Guess a number between 1 and `high`'
}

report: { guess Int
   description: when [
       { guess < answer}: {'too low'}
       { guess > answer}: {'too high'}
       { when.else     }: {'the answer!'}
   ]
   print('{guess} is {description}')
}

update: {
    guess Int
    answer Int
    guesses Int
    done Bool
    guess == answer ? {
        done = true
    } {
        guesses = guesses + 1
    }
}

```
