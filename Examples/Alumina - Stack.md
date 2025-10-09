#example

https://github.com/tibordp/alumina

```js
stack: {
    new: {
        with_capacity(0)
    }
    with_capacity: {
        capacity Int
        // capacity is not used but in the
        // oriinal example it was used to
        // create an array of size capacity
        Stack(data:[T](),len: 0)
    }
    Stack: {
        data [T]()
        len Int
        reserve: {
            additional Int
            additional > data.len() ? {
                max: std.cmp.max
                // no need to resize but left here as transcribed
                data = data.realloc(max (data.len * 2, len + additional))
            }
        }
        push: {
	    value T
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

    while ({v.is_empty() == false}, {
        print('`v.pop()`')
    })
}

```
