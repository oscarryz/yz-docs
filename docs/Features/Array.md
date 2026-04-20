#feature

Arrays are ordered, homogeneous collections. The type of an array whose elements are `T` is written `[T]`.

## Declaration

```yz
// Declare an array variable
a [Int]

// Declare and initialize with a literal
a [Int] = [1, 2, 3]

// Short declaration (type inferred)
a: [1, 2, 3]

// Empty array
a: [Int]()
```

## Literal syntax

Array literals use square brackets with comma-separated elements:

```yz
nums: [1, 2, 3]
words: ["hello", "world"]
```

## Indexed access

```yz
nums: [10, 20, 30]
print(nums[0])   // 10

nums[0] = 99
print(nums[0])   // 99
```

## Append

Use the `<<` non-word method to append an element:

```yz
a: [Int]()
a << 1
a << 2
print(a[0])   // 1
```

## Length

```yz
a: [1, 2, 3]
print(a.length())   // 3
```

## Higher-order methods

### filter

Returns a new array containing only the elements for which the boc returns true:

```yz
nums: [1, 2, 3, 10, 20]
big: nums.filter({ item Int; item > 10 })
// big = [20]
```

Trailing-block syntax works too:

```yz
big: nums.filter { item Int; item > 10 }
```

### each

Iterates over every element, calling the boc with each one:

```yz
nums: [1, 2, 3]
nums.each({ item Int; print(item) })
```

### map

Returns a new array by applying the boc to every element:

```yz
nums: [1, 2, 3]
doubled: nums.map({ item Int; item * 2 })
doubled.each({ item Int; print(item) })
// prints 2, 4, 6
```
