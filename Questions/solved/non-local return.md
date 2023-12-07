```js
// How would we do a break in a for loop? 
a: [1,2,3,4]

b: a.filter {it Int; it % 2 == 0}  // b: [2,4]

```

`break` and `continue`
TODO: have to define how exactly `break` and `continue` will work

```javascript

a.for_each { it Int 
   it > 2 ? {
       break // should finish the break and t
   } 
   print '{it}'
}
print 'this'
// 1 2 this
print 'hola'
a.for_each { it Int
    it == 2 ? { 
        continue // will skip the print and continue with 3 4 etc
    }
    print '{it}'
}
// 1 3 4    
```

On blocks: 


```javascript
say_hi: {
    is_time ? {
        return false
    }
    logic()
    logic()
    true
}

// return finishes the execution returns `it`
find: <A>{ collection []A; item A 
    collection.for_each { it A 
        it == item ? {
            return it       
        }
        some()
    }
    other_logic()
    return not_found
}
// `break` stops the iteration, and skips `some()` but `other_logic()` is executed 
find: <A>{ collection []A; item A 
    collection.for_each { it A 
        it == item ? {
            break
        }
        some()
    }
    other_logic()
    return not_found
}
// `continue` skips `some()` but  `other_logic()` is executed 
find: <A>{ collection []A; item A 
    collection.for_each { it A 
        it == item ? {
            continue
        }
        some()
    }
    other_logic()
    not_foun
}

```

```go

a.for_each { it Int 
   it > 2 ? {
       break // should finish the break and t
   } 
   print '{it}'
}

func (a *ArrayBlock) for_each {
    
    logic:= func(b *block) {
        if b.it > 2 {
            return
        } else {
            print()
        }
    }
    for x:= range a.items {
        b = &block{}
        b.it = x
        b.visit()
    }
}

a : {
    x()
    i : 1
    r: cond.value()
    r.?.ifTrue = {
        print 'Si'
    }
    r?r:= f.?.value()
}
type A struct { 
    i *Int
    r *a$1
    
}
type X struct {

}
func (x *X) value() {}
func (a *A) value() {
    x.value()
    x.i = Int{1}
    r:=cond.value()
    r.$qm.ifTrue = a$1{}
    a.$1= a$1r:= r.$qm.value()
}




```

 