#example

https://github.com/carbon-language/carbon-lang#readme

```js
sorting: {
    partition: {
        s [T]()
        i: -1
        s.for_each({
            it T
            it <= s[s.len()-1] ? {
                i = i  + 1
                swap(s[i], e)
            }
        })
        i
    }
    quick_sort: {
        s [T]()
        s.len() <= 1 ? {
            return
        }
        p: partition(s)
        quick_sort(s.sub_array(0, p - 1))
        quick_sort(s.sub_array(p +, 1 ))
    }
}
```


```js
geometry: {
    Circle: {
        r Float
    }
    print_total_area: {
        circles [Circle]
        area Float = 0.0
        circles.for_each({it Circle
            area = area + math.pi * c.r * c.r
        })
        print("Total area: `area`")
    }
    main: {
        circles: [Circle(1.0) Circle(2.0)]
        print_total_area(circles)
    }
}
```
