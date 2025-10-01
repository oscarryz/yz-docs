#example

https://youtu.be/Bxel30L6C5U?t=324

```js
Rgb: {
  r Decimal
  g Decimal
  b Decimal
}
darken: {rgb Rgb
  scale: 0.5
  Rgb(
    rgb.r * scale,
    rgb.g * scale,
    rgb.r * scale
  )
}

NAMED_COLORS: [
  'red':    Rgb( 1.0, 0.0, 0.0 )
  'yellow': Rgb( 1.0, 1.0, 0.0 )
  'blue':   Rgb( 0.0, 0.0, 1.0 )
]
names: ['red' 'yellow' 'blue']
rgbs: names.map { name String; NAMED_COLORS[name] }
darks: rgbs.map(darken)
names.len().times({ i Int
  print('`names[i]` `rgbs[i]` `darks[i]`')
})
```


// Enums / sum types exploration
```js
Rgb: {
  r Decimal
  g Decimal
  b Decimal
  Red(1.0, 0.0, 0.0)
  Yellow(1.0, 1.0, 0.0)
  Blue(1.0, 0.0, 0.0)
}
named_colors: [
  'red': Red
  'yellow': Yellow
  'blue': Blue
]
names: ['red' 'yellow' 'blue']
rbgs: names.map { s String
  named_color[name]
}
darks: rbgs.map { rbg Rbg
  scale: 5.0
  Rbg(r * scale, g * scale, b * scale)
}
1.to(names.len()).each({ i Int
  print("`names[i]` `rbgs[i]` `darks[i]`")
}



```
