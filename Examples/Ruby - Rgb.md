https://youtu.be/Bxel30L6C5U?t=324

```js
Rgb: {
  r Decimal
  g Decimal
  b Decimal
}
darken: {rgb Rgb
  scale: 0.5
  Rgb{ rtb.to_a.map { x Decimal; scale * x }}
}

NAMED_COLORS: [
  'red':    Rgb{ 1.0 0.0 0.0 }
  'yellow': Rgb{ 1.0 1.0 0.0 }
  'blue':   Rgb{ 0.0 0.0 1.0 }
]
names: ['red', 'yellow', 'blue']
rgbs: names.map { name String; NAMED_COLORS[name] }
darks: rgbs.map darken
1 .to names.len(), { i Int
  print '{names[i]} {rgbs[i]} {darks[i]}'
}
```