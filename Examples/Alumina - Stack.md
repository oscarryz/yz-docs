https://github.com/tibordp/alumina

```javascript
stack: {
    new: {
        with_capacity(0)
    }
    with_capacity: { 
        capacity Int
        Stack {
            data:[]
            len: 0
        }
    }
    Stack {
        data []
        len Int
        reserve: {
            additional Int
            len = additional > data.len() ? {
                max: core.cmp.max
                data = data.realloc max  data.len * 2 len + additional  
            }
        }
        push: { 
	        value  
            reserve(1)
            data[len] = value
            len = len + 1
        }
        << : push 
        pop: {
            len = len -1
            data[len]
        }
    
        is_empty: {
            len == 0
        }
        free: {
            data.free()
        }
    }
}
main: {
    v: stack.new()
    v << 'Stack\n'
    v << 'a '
    v << 'am '
    v << 'I '

    while {v.is_empty() == false} {
        print '$(v.pop())'
    }
}

```