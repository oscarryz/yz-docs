#example

https://www.golangprograms.com/go-language/concurrency.html

```go
package main

import (
	"hash/fnv"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Number of philosophers is simply the length of this list.
ph: ["Mark" "Russell" "Rocky" "Haris" "Root"]
hunger: 3                // Number of times each philosopher eats
think: time.second / 100 // Mean think time
eat: time.second / 100   // Mean eat time

fmt: log.new(os.Stdout  "" 0)

dining sync.WaitGroup

diningProblem:{
    phName String
    dominantHand sync.Mutex
    otherHand sync.Mutex

	print("{phName} Seated")
	h : fnv.new64a()
	h.Write(phName.to_byte_array())
	rg: rand.new(rand.new_source(int64{h.Sum64()}))
	rSleep : {
	    t time.Duration
	    time.sleep(t/2 + time.Duration(rg.Int63n(int64(t))))
	}
	h: hunger
	while { h > 0 } {
	   h = h - 1
		fmt.Println(phName, "Hungry")
		dominantHand.Lock() // pick up forks
		otherHand.Lock()
		fmt.Println(phName, "Eating")
		rSleep(eat)
		dominantHand.Unlock() // put down forks
		otherHand.Unlock()
		fmt.Println(phName, "Thinking")
		rSleep(think)
	}
	fmt.Println(phName, "Satisfied")
	dining.Done()
	fmt.Println(phName, "Left the table")
}

main:{
	fmt.Println("Table empty")
	dining.Add(5)
	fork0 : &sync.Mutex{}
	forkLeft : fork0
	for i := 1; i < len(ph); i++ {
		forkRight := &sync.Mutex{}
		go diningProblem(ph[i], forkLeft, forkRight)
		forkLeft = forkRight
	}
	go diningProblem(ph[0], fork0, forkLeft)
	dining.Wait() // wait for philosphers to finish
	fmt.Println("Table empty")
}
```




```js
ph: ['Mark' 'Russel' 'Rocky' 'Haris' 'Root']
hunger: 3
thing = 100
eat: 100

dining_problem: {
    ph_name String
    dominant_hand ?
    other_hand ?

    r: random.int(len(ph))

}
```
