#example
https://rosettacode.org/wiki/Dot_product

```javascript
dot_product: {
  a []Int
  b []Int

  a.len() != b.len() ? {
   return Err('vectors must be same length')
  }
  sum: 0
  0.to(a.len()-1).do { i Int
    sum = sum + (a[i] * b[i])
  }
  Ok(sum)
}
sum: dot_product([1 3 -5] [4 -2 -1]).or_error()



```

