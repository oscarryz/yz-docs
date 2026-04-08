#feature

This might not be straightforward but I might need to implement something like this to break loops


Problem: 

Iterate a list and if a condition is met, exit. 
The following won't work becuase the return with non-local return will exit the whole thing

```yz
    max_from_list: {
      list List
      m Integer
      list.for_each {
        item Integer
        item < 0 ? {
          return // will exit the max_from_list function 
        },{}
        m = max(m, item)
      }
    }
```

The solution (albeit a bit complex) is to use the non-local return to exit another function and make the whole thing run inside another block that has the exit.
    
```yz
    max_from_list: {
      
      list List(String)
      m Integer
      b : Block {
        code: {
            list.for_each {
              item Integer
                item < 0 ? {
                b.exit() // will exit the max_from_list function 
              }
              m = max(m, item)
            }
        }
      }
      b.run()
    
    }
    
    Block: {
      
      code {}
      exit {}
      run: {
        exit: {
          return // exits from run method 
        }
        code()
      }
      exit: {
        exit()
      }
      with_exit: {
    
      }
    
    }
```

Probably the following would work too.

  
```yz
    list List(Int)
    m Integer
    b: {
        code: {
            list.for_each {
                item Integer
                item < 0 ? { 
                    b.break()
                }
                m: max(m, item)
            }
        }
        run: {
            break: {
                return
            }
            code()
        }
    }
    b.run()
    return m
```


Revisited: just add a `break` keyword

```yz
max_from_list: {
  list []
  m Int
  list.for_each {
    item Int
    item < 0 ? {
      break // will exit the loop 
    }
    m = max(max, item)
  }
  m // m is undefined
}
```


Other things

```yz
stop: 50
a:0
b:0
c:1
print a + c + c

1 .to stop  {
    a = b
    b = c
    c = a + b
    print c
}
```

